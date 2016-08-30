package main

import (
	"flag"

	"baymax/club_srv/db"
	"baymax/club_srv/handler"
	"baymax/rpc/server"

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
	db.Init(Config.Database.Address)

	s := server.NewServer()
	s.RegisterName("Club", new(handler.ClubHandler))

	// 启动服务
	s.Serve("tcp", Config.Server.Addr)
}
