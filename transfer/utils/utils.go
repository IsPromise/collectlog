package utils

import (
	"gopkg.in/ini.v1"
	"logsys/transfer/logs"
	"strings"
)

type Config struct {
	EsConf
	KafkaConf
}

type KafkaConf struct {
	Address []string
	Topic   string
}

type EsConf struct {
	Address      string
	Index        string
	ChanSize     int
	GoroutineNum int
}

var Cfg Config

func InitConf() {
	cfg, err := ini.Load("./conf/transfer.ini")
	if err != nil {
		logs.Warning.Fatalln("init transfer config failed,err ", err)
	}

	Cfg.EsConf.Address = cfg.Section("es").Key("Address").String()
	Cfg.EsConf.Index = cfg.Section("es").Key("Index").String()
	Cfg.EsConf.ChanSize = cfg.Section("es").Key("ChanSize").MustInt(10000)
	Cfg.EsConf.GoroutineNum = cfg.Section("es").Key("GoroutineNum").MustInt(16)

	str := cfg.Section("kafka").Key("Address").String()
	Cfg.KafkaConf.Address = strings.Split(str, ",")
	Cfg.KafkaConf.Topic = cfg.Section("kafka").Key("Topic").String()

	logs.Info.Println("load transfer config success")
}
