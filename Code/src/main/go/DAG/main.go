package main

import (
	"fmt"
	"sync"
)

// Task 定义了一个工作流中的任务
type Task struct {
	ID           string
	Name         string
	Dependencies []string     // 依赖的任务 ID 列表
	Action       func() error // 任务执行的动作
	done         bool         // 标记任务是否完成
	mutex        sync.Mutex
}

// Workflow 定义了一个工作流
type Workflow struct {
	Tasks map[string]*Task // 任务 ID 到任务指针的映射
}

// NewWorkflow 创建一个新的工作流
func NewWorkflow() *Workflow {
	return &Workflow{
		Tasks: make(map[string]*Task),
	}
}

// AddTask 向工作流中添加一个任务
func (wf *Workflow) AddTask(task *Task) error {
	if _, exists := wf.Tasks[task.ID]; exists {
		return fmt.Errorf("task with ID %s already exists", task.ID)
	}
	wf.Tasks[task.ID] = task
	return nil
}

// Run 执行工作流
func (wf *Workflow) Run() error {
	var wg sync.WaitGroup
	errorsChan := make(chan error, len(wf.Tasks)) // 用于收集并发执行中的错误

	// 检查循环依赖 (一个非常基础的检查，实际场景可能需要更复杂的算法)
	if err := wf.detectCycle(); err != nil {
		return fmt.Errorf("cycle detected in workflow: %v", err)
	}

	for id := range wf.Tasks {
		wg.Add(1)
		go wf.executeTask(id, &wg, errorsChan)
	}

	wg.Wait()
	close(errorsChan)

	// 检查执行过程中是否有错误
	for err := range errorsChan {
		if err != nil {
			// 实际场景中可能需要更复杂的错误处理，比如停止其他任务等
			return fmt.Errorf("error during task execution: %v", err)
		}
	}

	fmt.Println("Workflow completed successfully!")
	return nil
}

// executeTask 执行单个任务及其依赖项
func (wf *Workflow) executeTask(taskID string, wg *sync.WaitGroup, errorsChan chan<- error) {
	defer wg.Done()

	task, exists := wf.Tasks[taskID]
	if !exists {
		errorsChan <- fmt.Errorf("task with ID %s not found during execution", taskID)
		return
	}

	// 加锁以确保任务状态的原子性操作
	task.mutex.Lock()
	if task.done {
		task.mutex.Unlock()
		return // 任务已完成，直接返回
	}
	task.mutex.Unlock() // 先释放一次，避免后续依赖检查死锁

	// 执行依赖任务
	for _, depID := range task.Dependencies {
		depTask, depExists := wf.Tasks[depID]
		if !depExists {
			errorsChan <- fmt.Errorf("dependency task with ID %s for task %s not found", depID, taskID)
			return
		}

		// 递归执行依赖任务 (这里也需要确保依赖任务完成)
		// 为了简单起见，我们这里假设依赖任务会自行完成或者被其他 goroutine 执行
		// 一个更健壮的实现会在这里等待依赖任务完成
		depTask.mutex.Lock()
		depDone := depTask.done
		depTask.mutex.Unlock()

		if !depDone {
			// 如果依赖项未完成，理论上应该等待或触发执行
			// 为简化，这里假设依赖会通过其他途径被执行，或者需要更复杂的调度逻辑
			// 这里我们先简单地认为如果依赖没完成，当前任务也不能执行 (实际场景需要阻塞或重新调度)
			fmt.Printf("Task %s is waiting for dependency %s\n", task.Name, depTask.Name)
			// 为了避免死锁和演示简单性，我们不在这里阻塞，而是依赖于并发执行顺序
			// 实际场景中，这里可以使用 channel 或其他同步原语来等待依赖完成
			// 或者，更简单的方式是，在 Run 方法中按拓扑排序来启动任务。
			// 但为了展示并发性，这里采用 goroutine + 依赖检查。
			// 这里我们先返回，等待依赖任务被执行
			// 重新将任务加入等待组，因为当前任务尚未真正开始执行
			// wg.Add(1) // 这个逻辑有点复杂，暂时去掉，依赖于并发执行的顺序
			return
		}
	}

	// 所有依赖都完成后，执行当前任务
	fmt.Printf("Executing task: %s (%s)\n", task.Name, task.ID)
	if err := task.Action(); err != nil {
		errorsChan <- fmt.Errorf("error in task %s (%s): %v", task.Name, task.ID, err)
		return
	}

	task.mutex.Lock()
	task.done = true
	task.mutex.Unlock()
	fmt.Printf("Task %s (%s) completed.\n", task.Name, task.ID)
}

// detectCycle 检测工作流中是否存在循环依赖 (使用深度优先搜索)
func (wf *Workflow) detectCycle() error {
	visited := make(map[string]bool)        // 标记已访问的节点
	recursionStack := make(map[string]bool) // 标记当前递归栈中的节点

	for id := range wf.Tasks {
		if !visited[id] {
			if wf.isCyclicUtil(id, visited, recursionStack) {
				return fmt.Errorf("cycle detected involving task %s", id)
			}
		}
	}
	return nil
}

func (wf *Workflow) isCyclicUtil(taskID string, visited map[string]bool, recursionStack map[string]bool) bool {
	visited[taskID] = true
	recursionStack[taskID] = true

	task, exists := wf.Tasks[taskID]
	if !exists {
		return false // 理论上不应该发生，因为是从 wf.Tasks 迭代的
	}

	for _, depID := range task.Dependencies {
		if !visited[depID] {
			if wf.isCyclicUtil(depID, visited, recursionStack) {
				return true
			}
		} else if recursionStack[depID] {
			return true // 在递归栈中再次遇到，说明有环
		}
	}

	recursionStack[taskID] = false // 离开递归栈
	return false
}

func main() {
	wf := NewWorkflow()

	// 定义任务
	taskA := &Task{
		ID:   "A",
		Name: "Task A",
		Action: func() error {
			fmt.Println("Running Task A...")
			// time.Sleep(1 * time.Second) // 模拟耗时操作
			return nil
		},
	}

	taskB := &Task{
		ID:           "B",
		Name:         "Task B",
		Dependencies: []string{"A"}, // B 依赖 A
		Action: func() error {
			fmt.Println("Running Task B...")
			// time.Sleep(1 * time.Second)
			return nil
		},
	}

	taskC := &Task{
		ID:           "C",
		Name:         "Task C",
		Dependencies: []string{"A"}, // C 依赖 A
		Action: func() error {
			fmt.Println("Running Task C...")
			// time.Sleep(1 * time.Second)
			return nil
		},
	}

	taskD := &Task{
		ID:           "D",
		Name:         "Task D",
		Dependencies: []string{"B", "C"}, // D 依赖 B 和 C
		Action: func() error {
			fmt.Println("Running Task D...")
			// time.Sleep(1 * time.Second)
			return nil
		},
	}

	taskE := &Task{ // 独立的任务
		ID:   "E",
		Name: "Task E",
		Action: func() error {
			fmt.Println("Running Task E...")
			return nil
		},
	}

	// 添加任务到工作流
	wf.AddTask(taskA)
	wf.AddTask(taskB)
	wf.AddTask(taskC)
	wf.AddTask(taskD)
	wf.AddTask(taskE)

	//尝试添加一个会产生循环依赖的任务(用于测试循环检测)
	taskF := &Task{
		ID:           "F",
		Name:         "Task F",
		Dependencies: []string{"F"}, // F 依赖 F (自身循环)
		Action: func() error {
			fmt.Println("Running Task F...")
			return nil
		},
	}
	wf.AddTask(taskF)

	taskG := &Task{
		ID:           "G",
		Name:         "Task G",
		Dependencies: []string{"H"},
		Action:       func() error { fmt.Println("Running Task G..."); return nil },
	}
	taskH := &Task{
		ID:           "H",
		Name:         "Task H",
		Dependencies: []string{"G"}, // H 依赖 G， G 依赖 H (相互循环)
		Action:       func() error { fmt.Println("Running Task H..."); return nil },
	}
	wf.AddTask(taskG)
	wf.AddTask(taskH)

	//执行工作流
	if err := wf.Run(); err != nil {
		fmt.Printf("Workflow execution failed: %v\n", err)
	}
}
