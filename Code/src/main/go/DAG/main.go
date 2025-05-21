package main

import (
	"context"
	"errors" // 保留以便兼容，尽管在简单函数中未使用
	"fmt"
	"sync"
	"time"
)

// NodeState 表示节点执行状态
type NodeState int

const (
	Pending   NodeState = iota // 等待执行
	Running                    // 正在执行
	Completed                  // 执行完成
	Failed                     // 执行失败
	Skipped                    // 因短路或上游失败而被跳过
)

// String 将 NodeState 转换为可读字符串
func (s NodeState) String() string {
	switch s {
	case Pending:
		return "Pending" // 等待中
	case Running:
		return "Running" // 运行中
	case Completed:
		return "Completed" // 已完成
	case Failed:
		return "Failed" // 已失败
	case Skipped:
		return "Skipped" // 已跳过
	default:
		return "Unknown" // 未知状态
	}
}

// NodeOutput 是节点执行的输出，键为输出参数名，值为参数值
type NodeOutput map[string]interface{}

// NodeInput 是节点执行的输入，键通常为依赖的节点ID或全局输入名，值为对应的NodeOutput或全局值
type NodeInput map[string]interface{}

// Node 定义了工作流中的一个节点
type Node struct {
	ID          string                 // 节点的唯一标识符
	Name        string                 // 节点的可读名称
	Description string                 // 节点的描述信息
	Config      map[string]interface{} // 节点特定配置

	// ExecuteFunc 是节点的实际执行逻辑
	// ctx: 上下文，可用于超时或取消
	// inputs: 从上游节点获取的输入数据
	// config: 节点自身的配置信息
	// 返回: 输出数据 (NodeOutput), 是否触发短路 (bool), 错误信息 (error)
	ExecuteFunc func(ctx context.Context, inputs NodeInput, config map[string]interface{}) (output NodeOutput, shortCircuit bool, err error)

	Dependencies []string // 此节点依赖的其他节点的ID列表

	// 内部状态字段，在执行期间由引擎更新
	state        NodeState  // 节点的当前状态
	output       NodeOutput // 节点的输出结果
	err          error      // 节点执行时发生的错误
	shortCircuit bool       // 标记此节点是否触发了短路（根据ExecuteFunc的返回值）
}

// Workflow 定义了一个工作流，包含所有节点及其依赖关系
type Workflow struct {
	ID          string              // 工作流的唯一标识符
	Name        string              // 工作流的可读名称
	Nodes       map[string]*Node    // 工作流中所有节点的映射，键为节点ID
	adj         map[string][]string // 邻接表 (节点ID -> 后继节点ID列表)，表示节点的输出连接
	revAdj      map[string][]string // 反向邻接表 (节点ID -> 前驱节点ID列表)，表示节点的输入依赖
	entryPoints []string            // 入口节点ID列表 (没有依赖的节点)
}

// NewWorkflow 创建一个新的工作流实例
func NewWorkflow(id, name string) *Workflow {
	return &Workflow{
		ID:          id,
		Name:        name,
		Nodes:       make(map[string]*Node),
		adj:         make(map[string][]string),
		revAdj:      make(map[string][]string),
		entryPoints: []string{}, // 初始化为空切片
	}
}

// AddNode 向工作流中添加一个节点
func (wf *Workflow) AddNode(node *Node) error {
	if _, exists := wf.Nodes[node.ID]; exists {
		return fmt.Errorf("ID为 %s 的节点已存在", node.ID)
	}
	wf.Nodes[node.ID] = node
	wf.adj[node.ID] = []string{}    // 初始化邻接表条目
	wf.revAdj[node.ID] = []string{} // 初始化反向邻接表条目
	return nil
}

// Prepare 在执行前准备工作流，例如计算入口节点和重置状态
func (wf *Workflow) Prepare() error {
	// 重置入口点和所有节点的状态
	wf.entryPoints = []string{}
	for _, node := range wf.Nodes {
		node.state = Pending
		node.output = nil
		node.err = nil
		node.shortCircuit = false
	}

	// 清空并重建邻接表和反向邻接表，基于 Node.Dependencies
	wf.adj = make(map[string][]string)
	wf.revAdj = make(map[string][]string)
	for id := range wf.Nodes { // 为每个节点初始化邻接表和反向邻接表的空切片
		wf.adj[id] = []string{}
		wf.revAdj[id] = []string{}
	}

	// 根据 Node.Dependencies 填充邻接表和反向邻接表
	for id, node := range wf.Nodes {
		for _, depID := range node.Dependencies {
			if _, exists := wf.Nodes[depID]; !exists {
				return fmt.Errorf("节点 %s 的依赖 %s 不存在", id, depID)
			}
			// depID 是 node 的前驱节点，id 是 depID 的后继节点
			wf.adj[depID] = append(wf.adj[depID], id)    // depID -> id
			wf.revAdj[id] = append(wf.revAdj[id], depID) // id 依赖于 depID
		}
	}

	// 查找入口节点 (没有依赖的节点)
	for id, node := range wf.Nodes {
		if len(node.Dependencies) == 0 {
			wf.entryPoints = append(wf.entryPoints, id)
		}
	}

	// 如果图不为空但没有入口节点，则可能存在循环或所有节点都有依赖（这是无效的DAG）
	if len(wf.entryPoints) == 0 && len(wf.Nodes) > 0 {
		return errors.New("工作流中未找到入口节点（可能存在循环或所有节点都有依赖）")
	}
	// 注意: 更健壮的环检测应该在AddNode或AddEdge时进行，或者在这里进行一次完整的DFS。
	return nil
}

// ExecutionContext 保存工作流执行期间的上下文信息
type ExecutionContext struct {
	workflow                 *Workflow             // 当前执行的工作流
	nodeResults              map[string]NodeOutput // 存储每个成功节点的输出结果
	nodeStates               map[string]NodeState  // 存储每个节点的状态
	nodeErrors               map[string]error      // 存储每个失败节点的错误信息
	mu                       sync.Mutex            // 互斥锁，用于保护共享资源 (如nodeResults, nodeStates, nodeErrors) 的并发访问
	wg                       sync.WaitGroup        // 用于等待所有并发的节点goroutine完成
	globalInputs             NodeInput             // 工作流的全局初始输入数据
	shortCircuitCh           chan string           // 用于通知全局短路事件，值为触发短路的节点ID。这是一个信号通道。
	shortCircuited           bool                  // 标记整个工作流是否已发生短路
	shortCircuitSourceNodeID string                // 标记触发工作流短路的源节点ID
	ctx                      context.Context       // 主上下文，用于控制工作流的取消（例如超时）
	cancelFunc               context.CancelFunc    // 主上下文的取消函数
}

// NewExecutionContext 创建一个新的执行上下文实例
func NewExecutionContext(wf *Workflow, globalInputs NodeInput, parentCtx context.Context) *ExecutionContext {
	ctx, cancel := context.WithCancel(parentCtx) // 创建可取消的子上下文
	// shortCircuitCh 的缓冲区大小至少为1，或者如果许多节点可以同时触发则为len(wf.Nodes)。
	// 为简单起见，小型缓冲区或在处理前预期只有一个触发器时使用无缓冲通道。
	// 大小为1的缓冲通道可以很好地防止在接收器尚未准备好时发送者阻塞。
	return &ExecutionContext{
		workflow:       wf,
		nodeResults:    make(map[string]NodeOutput),
		nodeStates:     make(map[string]NodeState),
		nodeErrors:     make(map[string]error),
		globalInputs:   globalInputs,
		shortCircuitCh: make(chan string, 1), // 缓冲通道，用于在已短路时非阻塞发送
		ctx:            ctx,
		cancelFunc:     cancel,
	}
}

// Run 执行工作流
func (wf *Workflow) Run(ctx context.Context, globalInputs NodeInput) (*ExecutionContext, error) {
	if err := wf.Prepare(); err != nil { // 首先准备工作流
		return nil, fmt.Errorf("工作流准备失败: %w", err)
	}

	execCtx := NewExecutionContext(wf, globalInputs, ctx) // 创建执行上下文
	// execCtx.cancelFunc() // cancelFunc会在出错或executeNode路径结束时调用

	fmt.Printf("工作流 '%s' 开始执行，有 %d 个入口节点: %v\n", wf.Name, len(wf.entryPoints), wf.entryPoints)

	if len(wf.entryPoints) == 0 && len(wf.Nodes) > 0 {
		fmt.Printf("工作流 '%s' 有节点但没有入口节点。无法运行。\n", wf.Name)
		// 如果因没有入口节点而提前退出，请确保取消上下文
		execCtx.cancelFunc()
		return execCtx, errors.New("没有入口节点来启动工作流")
	}
	if len(wf.Nodes) == 0 {
		fmt.Printf("工作流 '%s' 没有节点。无需运行。\n", wf.Name)
		execCtx.cancelFunc()
		return execCtx, nil // 或者返回一个错误，指示没有可运行的内容
	}

	// 为每个入口节点启动一个goroutine执行
	for _, nodeID := range wf.entryPoints {
		execCtx.wg.Add(1) // WaitGroup计数器加1
		go execCtx.executeNode(nodeID)
	}

	// 等待所有goroutine完成或工作流上下文被取消/超时
	waitDone := make(chan struct{})
	go func() {
		execCtx.wg.Wait() // 等待所有通过wg.Add(1)添加的任务完成
		close(waitDone)   // 关闭waitDone通道，通知主goroutine
	}()

	select {
	case <-waitDone:
		fmt.Printf("工作流 '%s' 完成了所有可达节点的处理。\n", wf.Name)
	case <-execCtx.ctx.Done(): // 如果工作流的主上下文被取消 (例如超时)
		fmt.Printf("工作流 '%s' 被取消或超时: %v\n", wf.Name, execCtx.ctx.Err())
		execCtx.mu.Lock()
		if !execCtx.shortCircuited { // 如果尚未被某个节点短路
			execCtx.shortCircuited = true
			execCtx.shortCircuitSourceNodeID = "WorkflowContextDone" // 标记为上下文取消
			// 尝试关闭通道以通知其他goroutine，确保它尚未关闭。
			// 如果多个goroutine可能尝试关闭，这部分可能很棘手。
			// 使用select和default或sync.Once进行关闭可能更安全。
			// 目前，我们假设cancelFunc和标志处理了大多数情况。
			// 如果已经由触发短路的节点关闭，则在此处关闭它可能会导致panic。
			// 相反，依赖于各个goroutine检查execCtx.ctx.Done()。
		}
		execCtx.mu.Unlock()
	}

	// 如果尚未通过短路节点或超时完成，则发送最终取消信号。
	// 这可确保资源得到清理。
	execCtx.cancelFunc()

	// 收集最终结果和状态，检查是否有未处理的错误
	var encounteredError bool
	execCtx.mu.Lock() // 为最终读取nodeErrors加锁
	for nodeID, err := range execCtx.nodeErrors {
		if err != nil {
			fmt.Printf("节点 %s 中发生错误: %v\n", nodeID, err)
			encounteredError = true
			// break // 即使有错误，也继续收集所有错误信息
		}
	}
	execCtx.mu.Unlock()

	if encounteredError {
		return execCtx, errors.New("工作流执行遇到一个或多个错误")
	}

	return execCtx, nil // 工作流成功完成（所有可达节点）
}

// skipDownstreamNodeIfParentTriggered 递归地将下游节点标记为Skipped (如果其父节点触发了跳过)
// reasonNodeID 是导致跳过的上游节点的ID
func (ec *ExecutionContext) skipDownstreamNodeIfParentTriggered(nodeID string, reasonNodeID string) {
	ec.mu.Lock()
	// 检查节点是否已被处理或跳过，以避免冗余工作/日志记录
	if state, exists := ec.nodeStates[nodeID]; exists && state != Pending && state != Running {
		ec.mu.Unlock()
		return // 如果节点已完成、失败或已跳过，则不执行任何操作
	}

	// 确保 workflow.Nodes[nodeID] 存在
	nodeToSkip, nodeExists := ec.workflow.Nodes[nodeID]
	if !nodeExists {
		// 这不应该发生，如果工作流已正确准备
		fmt.Printf("严重错误: 在skipDownstreamNodeIfParentTriggered中找不到节点 %s\n", nodeID)
		ec.mu.Unlock()
		return
	}

	fmt.Printf("节点 '%s': 尝试跳过，因为上游 '%s' 失败/短路。\n", nodeToSkip.Name, reasonNodeID)
	ec.nodeStates[nodeID] = Skipped
	nodeToSkip.state = Skipped
	nodeToSkip.err = fmt.Errorf("由于上游 %s 而跳过", reasonNodeID) // 记录跳过的原因
	ec.nodeErrors[nodeID] = nodeToSkip.err                   // 也将其视为一种“错误”状态，尽管是预期的跳过
	ec.mu.Unlock()

	// 递归地跳过其子节点
	// 重要提示：如果在递归调用之前未向下传递锁，则在递归调用之前获取adj列表。
	// 但是，adj对于工作流运行是静态的。
	for _, childNodeID := range ec.workflow.adj[nodeID] {
		// 将原始原因或当前nodeID作为新原因传递
		ec.skipDownstreamNodeIfParentTriggered(childNodeID, nodeID)
	}
}

// areDependenciesMetAndNotSkipped 检查一个节点的所有依赖是否都已'Completed'状态
func (ec *ExecutionContext) areDependenciesMetAndNotSkipped(nodeID string) bool {
	ec.mu.Lock() // 为安全读取共享状态加锁
	defer ec.mu.Unlock()

	node := ec.workflow.Nodes[nodeID]
	if node == nil {
		// 如果工作流准备正确，则不应发生这种情况
		fmt.Printf("错误: 在依赖项检查期间，工作流中找不到节点 %s。\n", nodeID)
		return false
	}
	if len(node.Dependencies) == 0 {
		return true // 没有依赖项，因此它们已满足。
	}

	for _, depID := range node.Dependencies {
		depState, stateExists := ec.nodeStates[depID]
		// 如果任何依赖项未完成，或者在状态中不存在（应表示待处理或有问题）
		if !stateExists || depState != Completed {
			return false // 任何一个依赖不是Completed，则此节点不能执行
		}
	}
	return true // 所有依赖项均为“Completed”
}

// executeNode 是执行单个节点的核心逻辑单元
func (ec *ExecutionContext) executeNode(nodeID string) {
	defer ec.wg.Done() // 确保WaitGroup计数器在goroutine退出时递减

	node := ec.workflow.Nodes[nodeID] // 获取节点实例

	// 1. 初始检查：检查工作流是否已被取消或发生全局短路
	select {
	case <-ec.ctx.Done(): // 如果工作流主上下文被取消
		ec.mu.Lock()
		// 仅当节点尚未处于终端状态时更新
		if ec.nodeStates[nodeID] != Completed && ec.nodeStates[nodeID] != Failed && ec.nodeStates[nodeID] != Skipped {
			fmt.Printf("节点 '%s' 的执行因工作流上下文取消而被取消: %v\n", node.Name, ec.ctx.Err())
			ec.nodeStates[nodeID] = Skipped // 将状态标记为Skipped
			node.state = Skipped
			node.err = ec.ctx.Err()
			ec.nodeErrors[nodeID] = node.err // 记录错误
			// 递归跳过下游节点
			downstreamIDs := ec.workflow.adj[nodeID]
			ec.mu.Unlock() // 在递归调用前解锁
			for _, dsNodeID := range downstreamIDs {
				ec.skipDownstreamNodeIfParentTriggered(dsNodeID, nodeID)
			}
			return // 退出此goroutine
		}
		ec.mu.Unlock()
		return // 节点已处于终端状态或已处理
	default:
		// 如果没有立即取消，则继续执行
	}

	ec.mu.Lock()
	// 在获取锁后再次检查：如果在等待锁期间发生了全局短路
	if ec.shortCircuited && ec.shortCircuitSourceNodeID != nodeID { // 如果已短路且非本节点触发
		if ec.nodeStates[nodeID] == Pending || ec.nodeStates[nodeID] == Running { // 仅当节点尚未完成时跳过
			fmt.Printf("节点 '%s' 因由 '%s' 触发的全局工作流短路而被跳过。\n", node.Name, ec.shortCircuitSourceNodeID)
			ec.nodeStates[nodeID] = Skipped
			node.state = Skipped
			node.err = fmt.Errorf("工作流被 %s 短路", ec.shortCircuitSourceNodeID)
			ec.nodeErrors[nodeID] = node.err
			// 递归跳过下游节点
			downstreamIDs := ec.workflow.adj[nodeID]
			ec.mu.Unlock() // 在递归调用前解锁
			for _, dsNodeID := range downstreamIDs {
				ec.skipDownstreamNodeIfParentTriggered(dsNodeID, nodeID)
			}
			return
		}
		ec.mu.Unlock()
		return // 节点已处理或处于终端状态
	}
	ec.mu.Unlock() // 在可能长时间运行的依赖项检查/输入收集之前释放锁

	// 2. 依赖检查和输入数据收集
	inputs := make(NodeInput)         // 初始化当前节点的输入数据
	allDepsCompletedOrSkipped := true // 标记是否所有必需的依赖都已“完成”。
	anyDepFailedOrSkipped := false    // 标记是否有任何依赖失败或被跳过
	reasonForSkipOrFailure := ""      // 记录跳过或失败的原因节点名称

	ec.mu.Lock() // 加锁以安全访问依赖节点的状态和结果
	for _, depID := range node.Dependencies {
		depNode, depNodeExists := ec.workflow.Nodes[depID]
		if !depNodeExists { // 如果Prepare()正确，则不应发生
			fmt.Printf("严重错误: %s 的依赖节点 %s 未找到！\n", node.Name, depID)
			node.state = Failed
			node.err = fmt.Errorf("依赖 %s 未找到", depID)
			ec.nodeStates[nodeID] = Failed
			ec.nodeErrors[nodeID] = node.err
			ec.mu.Unlock()
			return // 严重错误，无法继续
		}

		depState, stateExists := ec.nodeStates[depID]
		depOutput, _ := ec.nodeResults[depID] // 如果依赖失败/跳过，输出可能为nil

		if !stateExists || depState == Pending || depState == Running {
			allDepsCompletedOrSkipped = false // 某个依赖项尚未处于终端状态
			break                             // 无需再检查其他依赖
		}
		if depState == Failed || depState == Skipped {
			anyDepFailedOrSkipped = true
			reasonForSkipOrFailure = depNode.Name // 使用名称以便于阅读
			// 如果任何依赖项失败或被跳过，则此节点无法运行，应被跳过。
			break // 无需再检查其他依赖
		}
		// 如果 depState == Completed
		inputs[depID] = depOutput // 从已完成的依赖项收集输出
	}
	ec.mu.Unlock() // 释放锁

	if !allDepsCompletedOrSkipped {
		// 此情况意味着某个依赖项仍处于Pending或Running状态。
		// 当前模型依赖于依赖项完成时触发其子节点。
		// 如果一个节点在其依赖项进入终端状态之前被调度，则这是一个问题，
		// 或者意味着它被调度得太早/顺序错误。
		// 目前，我们假设这意味着它无法继续。
		// 在更强大的系统中，此节点可能会等待或重新排队。
		// 考虑到当前的逻辑，这暗示调度或状态更新存在问题。
		// 但是，如果一个节点被添加到WaitGroup但其依赖项尚未被另一个goroutine标记为已完成，则可能存在竞争条件导致此情况。
		// 重新检查或短暂等待可能是一个选项，但为简单起见，我们将假设
		// 在goroutine启动之前的`areDependenciesMetAndNotSkipped`主要处理了此问题。
		fmt.Printf("节点 '%s' 尚不能运行，依赖项未处于终端状态。\n", node.Name)
		return // 无法执行，等待依赖完成
	}

	if anyDepFailedOrSkipped { // 如果有任何依赖失败或被跳过
		ec.mu.Lock()
		if ec.nodeStates[nodeID] != Skipped { // 避免因全局短路而已被跳过时产生冗余日志
			fmt.Printf("节点 '%s' 因其依赖 '%s' 失败或被跳过而被跳过。\n", node.Name, reasonForSkipOrFailure)
			ec.nodeStates[nodeID] = Skipped
			node.state = Skipped
			node.err = fmt.Errorf("依赖 %s 失败/跳过", reasonForSkipOrFailure)
			ec.nodeErrors[nodeID] = node.err
		}
		// 递归跳过下游节点
		downstreamIDs := ec.workflow.adj[nodeID]
		ec.mu.Unlock() // 在递归调用前解锁
		for _, dsNodeID := range downstreamIDs {
			ec.skipDownstreamNodeIfParentTriggered(dsNodeID, nodeID)
		}
		return // 当前节点被跳过
	}

	// 如果是入口节点且存在全局输入，则合并它们
	if len(node.Dependencies) == 0 && ec.globalInputs != nil {
		for k, v := range ec.globalInputs {
			if _, exists := inputs[k]; !exists { // 如果键冲突，优先使用特定节点的输出
				inputs[k] = v
			}
		}
	}

	// 3. 执行节点逻辑
	fmt.Printf("节点 '%s' (ID: %s) 开始执行。\n", node.Name, node.ID)
	ec.mu.Lock()
	ec.nodeStates[nodeID] = Running // 更新状态为Running
	node.state = Running
	ec.mu.Unlock()

	// 节点的执行上下文（如果需要可以更具体）
	nodeCtx := ec.ctx // 继承工作流上下文以传播取消信号

	// 调用节点定义的ExecuteFunc
	nodeOutput, shouldShortCircuit, execErr := node.ExecuteFunc(nodeCtx, inputs, node.Config)

	// 4. 处理执行结果并更新状态 (在锁保护下)
	ec.mu.Lock()

	// 检查在节点执行*期间*上下文是否被取消
	if nodeCtx.Err() != nil && node.state == Running { // 检查状态以避免覆盖已处理的状态
		fmt.Printf("节点 '%s' 的执行因上下文取消而被中断: %v\n", node.Name, nodeCtx.Err())
		ec.nodeStates[nodeID] = Failed // 或Skipped，取决于中断策略
		node.state = Failed
		node.err = nodeCtx.Err()
		ec.nodeErrors[nodeID] = node.err
		execErr = node.err // 确保此错误被考虑用于短路
		// 如果被取消，则不继续触发下游
	} else if execErr != nil { // 如果ExecuteFunc返回错误
		fmt.Printf("节点 '%s' (ID: %s) 执行失败: %v\n", node.Name, node.ID, execErr)
		ec.nodeStates[nodeID] = Failed
		node.state = Failed
		node.err = execErr
		ec.nodeErrors[nodeID] = execErr
	} else { // ExecuteFunc成功执行
		fmt.Printf("节点 '%s' (ID: %s) 执行完成。\n", node.Name, node.ID)
		ec.nodeStates[nodeID] = Completed
		node.state = Completed
		ec.nodeResults[nodeID] = nodeOutput // 存储输出结果
		node.output = nodeOutput
	}
	node.shortCircuit = shouldShortCircuit // 记录节点是否希望触发短路（来自ExecuteFunc的返回值）

	// 5. 处理短路逻辑 (如果此节点触发或发生错误)
	// 节点中的错误被视为短路其依赖路径的原因。
	// `node.shortCircuit` 是指节点*明确*请求更广泛的工作流短路。
	triggeredGlobalShortCircuit := false                                        // 标记此节点是否是*首次*触发全局短路
	if node.state == Failed || (node.state == Completed && node.shortCircuit) { // 如果节点失败，或成功但要求短路
		if !ec.shortCircuited { // 仅当全局尚未短路时，此节点才能成为第一个触发者
			ec.shortCircuited = true             // 标记全局短路发生
			ec.shortCircuitSourceNodeID = nodeID // 记录是此节点触发的
			triggeredGlobalShortCircuit = true
			fmt.Printf("节点 '%s' 正在启动工作流短路。\n", node.Name)
			// 安全地尝试在通道上发送/关闭通道。关闭是广播。
			// 使用select发送可确保如果通道已满或尚无接收者，则不会阻塞。
			// 关闭更为明确。
			select {
			case <-ec.shortCircuitCh: // 通道已关闭或已发送值
			default:
				close(ec.shortCircuitCh) // 关闭通道以通知其他goroutine
			}
			ec.cancelFunc() // 同时取消主上下文以停止其他长时间运行的任务
		}
	}

	// 6. 触发下游节点或跳过它们
	downstreamIDs := ec.workflow.adj[nodeID] // 获取所有下游节点
	ec.mu.Unlock()                           // 在可能调度新goroutine或递归调用之前释放锁

	if node.state == Failed || (node.state == Completed && node.shortCircuit && triggeredGlobalShortCircuit) || (ec.shortCircuited && ec.shortCircuitSourceNodeID == nodeID) {
		// 如果此节点失败，或者它完成并触发了*新的*全局短路，
		// 或者它是现有全局短路的源头，则跳过其直接子节点。
		// 这也涵盖了节点在执行期间被取消的情况。
		fmt.Printf("节点 '%s' 导致其下游路径被跳过。\n", node.Name)
		for _, dsNodeID := range downstreamIDs {
			ec.skipDownstreamNodeIfParentTriggered(dsNodeID, nodeID) // 跳过下游
		}
	} else if node.state == Completed && !node.shortCircuit { // 如果节点成功完成且未请求短路
		// 如果由于*另一个*节点导致全局短路，则此节点的子节点仍可能被跳过。
		// 该检查位于这些子节点的executeNode的开头。
		// 在这里，我们仅当此节点正常完成时才调度子节点。
		for _, nextNodeID := range downstreamIDs {
			if ec.areDependenciesMetAndNotSkipped(nextNodeID) { // 检查下游节点的依赖是否都已满足
				// 在启动前，检查是否由其他分支引起了全局短路
				ec.mu.Lock()
				globalSC := ec.shortCircuited
				ec.mu.Unlock()
				if globalSC {
					// 如果发生全局短路，请确保也跳过此下游节点
					// 它可能尚未被源头的skipDownstreamNodeIfParentTriggered触及。
					fmt.Printf("由于工作流已全局短路，跳过启动 '%s'。\n", ec.workflow.Nodes[nextNodeID].Name)
					ec.skipDownstreamNodeIfParentTriggered(nextNodeID, ec.shortCircuitSourceNodeID) // 传递原始短路源
				} else {
					ec.wg.Add(1) // 为下游节点启动新的goroutine
					go ec.executeNode(nextNodeID)
				}
			}
		}
	}
}

// --- 简单DAG的示例ExecuteFuncs ---

// simpleStartFunc 是起始节点的简单执行函数
func simpleStartFunc(ctx context.Context, inputs NodeInput, config map[string]interface{}) (NodeOutput, bool, error) {
	nodeName := "UnknownStartNode"                    // 默认节点名
	if name, ok := config["node_name"].(string); ok { // 从配置中获取节点名
		nodeName = name
	}
	fmt.Printf("--- [%s] 正在执行起始节点 ---\n", nodeName)
	time.Sleep(100 * time.Millisecond) // 模拟工作耗时

	initialData := "来自A的初始数据"                             // 默认初始数据
	if globalVal, ok := inputs["global_start_data"]; ok { // 检查是否有全局输入
		initialData = fmt.Sprintf("来自全局输入的初始数据: %v", globalVal)
	}

	output := make(NodeOutput)     // 创建输出map
	output["data_a"] = initialData // 设置输出数据
	fmt.Printf("--- [%s] 输出: %s ---\n", nodeName, output["data_a"])
	return output, false, nil // 返回输出，不短路，无错误
}

// simpleProcessBFunc 是处理节点B的简单执行函数
func simpleProcessBFunc(ctx context.Context, inputs NodeInput, config map[string]interface{}) (NodeOutput, bool, error) {
	nodeName := "UnknownProcessBNode"
	if name, ok := config["node_name"].(string); ok {
		nodeName = name
	}
	fmt.Printf("--- [%s] 正在执行处理节点B ---\n", nodeName)
	time.Sleep(200 * time.Millisecond) // 模拟工作耗时

	var dataFromA string
	// 来自节点A的输入将位于键 "A_id" (或A的ID) 下
	// 其内容将是一个NodeOutput映射。
	for depID, depOutputInterface := range inputs { // 应该只有来自A的输入
		depOutput, ok := depOutputInterface.(NodeOutput) // 类型断言
		if !ok {
			return nil, false, fmt.Errorf("[%s] 严重错误: 来自 %s 的输入不是NodeOutput类型", nodeName, depID)
		}
		if val, valOk := depOutput["data_a"].(string); valOk { // 从A的输出中获取"data_a"
			dataFromA = val
			break // 找到来自A的数据
		}
	}

	if dataFromA == "" { // 如果没有收到来自A的数据
		return nil, false, fmt.Errorf("[%s] 错误: 未从上游接收到data_a", nodeName)
	}

	output := make(NodeOutput)
	output["data_b"] = fmt.Sprintf("B已处理: [%s]", dataFromA) // 处理数据
	fmt.Printf("--- [%s] 输出: %s ---\n", nodeName, output["data_b"])
	return output, false, nil
}

// simpleProcessCFunc 是处理节点C的简单执行函数
func simpleProcessCFunc(ctx context.Context, inputs NodeInput, config map[string]interface{}) (NodeOutput, bool, error) {
	nodeName := "UnknownProcessCNode"
	if name, ok := config["node_name"].(string); ok {
		nodeName = name
	}
	fmt.Printf("--- [%s] 正在执行处理节点C ---\n", nodeName)
	time.Sleep(300 * time.Millisecond) // 模拟稍长的工作耗时

	var dataFromA string
	for depID, depOutputInterface := range inputs { // 应该只有来自A的输入
		depOutput, ok := depOutputInterface.(NodeOutput)
		if !ok {
			return nil, false, fmt.Errorf("[%s] 严重错误: 来自 %s 的输入不是NodeOutput类型", nodeName, depID)
		}
		if val, valOk := depOutput["data_a"].(string); valOk {
			dataFromA = val
			break
		}
	}

	if dataFromA == "" {
		return nil, false, fmt.Errorf("[%s] 错误: 未从上游接收到data_a", nodeName)
	}

	output := make(NodeOutput)
	output["data_c"] = fmt.Sprintf("C已差异化处理: [%s]", dataFromA) // 另一种处理方式
	fmt.Printf("--- [%s] 输出: %s ---\n", nodeName, output["data_c"])
	return output, false, nil
}

// simpleCombineDFunc 是合并节点D的简单执行函数
func simpleCombineDFunc(ctx context.Context, inputs NodeInput, config map[string]interface{}) (NodeOutput, bool, error) {
	nodeName := "UnknownCombineDNode"
	if name, ok := config["node_name"].(string); ok {
		nodeName = name
	}
	fmt.Printf("--- [%s] 正在执行合并节点D ---\n", nodeName)
	time.Sleep(50 * time.Millisecond) // 模拟工作耗时

	var dataFromB, dataFromC string

	// 从输入中提取来自B和C的数据
	for depID, depOutputInterface := range inputs {
		depOutput, ok := depOutputInterface.(NodeOutput)
		if !ok {
			return nil, false, fmt.Errorf("[%s] 严重错误: 来自 %s 的输入不是NodeOutput类型", nodeName, depID)
		}
		// 假设B的ID是"B"，C的ID是"C"以直接访问，
		// 或检查哪个键存在。
		if val, valOk := depOutput["data_b"].(string); valOk { // 来自B的输出
			dataFromB = val
		}
		if val, valOk := depOutput["data_c"].(string); valOk { // 来自C的输出
			dataFromC = val
		}
	}

	if dataFromB == "" && dataFromC == "" { // 如果B和C都没有提供数据
		// 注意：如果其中一个节点失败并被跳过，这里可能只会收到一个输入。
		// 这里的逻辑需要根据期望的行为调整（例如，如果一个输入缺失是否算错误）。
		// 为了简单起见，我们允许部分数据。
		fmt.Printf("[%s] 警告: 可能未从B或C接收到全部数据 (B: '%s', C: '%s')\n", nodeName, dataFromB, dataFromC)
		// return nil, false, fmt.Errorf("[%s] 错误: 未从B或C接收到数据", nodeName)
	}

	output := make(NodeOutput)
	output["final_result"] = fmt.Sprintf("D已合并: {%s} 与 {%s}", dataFromB, dataFromC) // 合并结果
	fmt.Printf("--- [%s] 输出: %s ---\n", nodeName, output["final_result"])
	return output, false, nil
}

// printResults 打印工作流执行结果的辅助函数
func printResults(ec *ExecutionContext) {
	if ec == nil {
		fmt.Println("执行上下文为nil。")
		return
	}
	fmt.Printf("\n======= 工作流 '%s' 执行结果 =======\n", ec.workflow.Name)
	ec.mu.Lock() // 加锁以安全访问共享的执行结果
	defer ec.mu.Unlock()

	// 按ID对节点排序以便一致地打印输出
	var nodeIDs []string
	for id := range ec.workflow.Nodes {
		nodeIDs = append(nodeIDs, id)
	}
	// 简单的冒泡排序
	for i := 0; i < len(nodeIDs); i++ {
		for j := i + 1; j < len(nodeIDs); j++ {
			if nodeIDs[i] > nodeIDs[j] {
				nodeIDs[i], nodeIDs[j] = nodeIDs[j], nodeIDs[i]
			}
		}
	}

	for _, nodeID := range nodeIDs {
		node := ec.workflow.Nodes[nodeID]
		// 使用节点对象本身的状态，因为它由executeNode更新
		state := node.state   // 或者使用 ec.nodeStates[nodeID]
		output := node.output // 或者使用 ec.nodeResults[nodeID]
		err := node.err       // 或者使用 ec.nodeErrors[nodeID]

		fmt.Printf("  节点: %-20s (ID: %-3s) | 状态: %-10s\n", node.Name, node.ID, state)
		if output != nil {
			fmt.Printf("    输出: %v\n", output)
		}
		if err != nil {
			fmt.Printf("    错误: %v\n", err)
		}
		if node.shortCircuit { // 这指的是节点*自身*决定是否短路（来自ExecuteFunc的返回值）
			fmt.Printf("    显式触发短路: 是\n")
		}
	}
	if ec.shortCircuited {
		fmt.Printf("整个工作流被节点 '%s' 短路。\n", ec.shortCircuitSourceNodeID)
	} else {
		fmt.Printf("整个工作流: 未短路。\n")
	}
	fmt.Println("==============================================")
}

func main() {
	fmt.Println("--- 简单DAG演示 ---")

	// 1. 创建工作流
	wfSimple := NewWorkflow("wfSimple", "简单A->(B,C)->D演示")

	// 2. 定义节点
	nodeA := &Node{
		ID:          "A",                                          // 节点ID
		Name:        "起始节点 A",                                     // 节点名称
		ExecuteFunc: simpleStartFunc,                              // 节点的执行函数
		Config:      map[string]interface{}{"node_name": "NodeA"}, // 节点配置，传递给ExecuteFunc
	}
	nodeB := &Node{
		ID:           "B",
		Name:         "处理节点 B",
		Dependencies: []string{"A"}, // 依赖于节点A
		ExecuteFunc:  simpleProcessBFunc,
		Config:       map[string]interface{}{"node_name": "NodeB"},
	}
	nodeC := &Node{
		ID:           "C",
		Name:         "处理节点 C",
		Dependencies: []string{"A"}, // 依赖于节点A
		ExecuteFunc:  simpleProcessCFunc,
		Config:       map[string]interface{}{"node_name": "NodeC"},
	}
	nodeD := &Node{
		ID:           "D",
		Name:         "合并节点 D",
		Dependencies: []string{"B", "C"}, // 依赖于节点B和C
		ExecuteFunc:  simpleCombineDFunc,
		Config:       map[string]interface{}{"node_name": "NodeD"},
	}

	// 3. 将节点添加到工作流
	wfSimple.AddNode(nodeA)
	wfSimple.AddNode(nodeB)
	wfSimple.AddNode(nodeC)
	wfSimple.AddNode(nodeD)

	// 4. 运行工作流
	// 如果起始节点需要，我们可以传递全局输入
	globalInputs := NodeInput{"global_start_data": "来自Main的问候！"}
	// 为安全起见，可以为整个工作流执行设置超时
	// workflowCtx, cancelWorkflow := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancelWorkflow() // 确保取消函数被调用

	// execCtx, err := wfSimple.Run(workflowCtx, globalInputs)
	// 为这个简单的演示运行不设置显式超时
	execCtx, err := wfSimple.Run(context.Background(), globalInputs)
	if err != nil {
		fmt.Printf("\n简单工作流执行完成并出现错误: %v\n", err)
	} else {
		fmt.Println("\n简单工作流执行完成。")
	}

	// 5. 打印结果
	if execCtx != nil {
		printResults(execCtx)
	}

	fmt.Println("-------------------------------\n")

	// --- 演示：简单DAG中的节点失败 ---
	fmt.Println("--- 简单DAG演示（节点B失败） ---")
	wfFail := NewWorkflow("wfFail", "节点B失败的简单DAG") // 工作流ID和名称
	nodeAF := &Node{ID: "AF", Name: "起始节点 AF", ExecuteFunc: simpleStartFunc, Config: map[string]interface{}{"node_name": "NodeAF"}}
	nodeBF := &Node{ID: "BF", Name: "处理节点 BF (将失败)", Dependencies: []string{"AF"}, Config: map[string]interface{}{"node_name": "NodeBF"},
		ExecuteFunc: func(ctx context.Context, inputs NodeInput, config map[string]interface{}) (NodeOutput, bool, error) {
			fmt.Println("--- [NodeBF] 正在执行 (并将失败)... ---")
			time.Sleep(50 * time.Millisecond)
			return nil, false, errors.New("[NodeBF] 模拟的失败") // 返回错误
		},
	}
	nodeCF := &Node{ID: "CF", Name: "处理节点 CF", Dependencies: []string{"AF"}, ExecuteFunc: simpleProcessCFunc, Config: map[string]interface{}{"node_name": "NodeCF"}}
	nodeDF := &Node{ID: "DF", Name: "合并节点 DF", Dependencies: []string{"BF", "CF"}, ExecuteFunc: simpleCombineDFunc, Config: map[string]interface{}{"node_name": "NodeDF"}}

	wfFail.AddNode(nodeAF)
	wfFail.AddNode(nodeBF)
	wfFail.AddNode(nodeCF)
	wfFail.AddNode(nodeDF)

	execCtxFail, errFail := wfFail.Run(context.Background(), NodeInput{"global_start_data": "用于失败流程的数据"})
	if errFail != nil {
		fmt.Printf("\n失败的工作流按预期执行完成并出现错误: %v\n", errFail) // 预期会看到错误
	}
	if execCtxFail != nil {
		printResults(execCtxFail) // 打印结果，应显示BF失败，DF跳过
	}
	fmt.Println("-------------------------------\n")
}
