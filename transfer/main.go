package main

import (
	"logsys/transfer/es"
	"logsys/transfer/kafka"
	"logsys/transfer/logs"
	"logsys/transfer/utils"
)

func main() {
	// 初始化配置
	utils.InitConf()
	// 连接kafka 开始消费
	kafka.TransferTask()
	defer kafka.Close()
	// 初始化es
	es.InitEs()
	defer logs.Close()

	select {}
}
