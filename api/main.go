package main

import (
	"flag"
	"github.com/jinzhu/configor"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
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

	StartServer(Config.Server.Addr)
}
