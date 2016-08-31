package main

import (
	"flag"

	"baymax/club_srv/model"
	"baymax/club_srv/handler"
	"baymax/rpc"

	"github.com/jinzhu/configor"
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

	// 连接数据库
	model.Init(Config.Database.DSN)

	server := rpc.NewServer()
	server.RegisterName("Club", new(handler.ClubHandler))

	// 启动服务
	server.Serve("tcp", Config.Server.Addr)
}
