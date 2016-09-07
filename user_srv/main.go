package main

import (
	"flag"
	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/configor"
	"baymax/user_srv/handler"
	"baymax/user_srv/model"
	"baymax/rpc"
)

func init() {
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
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	model.Init(Config.Database.Address, Config.Debug)

	rpcServer := rpc.NewServer()
	//rpcServer := rpc.NewServer()
	handler.RegisterRPCService(rpcServer)

	logrus.WithField("Address", Config.Server.Addr).Info("RPCServer Listening")
	//listener, _ := net.Listen("tcp", Config.Server.Addr)
	//rpcServer.Accept(listener)
	rpcServer.Serve("tcp", Config.Server.Addr)
}
