package main

import (
	"logsys/logAgent/agent"
	"logsys/logAgent/kafka"
	"logsys/logAgent/logs"
	"logsys/logAgent/utils"
)

// 日志收集的客户端
// 需要实现收集指定目录下的日志文件，发送到Kafka中

func main() {
	//	1、读取配置文件
	utils.InitConf()

	// 初始化kafka
	kafka.InitKafka()
	defer kafka.Close()
	// 程序整个退出时关闭日志文件
	defer logs.Close()

	// 初始化etcd
	utils.InitEtcd()

	//	2、初始化收集日志
	utils.InitTail()
	//	3、发送数据到Kafka
	agent.Start()

}
