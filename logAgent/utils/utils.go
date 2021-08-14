package utils

import (
	"fmt"
	"gopkg.in/ini.v1"
	"logsys/logAgent/logs"
	"net"
	"strings"
)

type KafkaConfig struct {
	Address  []string
	ChanSize int
}

type Etcd struct {
	Address    []string
	CollectKey string
}

type Config struct {
	KafkaConfig
	Etcd
}

var Cfg Config

func InitConf() {
	// 加载配置文件
	cfg, err := ini.Load("conf/config.ini")
	if err != nil {
		logs.Error.Fatalln("load config failed,err:", err)
	}

	// KafkaConfig
	str := cfg.Section("kafka").Key("Address").String()
	Cfg.KafkaConfig.Address = strings.Split(str, ",")
	if len(Cfg.KafkaConfig.Address) == 0 || Cfg.KafkaConfig.Address[0] == "" {
		logs.Warning.Fatalln("kafka address not placed")
	}
	Cfg.ChanSize = cfg.Section("kafka").Key("ChanSize").MustInt(10000)

	//	Etcd
	str = cfg.Section("etcd").Key("Address").String()
	Cfg.Etcd.Address = strings.Split(str, ",")
	if len(Cfg.Etcd.Address) == 0 || Cfg.Etcd.Address[0] == "" {
		logs.Warning.Fatalln("etcd cluster address not placed")
	}
	key := cfg.Section("etcd").Key("CollectKey").String()
	Cfg.CollectKey = fmt.Sprintf(key, GetIP())

	fmt.Println(Cfg.CollectKey)
}

// 获取本机IP 方便不同服务器拉去自己的配置
func GetIP() string {
	conn, err := net.Dial("udp", "127.0.0.1:8888")
	if err != nil {
		logs.Warning.Fatalln("Get Local IP failed,err:", err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return strings.Split(localAddr.IP.String(), ":")[0]
}
