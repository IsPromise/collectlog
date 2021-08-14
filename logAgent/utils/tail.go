package utils

import (
	"context"
	"github.com/hpcloud/tail"
	"logsys/logAgent/logs"
	"os"
)

type Task struct {
	Tail  *tail.Tail
	Topic string
}

var (
	tasks []*Task
	//err  error
	allConf []*CollectEntry

	Ctx    context.Context
	Cancel context.CancelFunc

	Flag chan bool
)

// task创建
func newTask(tail *tail.Tail, topic string) *Task {
	return &Task{
		Tail:  tail,
		Topic: topic,
	}
}

func GetTaskNum() int {
	return len(tasks)
}

//task对外接口
func GetTask(i int) *Task {
	return tasks[i]
}

// 使用tail包进行初始化
func InitTail() {
	config := tail.Config{
		ReOpen:    true,
		MustExist: false,
		Poll:      true,
		Follow:    true,
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: os.SEEK_END,
		},
	}
	allConf = getConf(Cfg.CollectKey)

	for _, value := range allConf {
		ta, err := tail.TailFile(value.Path, config)
		if err != nil {
			logs.Warning.Print("tail create failed for path:%s ,err:%v\n", value, err)
			continue
		}
		t := newTask(ta, value.Topic)
		tasks = append(tasks, t)
	}

	Ctx, Cancel = context.WithCancel(context.Background())

	Flag = make(chan bool)

	logs.Info.Println("init tail success")

}

func Notice() {
	// 产生更新通知之前的task goroutine 退出
	Cancel()
	// 任务设置为空
	tasks = nil
	// 重新加载task
	InitTail()
	// 通知重启
	Flag <- true
	logs.Info.Println("config update success")
}
