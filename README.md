# collectlog
日志收集项目

##logAgent 
logAgent从etcd拉取配置项，根据配置项创建任务列表读取日志文件，发送到kafka中

## Transfer
Transfer从kafka消费信息，通过一个channel发送到es中便于处理
