package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

// NodeState 表示节点执行状态
type NodeState int

const (
	Pending NodeState = iota
	Running
	Completed
	Failed
	Skipped // 因短路或上游失败而被跳过
)

func (s NodeState) String() string {
	switch s {
	case Pending:
		return "Pending"
	case Running:
		return "Running"
	case Completed:
		return "Completed"
	case Failed:
		return "Failed"
	case Skipped:
		return "Skipped"
	default:
		return "Unknown"
	}
}

// NodeOutput 是节点执行的输出
type NodeOutput map[string]interface{}

// NodeInput 是节点执行的输入
type NodeInput map[string]interface{}

// Node 定义了工作流中的一个节点
type Node struct {
	ID          string
	Name        string
	Description string
	Config      map[string]interface{} // 节点特定配置

	// ExecuteFunc 是节点的实际执行逻辑
	// ctx: 上下文，可用于超时或取消
	// inputs: 从上游节点获取的输入
	// config: 节点自身的配置
	// 返回: 输出数据，是否短路后续节点，错误
	ExecuteFunc func(ctx context.Context, inputs NodeInput, config map[string]interface{}) (output NodeOutput, shortCircuit bool, err error)

	// Dependencies 列出了该节点依赖的其他节点的ID
	Dependencies []string

	// internal state
	state        NodeState
	output       NodeOutput
	err          error
	shortCircuit bool // 标记此节点是否触发了短路
}

// Workflow 定义了一个工作流
type Workflow struct {
	ID          string
	Name        string
	Nodes       map[string]*Node
	adj         map[string][]string // 邻接表 (id -> []downstream_ids)
	revAdj      map[string][]string // 反向邻接表 (id -> []upstream_ids)
	entryPoints []string            // 入口节点 (没有依赖的节点)
}

// NewWorkflow 创建一个新的工作流
func NewWorkflow(id, name string) *Workflow {
	return &Workflow{
		ID:          id,
		Name:        name,
		Nodes:       make(map[string]*Node),
		adj:         make(map[string][]string),
		revAdj:      make(map[string][]string),
		entryPoints: []string{},
	}
}

// AddNode 向工作流中添加一个节点
func (wf *Workflow) AddNode(node *Node) error {
	if _, exists := wf.Nodes[node.ID]; exists {
		return fmt.Errorf("node with ID %s already exists", node.ID)
	}
	wf.Nodes[node.ID] = node
	wf.adj[node.ID] = []string{}
	wf.revAdj[node.ID] = []string{} // 初始化反向邻接表
	return nil
}

// AddEdge 添加节点间的依赖关系 (from -> to)
func (wf *Workflow) AddEdge(fromID, toID string) error {
	fromNode, fromExists := wf.Nodes[fromID]
	toNode, toExists := wf.Nodes[toID]

	if !fromExists {
		return fmt.Errorf("source node with ID %s does not exist", fromID)
	}
	if !toExists {
		return fmt.Errorf("target node with ID %s does not exist", toID)
	}

	// 更新邻接表
	wf.adj[fromID] = append(wf.adj[fromID], toID)
	// 更新反向邻接表
	wf.revAdj[toID] = append(wf.revAdj[toID], fromID)

	// 更新节点的依赖列表 (虽然我们在Node结构体中有Dependencies，但这里确保图结构正确)
	// 通常，我们会基于Node.Dependencies来构建图，或者反过来。这里我们假设 AddEdge 是权威的。
	found := false
	for _, dep := range toNode.Dependencies {
		if dep == fromID {
			found = true
			break
		}
	}
	if !found {
		toNode.Dependencies = append(toNode.Dependencies, fromID)
	}
	_ = fromNode // use fromNode

	// 检测环 (简化版，实际需要更复杂的DFS)
	// For a full cycle detection, a DFS traversal marking visited and recursion stack nodes is needed.
	// This is a simplified check.
	if wf.hasPath(toID, fromID) {
		// Remove the edge if it creates a cycle
		wf.adj[fromID] = wf.adj[fromID][:len(wf.adj[fromID])-1]
		wf.revAdj[toID] = wf.revAdj[toID][:len(wf.revAdj[toID])-1]
		// also remove from toNode.Dependencies
		var newDeps []string
		for _, dep := range toNode.Dependencies {
			if dep != fromID {
				newDeps = append(newDeps, dep)
			}
		}
		toNode.Dependencies = newDeps
		return fmt.Errorf("adding edge from %s to %s creates a cycle", fromID, toID)
	}

	return nil
}

// hasPath (辅助函数，用于简化环检测，实际应使用DFS)
func (wf *Workflow) hasPath(startNodeID, endNodeID string) bool {
	q := []string{startNodeID}
	visited := make(map[string]bool)
	visited[startNodeID] = true

	for len(q) > 0 {
		curr := q[0]
		q = q[1:]

		if curr == endNodeID {
			return true
		}

		for _, neighbor := range wf.adj[curr] {
			if !visited[neighbor] {
				visited[neighbor] = true
				q = append(q, neighbor)
			}
		}
	}
	return false
}

// Prepare 在执行前准备工作流，例如计算入口节点
func (wf *Workflow) Prepare() error {
	// 重置状态和入口点
	wf.entryPoints = []string{}
	for _, node := range wf.Nodes {
		node.state = Pending
		node.output = nil
		node.err = nil
		node.shortCircuit = false
	}

	// 重新构建邻接表和反向邻接表，基于 Node.Dependencies
	// 清空旧的邻接表信息
	wf.adj = make(map[string][]string)
	wf.revAdj = make(map[string][]string)
	for id := range wf.Nodes {
		wf.adj[id] = []string{}
		wf.revAdj[id] = []string{}
	}

	// 根据 Node.Dependencies 填充邻接表
	for id, node := range wf.Nodes {
		for _, depID := range node.Dependencies {
			if _, exists := wf.Nodes[depID]; !exists {
				return fmt.Errorf("node %s has a dependency %s that does not exist", id, depID)
			}
			wf.adj[depID] = append(wf.adj[depID], id)    // depID -> id
			wf.revAdj[id] = append(wf.revAdj[id], depID) // id depends on depID
		}
	}

	// 查找入口节点
	for id, node := range wf.Nodes {
		if len(node.Dependencies) == 0 {
			wf.entryPoints = append(wf.entryPoints, id)
		}
	}

	if len(wf.entryPoints) == 0 && len(wf.Nodes) > 0 {
		return errors.New("no entry points found in the workflow (possible cycle or all nodes have dependencies)")
	}

	// 简单的环检测：如果准备后没有入口节点且图不为空，则可能有环。
	// 更健壮的环检测应该在 AddEdge 时进行，或者在 Prepare 时进行一次完整的DFS。
	// (已在AddEdge中添加了简化版环检测)

	return nil
}

// ExecutionContext 保存工作流执行期间的上下文
type ExecutionContext struct {
	workflow                 *Workflow
	nodeResults              map[string]NodeOutput // 存储每个节点的输出
	nodeStates               map[string]NodeState  // 存储每个节点的状态
	nodeErrors               map[string]error      // 存储每个节点的错误
	mu                       sync.Mutex            // 保护共享资源
	wg                       sync.WaitGroup
	globalInputs             NodeInput          // 工作流的全局初始输入
	shortCircuitCh           chan string        // 用于通知全局短路事件，值为触发短路的节点ID
	shortCircuited           bool               // 标记整个工作流是否已短路
	shortCircuitSourceNodeID string             // 标记触发短路的节点ID
	ctx                      context.Context    // 主上下文，用于取消
	cancelFunc               context.CancelFunc // 取消函数
}

// NewExecutionContext 创建执行上下文
func NewExecutionContext(wf *Workflow, globalInputs NodeInput, parentCtx context.Context) *ExecutionContext {
	ctx, cancel := context.WithCancel(parentCtx)
	return &ExecutionContext{
		workflow:       wf,
		nodeResults:    make(map[string]NodeOutput),
		nodeStates:     make(map[string]NodeState),
		nodeErrors:     make(map[string]error),
		globalInputs:   globalInputs,
		shortCircuitCh: make(chan string, len(wf.Nodes)), // buffered channel
		ctx:            ctx,
		cancelFunc:     cancel,
	}
}

// Run 执行工作流
func (wf *Workflow) Run(ctx context.Context, globalInputs NodeInput) (*ExecutionContext, error) {
	if err := wf.Prepare(); err != nil {
		return nil, fmt.Errorf("workflow preparation failed: %w", err)
	}

	execCtx := NewExecutionContext(wf, globalInputs, ctx)
	defer execCtx.cancelFunc() // Ensure cancelFunc is called when Run exits

	fmt.Printf("Workflow '%s' starting with %d entry points.\n", wf.Name, len(wf.entryPoints))

	for _, nodeID := range wf.entryPoints {
		execCtx.wg.Add(1)
		go execCtx.executeNode(nodeID)
	}

	// 等待所有goroutine完成或超时/取消
	waitDone := make(chan struct{})
	go func() {
		execCtx.wg.Wait()
		close(waitDone)
	}()

	select {
	case <-waitDone:
		fmt.Printf("Workflow '%s' finished processing all reachable nodes.\n", wf.Name)
	case <-execCtx.ctx.Done(): // If the main context is cancelled
		fmt.Printf("Workflow '%s' cancelled by parent context.\n", wf.Name)
		// Ensure all goroutines are signaled if main context is cancelled
		execCtx.mu.Lock()
		if !execCtx.shortCircuited {
			execCtx.shortCircuited = true
			execCtx.shortCircuitSourceNodeID = "ContextCancelled"
			close(execCtx.shortCircuitCh) // Signal all pending goroutines
		}
		execCtx.mu.Unlock()
		// We might still wait for a short period for goroutines to wind down
		// or rely on their internal context checks.
	}

	// 在这里可以收集最终结果或状态
	// 如果需要，检查是否有未处理的错误
	var encounteredError bool
	for _, err := range execCtx.nodeErrors {
		if err != nil {
			encounteredError = true
			break
		}
	}
	if encounteredError {
		return execCtx, errors.New("workflow execution encountered errors")
	}

	return execCtx, nil
}

// skipDownstreamNodes 递归地将下游节点标记为Skipped
func (ec *ExecutionContext) skipDownstreamNodes(nodeID string, reasonNodeID string) {
	for _, downstreamNodeID := range ec.workflow.adj[nodeID] {
		ec.mu.Lock()
		currentState, ok := ec.nodeStates[downstreamNodeID]
		if !ok || currentState == Pending { // 只有Pending状态的可以被设置为Skipped
			fmt.Printf("Skipping node '%s' because upstream '%s' failed or short-circuited.\n", ec.workflow.Nodes[downstreamNodeID].Name, reasonNodeID)
			ec.nodeStates[downstreamNodeID] = Skipped
			ec.workflow.Nodes[downstreamNodeID].state = Skipped
			ec.mu.Unlock()
			ec.skipDownstreamNodes(downstreamNodeID, reasonNodeID) // 递归跳过更下游的节点
		} else {
			ec.mu.Unlock()
		}
	}
}

func (ec *ExecutionContext) executeNode(nodeID string) {
	defer ec.wg.Done()

	node := ec.workflow.Nodes[nodeID]

	// 检查主上下文是否已取消或是否已发生全局短路
	select {
	case <-ec.ctx.Done():
		ec.mu.Lock()
		if ec.nodeStates[nodeID] != Completed && ec.nodeStates[nodeID] != Failed {
			fmt.Printf("Node '%s' execution cancelled due to workflow context cancellation.\n", node.Name)
			ec.nodeStates[nodeID] = Skipped
			node.state = Skipped
			node.err = ec.ctx.Err()
			ec.nodeErrors[nodeID] = node.err
		}
		ec.mu.Unlock()
		return
	case scNodeID, ok := <-ec.shortCircuitCh:
		if ok { // Channel is still open, means a short circuit happened
			ec.mu.Lock()
			// Check if we should globally short circuit
			if !ec.shortCircuited {
				ec.shortCircuited = true
				ec.shortCircuitSourceNodeID = scNodeID
				fmt.Printf("Workflow short-circuit triggered by node '%s'. Propagating...\n", scNodeID)
				// Close channel to signal all other listeners
				// Be careful: only close it once. This is why shortCircuited flag is important.
				// This is tricky because multiple goroutines might read. A better way is a dedicated closer.
				// For simplicity, we'll allow multiple reads but only one effective signal.
				// No, closing here is bad if multiple goroutines read from it.
				// Instead, rely on the shortCircuited flag and parent context cancellation.
			}
			ec.mu.Unlock()
		} else { // Channel closed, means short circuit was already signaled and propagated
			ec.mu.Lock()
			ec.shortCircuited = true // Ensure flag is set
			ec.mu.Unlock()
		}
	default:
		// Continue if no immediate signal
	}

	ec.mu.Lock()
	// 如果整个工作流已经短路，并且当前节点不是触发短路的节点，则跳过
	if ec.shortCircuited && ec.shortCircuitSourceNodeID != nodeID {
		// 确保只有Pending的节点会被标记为Skipped
		if ec.nodeStates[nodeID] == Pending {
			fmt.Printf("Node '%s' skipped due to global workflow short-circuit triggered by '%s'.\n", node.Name, ec.shortCircuitSourceNodeID)
			ec.nodeStates[nodeID] = Skipped
			node.state = Skipped
		}
		ec.mu.Unlock()
		return
	}

	// 检查依赖是否满足
	inputs := make(NodeInput)
	allDepsCompleted := true
	anyDepFailedOrSkipped := false
	reasonForSkip := ""

	for _, depID := range node.Dependencies {
		depState, stateExists := ec.nodeStates[depID]
		depOutput, outputExists := ec.nodeResults[depID]

		if !stateExists || depState == Pending || depState == Running {
			allDepsCompleted = false // 依赖项尚未完成
			break
		}
		if depState == Failed || depState == Skipped {
			anyDepFailedOrSkipped = true
			reasonForSkip = depID
			break
		}
		if depState == Completed {
			if !outputExists {
				// 这不应该发生，如果节点完成，应该有输出
				fmt.Printf("Error: Dependency '%s' for node '%s' is completed but has no output.\n", depID, node.Name)
				allDepsCompleted = false // Treat as incomplete
				anyDepFailedOrSkipped = true
				reasonForSkip = depID
				node.err = fmt.Errorf("dependency %s completed without output", depID)
				ec.nodeErrors[nodeID] = node.err
				break
			}
			// 合并依赖项的输出作为当前节点的输入
			// 简单的合并：假设依赖的输出名不冲突，或者按约定处理
			// Dify等系统会有更复杂的输入映射机制
			inputs[depID] = depOutput // 例如: inputs["nodeA_id"] = nodeA_output
			// 或者更细致地： inputs["input_name_for_current_node"] = depOutput["specific_output_key_from_dep_node"]
		}
	}
	ec.mu.Unlock() // 释放锁，允许其他节点更新状态

	if !allDepsCompleted {
		// 理论上，这个分支不应该被频繁命中，因为节点是由其依赖项完成时触发的
		// 但如果并发模型允许直接调度，或者有竞争条件，则需要处理
		fmt.Printf("Node '%s' cannot run yet, dependencies not met.\n", node.Name)
		// ec.wg.Done() // This wg.Done was already deferred. No need to call it here.
		return // Should be re-queued or handled by dependency trigger
	}

	if anyDepFailedOrSkipped {
		ec.mu.Lock()
		fmt.Printf("Node '%s' skipped because its dependency '%s' failed or was skipped.\n", node.Name, ec.workflow.Nodes[reasonForSkip].Name)
		ec.nodeStates[nodeID] = Skipped
		node.state = Skipped
		ec.mu.Unlock()
		// 递归跳过下游节点
		ec.skipDownstreamNodes(nodeID, reasonForSkip)
		return
	}

	// 如果是入口节点且有全局输入，则合并
	if len(node.Dependencies) == 0 && ec.globalInputs != nil {
		for k, v := range ec.globalInputs {
			inputs[k] = v // 全局输入可以被节点ID的输出覆盖（如果命名冲突）
		}
	}

	// --- 开始执行节点 ---
	fmt.Printf("Node '%s' (ID: %s) starting execution.\n", node.Name, node.ID)
	ec.mu.Lock()
	ec.nodeStates[nodeID] = Running
	node.state = Running
	ec.mu.Unlock()

	// 执行节点的函数
	// 为此节点执行创建一个带超时的上下文 (示例: 30秒)
	// nodeCtx, nodeCancel := context.WithTimeout(ec.ctx, 30*time.Second)
	// defer nodeCancel()
	// 使用ec.ctx的子上下文，这样如果主流程取消，节点也能感知
	nodeOutput, shouldShortCircuit, err := node.ExecuteFunc(ec.ctx, inputs, node.Config)

	ec.mu.Lock()
	defer ec.mu.Unlock() // Ensure lock is released at the end of this block

	// 检查执行后全局是否已短路 (可能由其他并行节点触发)
	if ec.shortCircuited && ec.shortCircuitSourceNodeID != nodeID {
		if node.state != Completed && node.state != Failed { // 避免覆盖已完成/失败的状态
			fmt.Printf("Node '%s' execution was successful/failed, but workflow already short-circuited by '%s'. Marking as Skipped/Failed.\n", node.Name, ec.shortCircuitSourceNodeID)
			if err == nil { // 如果此节点本身没出错，但工作流短路了，标记为 Skipped
				ec.nodeStates[nodeID] = Skipped
				node.state = Skipped
			} else { // 如果此节点也出错了
				ec.nodeStates[nodeID] = Failed
				node.state = Failed
				node.err = err
				ec.nodeErrors[nodeID] = err
			}
		}
		// 不再触发下游，因为工作流已短路
		return
	}

	if err != nil {
		fmt.Printf("Node '%s' (ID: %s) execution failed: %v\n", node.Name, node.ID, err)
		ec.nodeStates[nodeID] = Failed
		node.state = Failed
		node.err = err
		ec.nodeErrors[nodeID] = err
		// 如果此节点失败，也认为是一种短路，影响下游
		node.shortCircuit = true // 标记此节点是短路源（即使shouldShortCircuit为false）
	} else {
		fmt.Printf("Node '%s' (ID: %s) execution completed.\n", node.Name, node.ID)
		ec.nodeStates[nodeID] = Completed
		node.state = Completed
		ec.nodeResults[nodeID] = nodeOutput
		node.output = nodeOutput
		node.shortCircuit = shouldShortCircuit // 记录节点是否希望短路
	}

	// 处理短路逻辑
	if node.shortCircuit { // 如果当前节点（成功或失败后）要求短路
		if !ec.shortCircuited { // 只有在全局尚未短路时才触发
			ec.shortCircuited = true
			ec.shortCircuitSourceNodeID = nodeID
			fmt.Printf("Node '%s' triggered short-circuit for the workflow.\n", node.Name)
			// Signal other goroutines. Closing channel is a broadcast.
			// Make sure it's closed only once. This critical section must be robust.
			// One way to ensure single close: use a sync.Once or another flag.
			// For this example, we rely on the ec.shortCircuited flag before closing.
			close(ec.shortCircuitCh) // Signal global short-circuit
			ec.cancelFunc()          // Также отменяем главный контекст, чтобы остановить длительные задачи
		}
	}

	// 如果当前节点失败或触发短路，则其下游节点应标记为Skipped
	if node.state == Failed || (node.state == Completed && node.shortCircuit) {
		reasonForSkip := node.ID
		if node.state == Failed {
			fmt.Printf("Node '%s' failed. Skipping downstream nodes.\n", node.Name)
		} else if node.shortCircuit {
			fmt.Printf("Node '%s' completed and triggered short-circuit. Skipping downstream nodes.\n", node.Name)
		}
		// 必须在锁外调用，以避免死锁，因为skipDownstreamNodes也会获取锁
		// 但要确保状态已更新
		// Unlock before recursive call, and relock if necessary, or pass necessary info.
		// The current structure: lock, update state, unlock. Then call dependent logic.
		// This seems fine as skipDownstreamNodes operates on its own logic based on current states.

		// 收集下游节点ID
		downstreamIDs := ec.workflow.adj[nodeID]
		// 释放当前锁，因为skipDownstreamNodes会尝试获取锁
		ec.mu.Unlock()
		for _, dsNodeID := range downstreamIDs {
			ec.skipDownstreamNodeIfParentTriggered(dsNodeID, reasonForSkip)
		}
		ec.mu.Lock() // 重新获取锁以完成此函数的其余部分（如果有的话）

	} else if node.state == Completed && !node.shortCircuit { // 如果成功且不短路，则触发下游节点
		// 收集下游节点ID
		downstreamIDs := ec.workflow.adj[nodeID]
		// 释放当前锁，因为executeNode会尝试获取锁
		ec.mu.Unlock()
		for _, nextNodeID := range downstreamIDs {
			// 检查下游节点的所有依赖是否都已完成
			if ec.areDependenciesMetAndNotSkipped(nextNodeID) {
				ec.wg.Add(1)
				go ec.executeNode(nextNodeID)
			}
		}
		ec.mu.Lock() // 重新获取锁
	}
}

// skipDownstreamNodeIfParentTriggered is a helper to skip a single node and its children
// This is called when a parent node has failed or short-circuited.
func (ec *ExecutionContext) skipDownstreamNodeIfParentTriggered(nodeID string, reasonNodeID string) {
	ec.mu.Lock()
	currentState, ok := ec.nodeStates[nodeID]
	// Only skip if pending. If it's already running, completed, failed, or skipped, don't change state.
	if ok && currentState != Pending {
		ec.mu.Unlock()
		return
	}
	if !ok { // Not yet in states map, implies Pending
		ec.nodeStates[nodeID] = Pending // Ensure it's in the map
	}

	fmt.Printf("Skipping node '%s' because its direct upstream '%s' failed or short-circuited.\n", ec.workflow.Nodes[nodeID].Name, ec.workflow.Nodes[reasonNodeID].Name)
	ec.nodeStates[nodeID] = Skipped
	ec.workflow.Nodes[nodeID].state = Skipped
	ec.mu.Unlock()

	// Recursively skip its children
	downstreamChildren := ec.workflow.adj[nodeID]
	for _, childNodeID := range downstreamChildren {
		ec.skipDownstreamNodeIfParentTriggered(childNodeID, nodeID) // Reason is now the current skipped node
	}
}

// areDependenciesMetAndNotSkipped 检查节点的所有依赖是否都已完成，并且没有一个依赖是Failed或Skipped
func (ec *ExecutionContext) areDependenciesMetAndNotSkipped(nodeID string) bool {
	ec.mu.Lock() // Lock to safely read shared states
	defer ec.mu.Unlock()

	node := ec.workflow.Nodes[nodeID]
	if node == nil {
		return false // Node not found
	}

	for _, depID := range node.Dependencies {
		depState, stateExists := ec.nodeStates[depID]
		if !stateExists || (depState != Completed) {
			// If any dependency is not completed (i.e., Pending, Running, Failed, Skipped),
			// then dependencies are not met for successful execution.
			return false
		}
	}
	return true // All dependencies are 'Completed'
}

// --- 示例节点执行函数 ---

// StartNodeFunc 简单的起始节点，可以从全局输入获取数据
func StartNodeFunc(ctx context.Context, inputs NodeInput, config map[string]interface{}) (NodeOutput, bool, error) {
	fmt.Printf("Executing Start Node (config: %v)\n", config)
	output := make(NodeOutput)
	if initialData, ok := inputs["initial_data"]; ok {
		output["data"] = fmt.Sprintf("Started with: %v", initialData)
	} else {
		output["data"] = "Started with no specific data."
	}
	time.Sleep(1 * time.Second) // 模拟工作
	return output, false, nil
}

// ProcessingNodeFunc 模拟数据处理
func ProcessingNodeFunc(ctx context.Context, inputs NodeInput, config map[string]interface{}) (NodeOutput, bool, error) {
	fmt.Printf("Executing Processing Node (config: %v)\n", config)
	output := make(NodeOutput)
	var combinedInputData []string
	for depID, depOutputInterface := range inputs {
		depOutput := depOutputInterface.(NodeOutput) // TYPE ASSERTION
		if data, ok := depOutput["data"].(string); ok {
			combinedInputData = append(combinedInputData, fmt.Sprintf("from %s: %s", depID, data))
		}
	}

	processedData := strings.Join(combinedInputData, " | ")
	output["data"] = fmt.Sprintf("Processed: [%s]", processedData)

	// 模拟可能发生的错误
	if failRate, ok := config["fail_rate"].(float64); ok {
		if time.Now().UnixNano()%int64(1.0/failRate) == 0 { // 简单的随机失败
			return nil, false, errors.New("simulated processing error")
		}
	}

	// 模拟短路条件
	if shortCircuitVal, ok := config["short_circuit_if_contains"].(string); ok {
		if strings.Contains(processedData, shortCircuitVal) {
			fmt.Printf("Processing Node is triggering short circuit because data contains '%s'\n", shortCircuitVal)
			output["reason"] = "Contains " + shortCircuitVal
			return output, true, nil // 短路，但没有错误
		}
	}

	time.Sleep(2 * time.Second) // 模拟工作
	return output, false, nil
}

// ConditionalNodeFunc 模拟条件判断，可能会触发短路
func ConditionalNodeFunc(ctx context.Context, inputs NodeInput, config map[string]interface{}) (NodeOutput, bool, error) {
	fmt.Printf("Executing Conditional Node (config: %v)\n", config)
	output := make(NodeOutput)
	var inputData string
	for _, depOutputInterface := range inputs { // 取第一个依赖的输出
		depOutput := depOutputInterface.(NodeOutput) // TYPE ASSERTION
		if data, ok := depOutput["data"].(string); ok {
			inputData = data
			break
		}
	}

	conditionField, _ := config["condition_field"].(string) // e.g., "data"
	mustContain, _ := config["must_contain"].(string)       // e.g., "secret"

	if strings.Contains(inputData, mustContain) {
		output["message"] = fmt.Sprintf("Condition met: '%s' contains '%s'.", conditionField, mustContain)
		output["data"] = inputData // 透传数据
		fmt.Printf("Conditional Node: condition met for '%s'.\n", inputData)
		time.Sleep(500 * time.Millisecond)
		return output, false, nil
	}

	errMsg := fmt.Sprintf("Condition NOT met: input '%s' does not contain '%s'. Triggering short-circuit.", inputData, mustContain)
	output["message"] = errMsg
	output["data"] = inputData
	fmt.Println(errMsg)
	// return output, true, errors.New(errMsg) // 既短路又报错
	return output, true, nil // 短路，但不报错
}

// FinalNodeFunc 最终节点，收集结果
func FinalNodeFunc(ctx context.Context, inputs NodeInput, config map[string]interface{}) (NodeOutput, bool, error) {
	fmt.Printf("Executing Final Node (config: %v)\n", config)
	output := make(NodeOutput)
	var finalResults []string
	for depID, depOutputInterface := range inputs {
		depOutput := depOutputInterface.(NodeOutput) // TYPE ASSERTION
		if data, ok := depOutput["data"]; ok {
			finalResults = append(finalResults, fmt.Sprintf("From %s: %v", depID, data))
		}
		if msg, ok := depOutput["message"]; ok {
			finalResults = append(finalResults, fmt.Sprintf("Message from %s: %v", depID, msg))
		}
	}
	output["summary"] = strings.Join(finalResults, "; ")
	time.Sleep(1 * time.Second)
	return output, false, nil
}

func main() {
	// --- Demo 1: 基本的线性流程 ---
	fmt.Println("--- Demo 1: Linear Workflow ---")
	wf1 := NewWorkflow("wf1", "Linear Demo")
	nodeA1 := &Node{ID: "A1", Name: "Start A1", ExecuteFunc: StartNodeFunc}
	nodeB1 := &Node{ID: "B1", Name: "Process B1", Dependencies: []string{"A1"}, ExecuteFunc: ProcessingNodeFunc}
	nodeC1 := &Node{ID: "C1", Name: "End C1", Dependencies: []string{"B1"}, ExecuteFunc: FinalNodeFunc}
	wf1.AddNode(nodeA1)
	wf1.AddNode(nodeB1)
	wf1.AddNode(nodeC1)

	initialInputs1 := NodeInput{"initial_data": "Hello Linear World"}
	execCtx1, err1 := wf1.Run(context.Background(), initialInputs1)
	if err1 != nil {
		fmt.Printf("Workflow 1 execution error: %v\n", err1)
	}
	printResults(execCtx1)
	fmt.Println("-------------------------------\n")
	time.Sleep(1 * time.Second)

	// --- Demo 2: 并发与合并 ---
	fmt.Println("--- Demo 2: Concurrent Execution & Merge ---")
	wf2 := NewWorkflow("wf2", "Concurrency Demo")
	nodeA2 := &Node{ID: "A2", Name: "Start A2", ExecuteFunc: StartNodeFunc}
	nodeB2 := &Node{ID: "B2", Name: "Process B2 (fast)", Dependencies: []string{"A2"}, ExecuteFunc: ProcessingNodeFunc, Config: map[string]interface{}{"task_name": "B2_fast"}}
	nodeC2 := &Node{ID: "C2", Name: "Process C2 (slow)", Dependencies: []string{"A2"}, ExecuteFunc: ProcessingNodeFunc, Config: map[string]interface{}{"task_name": "C2_slow"}}
	nodeD2 := &Node{ID: "D2", Name: "Final D2", Dependencies: []string{"B2", "C2"}, ExecuteFunc: FinalNodeFunc}
	wf2.AddNode(nodeA2)
	wf2.AddNode(nodeB2)
	wf2.AddNode(nodeC2)
	wf2.AddNode(nodeD2)

	initialInputs2 := NodeInput{"initial_data": "Hello Concurrent World"}
	execCtx2, err2 := wf2.Run(context.Background(), initialInputs2)
	if err2 != nil {
		fmt.Printf("Workflow 2 execution error: %v\n", err2)
	}
	printResults(execCtx2)
	fmt.Println("-------------------------------\n")
	time.Sleep(1 * time.Second)

	// --- Demo 3: 节点失败导致下游跳过 ---
	fmt.Println("--- Demo 3: Node Failure and Skipping ---")
	wf3 := NewWorkflow("wf3", "Failure Demo")
	nodeA3 := &Node{ID: "A3", Name: "Start A3", ExecuteFunc: StartNodeFunc}
	nodeB3 := &Node{ID: "B3", Name: "Process B3 (will fail)", Dependencies: []string{"A3"}, ExecuteFunc: ProcessingNodeFunc, Config: map[string]interface{}{"fail_rate": 1.0}} // 100% fail
	nodeC3 := &Node{ID: "C3", Name: "Process C3 (parallel to B3)", Dependencies: []string{"A3"}, ExecuteFunc: ProcessingNodeFunc, Config: map[string]interface{}{"task_name": "C3_ok"}}
	nodeD3 := &Node{ID: "D3", Name: "Final D3 (depends on B3, C3)", Dependencies: []string{"B3", "C3"}, ExecuteFunc: FinalNodeFunc}
	wf3.AddNode(nodeA3)
	wf3.AddNode(nodeB3)
	wf3.AddNode(nodeC3)
	wf3.AddNode(nodeD3)

	initialInputs3 := NodeInput{"initial_data": "Hello Failure World"}
	execCtx3, err3 := wf3.Run(context.Background(), initialInputs3)
	if err3 != nil {
		fmt.Printf("Workflow 3 execution finished with error indication: %v\n", err3) // Expected
	}
	printResults(execCtx3)
	fmt.Println("-------------------------------\n")
	time.Sleep(1 * time.Second)

	// --- Demo 4: 条件节点触发短路 ---
	fmt.Println("--- Demo 4: Conditional Short-Circuit ---")
	wf4 := NewWorkflow("wf4", "Short-Circuit Demo")
	nodeA4 := &Node{ID: "A4", Name: "Start A4", ExecuteFunc: StartNodeFunc}
	nodeB4 := &Node{ID: "B4", Name: "Prepare Data B4", Dependencies: []string{"A4"}, ExecuteFunc: ProcessingNodeFunc, Config: map[string]interface{}{"content": "some data with secret_token"}} // Let's ensure StartNode passes something that ProcessingNode can use. Modify StartNode or provide specific initial input.
	nodeC4 := &Node{ID: "C4", Name: "Conditional C4 (will short-circuit)", Dependencies: []string{"B4"}, ExecuteFunc: ConditionalNodeFunc, Config: map[string]interface{}{"condition_field": "data", "must_contain": "no_such_token"}}
	nodeD4 := &Node{ID: "D4", Name: "Process D4 (should be skipped)", Dependencies: []string{"C4"}, ExecuteFunc: ProcessingNodeFunc}
	nodeE4 := &Node{ID: "E4", Name: "Parallel Process E4 (should run if A4 works)", Dependencies: []string{"A4"}, ExecuteFunc: ProcessingNodeFunc, Config: map[string]interface{}{"task_name": "E4_parallel_independent"}}
	nodeF4 := &Node{ID: "F4", Name: "Final F4 (depends on D4, E4)", Dependencies: []string{"D4", "E4"}, ExecuteFunc: FinalNodeFunc}

	wf4.AddNode(nodeA4)
	wf4.AddNode(nodeB4)
	wf4.AddNode(nodeC4)
	wf4.AddNode(nodeD4)
	wf4.AddNode(nodeE4)
	wf4.AddNode(nodeF4)

	// For B4 to have meaningful data to pass to C4, StartNodeFunc's output "data" key will be used by ProcessingNodeFunc.
	// ConditionalNodeFunc will look for inputs[any_dep_id]["data"].
	// Ensure B4 receives "initial_data" from A4 and processes it.
	// StartNodeFunc provides `output["data"] = "Started with: %v", initialData`
	// ProcessingNodeFunc (B4) will receive `inputs["A4"]["data"]` and create `output["data"] = "Processed: [from A4: Started with: ... ]"`
	// ConditionalNodeFunc (C4) will receive `inputs["B4"]["data"]` and check its content.
	// We set "must_contain": "no_such_token" for C4, so it will trigger short-circuit.
	initialInputs4 := NodeInput{"initial_data": "confidential project info"}
	execCtx4, err4 := wf4.Run(context.Background(), initialInputs4)
	if err4 != nil {
		fmt.Printf("Workflow 4 execution finished with error indication: %v\n", err4)
	}
	printResults(execCtx4)
	fmt.Println("-------------------------------\n")
	time.Sleep(1 * time.Second)

	// --- Demo 5: 短路条件满足，不短路后续 ---
	fmt.Println("--- Demo 5: Short-Circuit Condition Met, but node continues (no short-circuit signal) ---")
	wf5 := NewWorkflow("wf5", "Conditional No-Short-Circuit Demo")
	nodeA5 := &Node{ID: "A5", Name: "Start A5", ExecuteFunc: StartNodeFunc}
	// B5 will process data from A5.
	nodeB5 := &Node{ID: "B5", Name: "Prepare Data B5", Dependencies: []string{"A5"}, ExecuteFunc: ProcessingNodeFunc, Config: map[string]interface{}{"content": "some data"}}
	// C5 will find the token, so it will NOT short-circuit.
	nodeC5 := &Node{ID: "C5", Name: "Conditional C5 (no short-circuit)", Dependencies: []string{"B5"}, ExecuteFunc: ConditionalNodeFunc, Config: map[string]interface{}{"condition_field": "data", "must_contain": "Processed"}} // "Processed" is likely in B5's output
	nodeD5 := &Node{ID: "D5", Name: "Process D5 (should run)", Dependencies: []string{"C5"}, ExecuteFunc: ProcessingNodeFunc}
	nodeE5 := &Node{ID: "E5", Name: "Final E5", Dependencies: []string{"D5"}, ExecuteFunc: FinalNodeFunc}

	wf5.AddNode(nodeA5)
	wf5.AddNode(nodeB5)
	wf5.AddNode(nodeC5)
	wf5.AddNode(nodeD5)
	wf5.AddNode(nodeE5)

	initialInputs5 := NodeInput{"initial_data": "workflow data"}
	execCtx5, err5 := wf5.Run(context.Background(), initialInputs5)
	if err5 != nil {
		fmt.Printf("Workflow 5 execution finished with error indication: %v\n", err5)
	}
	printResults(execCtx5)
	fmt.Println("-------------------------------\n")
}

func printResults(ec *ExecutionContext) {
	if ec == nil {
		fmt.Println("Execution context is nil.")
		return
	}
	fmt.Printf("Workflow '%s' Execution Results:\n", ec.workflow.Name)
	ec.mu.Lock()
	defer ec.mu.Unlock()

	// Sort nodes by ID for consistent printing order
	var nodeIDs []string
	for id := range ec.workflow.Nodes {
		nodeIDs = append(nodeIDs, id)
	}
	// Simple sort, can be more sophisticated if needed
	for i := 0; i < len(nodeIDs); i++ {
		for j := i + 1; j < len(nodeIDs); j++ {
			if nodeIDs[i] > nodeIDs[j] {
				nodeIDs[i], nodeIDs[j] = nodeIDs[j], nodeIDs[i]
			}
		}
	}

	for _, nodeID := range nodeIDs {
		node := ec.workflow.Nodes[nodeID]
		state := ec.nodeStates[nodeID] // Get the most up-to-date state from execution context
		output := ec.nodeResults[nodeID]
		err := ec.nodeErrors[nodeID]

		fmt.Printf("  Node: %s (ID: %s)\n", node.Name, node.ID)
		fmt.Printf("    State: %s\n", state)
		if output != nil {
			fmt.Printf("    Output: %v\n", output)
		}
		if err != nil {
			fmt.Printf("    Error: %v\n", err)
		}
		if node.shortCircuit { // This refers to the node's *own* decision to short-circuit
			fmt.Printf("    Triggered Short-Circuit: Yes\n")
		}
	}
	if ec.shortCircuited {
		fmt.Printf("Overall Workflow Short-Circuited by Node: %s\n", ec.shortCircuitSourceNodeID)
	} else {
		fmt.Printf("Overall Workflow: Not short-circuited.\n")
	}
}
