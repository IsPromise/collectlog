package utils

import (
	"context"
	"encoding/json"
	clientv3 "go.etcd.io/etcd/client/v3"
	"logsys/logAgent/logs"
	"time"
)

var (
	client *clientv3.Client
	err    error
)

type CollectEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

// etcd 管理配置项
// key-value
// 其中 value 为json格式字符串 格式如下
// [{"path":"d:/Log/1.log","topic":"test_log"},{"path":"D/log/2.log","topic":"web_log"}]

// 初始化Etcd
func InitEtcd() {
	client, err = clientv3.New(clientv3.Config{
		DialTimeout: 5 * time.Second,
		Endpoints:   Cfg.Etcd.Address,
	})
	if err != nil {
		logs.Warning.Fatalln("Create etcd Client failed,err: ", err)
	}

	logs.Info.Println("init etcd success")

	// 启动子goroutine 监听配置文件变化
	go WatchConf(Cfg.Etcd.CollectKey)
}

// etcd 获取配置信息
func getConf(key string) (list []*CollectEntry) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	gresp, err := client.Get(ctx, key)
	if err != nil {
		logs.Warning.Fatalln("get config from etcd failed,err: ", err)
	}

	if len(gresp.Kvs) == 0 {
		logs.Warning.Fatalln("config is nil")
	}

	ret := gresp.Kvs[0]
	//	反序列化
	err = json.Unmarshal(ret.Value, &list)
	if err != nil {
		logs.Warning.Fatalln("config unmarshal failed,err: ", err)
	}

	logs.Info.Println("config get success ")
	return
}

// 日志项管理 watch 去监控 etcd collect_conf 的变化
func WatchConf(key string) {
	wch := client.Watch(context.Background(), key)
	var newConf []*CollectEntry
	for resp := range wch {
		for _, evt := range resp.Events {
			logs.Info.Println("Config to be update")
			err := json.Unmarshal(evt.Kv.Value, &newConf)
			// 出错结束本层循环继续阻塞等待新的配置项
			if err != nil {
				logs.Warning.Printf("config update failed about config unmarshal failed,err:%v ", err)
				break
			}
			logs.Trace.Printf("update about key:%s ,value:%s\n", evt.Kv.Key, evt.Kv.Value)
			//	通知tail 配置项更新
			Notice()
		}
	}
}
