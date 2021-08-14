package agent

import (
	"context"
	"github.com/Shopify/sarama"
	"logsys/logAgent/kafka"
	"logsys/logAgent/logs"
	"logsys/logAgent/utils"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

// tail-->log-->client-->kafka的业务逻辑
func Run() {
	for i := 0; i < utils.GetTaskNum(); i++ {
		wg.Add(1)
		go Task(utils.GetTask(i), utils.Ctx)
	}
	wg.Wait()
}

// tail中拉去出日志 构造msg 调用发送接口
func Task(task *utils.Task, ctx context.Context) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			// 关闭 tail 对象 和 子goroutine
			task.Tail.Cleanup()
			logs.Trace.Println("goroutine to be end")
			return
		default:
			// 从tail中拉取出日志
			line, ok := <-task.Tail.Lines
			if !ok {
				logs.Warning.Printf("tail file close reopen,file name:%s\n", task.Tail.Filename)
				time.Sleep(time.Second)
				continue
			}
			// 空行略过
			if len(strings.Trim(line.Text, "\r")) == 0 {
				continue
			}
			//发送到kafka
			msg := &sarama.ProducerMessage{}
			msg.Topic = task.Topic
			msg.Value = sarama.StringEncoder(line.Text)
			kafka.SendMessage(msg)
		}

	}
}

func Start() {
	for {
		Run()
		<-utils.Flag
		logs.Info.Println("Server ReStart Success")
	}
}
