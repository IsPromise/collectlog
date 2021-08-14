package es

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic"
	"logsys/transfer/kafka"
	"logsys/transfer/logs"
	"logsys/transfer/utils"
)

type data struct {
	Data string `json:"data"`
}

func InitEs() {
	client, err := elastic.NewClient(elastic.SetURL("http://" + utils.Cfg.EsConf.Address))
	if err != nil {
		logs.Error.Fatalln("new es client failed,err: ", err)
	}

	for i := 0; i < utils.Cfg.GoroutineNum; i++ {
		go sendToEs(client)
	}
	logs.Info.Println("init es success")
}

func formatData(s string) (str string, err error) {
	d := data{Data: s}
	b, err := json.Marshal(&d)
	if err != nil {
		return "", err
	}
	str = string(b)
	return str, nil
}

func sendToEs(client *elastic.Client) {
	//	从通道中获取数据
	for value := range kafka.DataChan {
		d, err := formatData(value)
		if err != nil {
			logs.Warning.Println("data format failed,err: ", err)
			continue
		}

		// 发送到es
		put, err := client.Index().
			Index(utils.Cfg.Index).
			Type("log").
			BodyJson(d).
			Do(context.Background())
		if err != nil {
			logs.Warning.Println("send msg to es failed,err: ", err)
			continue
		}
		logs.Trace.Printf("msg send es success,index:%v,type:%v\n", put.Index, put.Index)
	}
}
