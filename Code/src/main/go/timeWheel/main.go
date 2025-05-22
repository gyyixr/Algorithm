package main

import (
	"container/list"
	"fmt"
	"time"
)

// 多层时间轮结构
type MultiLevelTimeWheel struct {
	wheels     []*TimeWheel
	levelCount int
}

// 时间轮结构
type TimeWheel struct {
	interval   time.Duration
	ticker     *time.Ticker
	slots      []*list.List
	currentPos int
	slotNum    int
	level      int
	nextWheel  *TimeWheel
}

// 定时任务
type Task struct {
	delay time.Duration
	round int
	key   interface{}
	job   func()
}

// 创建多层时间轮
func NewMultiLevelTimeWheel(baseInterval time.Duration, baseSlotNum int, levelCount int) *MultiLevelTimeWheel {
	mltw := &MultiLevelTimeWheel{
		wheels:     make([]*TimeWheel, levelCount),
		levelCount: levelCount,
	}

	for i := 0; i < levelCount; i++ {
		interval := baseInterval * time.Duration(pow(baseSlotNum, i))
		tw := &TimeWheel{
			interval:   interval,
			slots:      make([]*list.List, baseSlotNum),
			currentPos: 0,
			slotNum:    baseSlotNum,
			level:      i,
		}
		for j := 0; j < baseSlotNum; j++ {
			tw.slots[j] = list.New()
		}
		if i > 0 {
			mltw.wheels[i-1].nextWheel = tw
		}
		mltw.wheels[i] = tw
	}

	return mltw
}

// 启动多层时间轮
func (mltw *MultiLevelTimeWheel) Start() {
	for _, wheel := range mltw.wheels {
		wheel.ticker = time.NewTicker(wheel.interval)
		go wheel.run()
	}
}

// 停止多层时间轮
func (mltw *MultiLevelTimeWheel) Stop() {
	for _, wheel := range mltw.wheels {
		wheel.ticker.Stop()
	}
}

// 添加任务
func (mltw *MultiLevelTimeWheel) AddTask(delay time.Duration, key interface{}, job func()) {
	if delay < 0 {
		return
	}
	mltw.wheels[0].addTask(&Task{delay: delay, key: key, job: job})
}

// 时间轮运行
func (tw *TimeWheel) run() {
	for range tw.ticker.C {
		tw.tickHandler()
	}
}

// 处理一次 tick
func (tw *TimeWheel) tickHandler() {
	l := tw.slots[tw.currentPos]
	tw.scanAndRunTask(l)
	if tw.currentPos == tw.slotNum-1 {
		tw.currentPos = 0
		if tw.nextWheel != nil {
			tw.nextWheel.tickHandler()
		}
	} else {
		tw.currentPos++
	}
}

// 扫描并运行任务
func (tw *TimeWheel) scanAndRunTask(l *list.List) {
	for e := l.Front(); e != nil; {
		task := e.Value.(*Task)
		if task.round > 0 {
			task.round--
			e = e.Next()
			continue
		}

		go task.job()
		next := e.Next()
		l.Remove(e)
		e = next
	}
}

// 添加任务
func (tw *TimeWheel) addTask(task *Task) {
	delay := task.delay
	totalSlots := tw.slotNum * pow(tw.slotNum, tw.level)
	if int(delay.Seconds()) >= totalSlots {
		if tw.nextWheel != nil {
			tw.nextWheel.addTask(task)
		}
		fmt.Printf("delay:%d > tocalSlots:%d\n", int(delay.Seconds()), totalSlots)
		return
	}

	pos, round := tw.getPositionAndRound(delay)
	task.round = round

	tw.slots[pos].PushBack(task)
}

// 获取任务在时间轮中的位置和圈数
func (tw *TimeWheel) getPositionAndRound(d time.Duration) (pos int, round int) {
	delaySeconds := int(d.Seconds())
	intervalSeconds := int(tw.interval.Seconds())
	round = delaySeconds / (intervalSeconds * tw.slotNum)
	pos = (tw.currentPos + delaySeconds/intervalSeconds) % tw.slotNum

	return
}

// 辅助函数：整数幂运算
func pow(base, exp int) int {
	result := 1
	for i := 0; i < exp; i++ {
		result *= base
	}
	return result
}

func main() {
	mltw := NewMultiLevelTimeWheel(1*time.Second, 3, 3)
	mltw.Start()
	defer mltw.Stop()

	mltw.AddTask(5*time.Second, "task1", func() {
		fmt.Println("Task 1 executed")
	})

	mltw.AddTask(28*time.Second, "task2", func() {
		fmt.Println("Task 2 executed")
	})

	time.Sleep(100 * time.Second)
}
