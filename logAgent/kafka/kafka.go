package kafka

import (
	"logsys/logAgent/logs"
	"logsys/logAgent/utils"
	"github.com/Shopify/sarama"
)

var (
	client sarama.SyncProducer
	//	利用通道将同步代码改为异步的
	message chan *sarama.ProducerMessage
	err     error
)

// 初始化一个全局的kafka连接
func InitKafka() {
	// 初始化配置
	config := sarama.NewConfig()
	// 确认模式
	config.Producer.RequiredAcks = sarama.WaitForAll
	//	分区
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	//	成功交付信息
	config.Producer.Return.Successes = true

	//	连接kafka
	client, err = sarama.NewSyncProducer(utils.Cfg.KafkaConfig.Address, config)
	if err != nil {
		logs.Error.Fatalln("kafka producer closed ,err:", err)
	}

	logs.Info.Println("init kafka producer success")

	// 初始化通道
	message = make(chan *sarama.ProducerMessage, utils.Cfg.ChanSize)
	// 起一个子goroutine发送消息
	go sendKafka()

}

// Msg 中读取消息发送给kafka
func sendKafka() {
	for {
		select {
		case msg := <-message:
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				logs.Warning.Println("send msg to kafka failed,err:", err)
				return
			}
			logs.Trace.Printf("send msg to kafka success,topic:%v, pid:%v, offset:%v\n", msg.Topic, pid, offset)
		}
	}
}

// 需要执行关闭操作的变量在此集中关闭
// 在主进程在defer注册
func Close() {
	close(message)
	client.Close()
}

// 对外暴露发送消息的接口
func SendMessage(msg *sarama.ProducerMessage) {
	message <- msg
}
