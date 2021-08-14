package kafka

import (
	"github.com/Shopify/sarama"
	"logsys/transfer/logs"
	"logsys/transfer/utils"
)

// 用于同步读取消息进程和主进程
var (
	DataChan chan string = make(chan string, utils.Cfg.ChanSize)
	consumer sarama.Consumer
)

// kafka Consumer demo
func TransferTask() {
	//1、 创建消费者实例
	c, err := sarama.NewConsumer(utils.Cfg.KafkaConf.Address, nil)
	if err != nil {
		logs.Warning.Fatalln("fail to start consumer,err: ", err)
	}
	consumer = c

	//1、 通过 topic 获取所有 kafka 分区信息
	partitionList, err := consumer.Partitions(utils.Cfg.KafkaConf.Topic)
	if err != nil {
		logs.Warning.Printf("get partition failed,err:%v\n", err)
		return
	}

	//2、 遍历所有分区
	for partition := range partitionList {
		//3、	针对每个分区创建对应的分区消费者
		pc, err := consumer.ConsumePartition(utils.Cfg.KafkaConf.Topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			logs.Warning.Printf("filed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}

		//4、 异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				data := string(msg.Value)
				DataChan <- data
			}
		}(pc)
	}
	logs.Info.Println("init consumer task success")
}

func Close() {
	consumer.Close()
}
