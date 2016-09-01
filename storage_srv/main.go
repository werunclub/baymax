package main

import (
	"flag"

	"baymax/rpc"
	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/configor"
)

var log *logrus.Logger

func init() {
	log = logrus.New()

	var (
		addr   string
		config string
	)

	flag.StringVar(&config, "c", "", "config file")
	flag.StringVar(&addr, "addr", "", "addr, exmaple: 0.0.0.0:8080")
	flag.Parse()

	if config != "" {
		configor.Load(&Config, config)
	}

	if addr != "" {
		Config.Server.Addr = addr
	}
}

func main() {
	server := rpc.NewServer()
	server.RegisterName("Storage", new(storageHandler))

	// 启动服务
	server.Serve("tcp", Config.Server.Addr)
}
