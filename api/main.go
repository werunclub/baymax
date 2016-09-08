package main

import (
	"flag"
	"github.com/jinzhu/configor"
	"runtime"
	"github.com/Sirupsen/logrus"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// settingLogrus 这里设置的是 logrus 的 std logger
func settingLogrus() {
	debugLevel, err:= logrus.ParseLevel(Config.Logger.Level)
	if err != nil {
		debugLevel = logrus.InfoLevel
		logrus.WithError(err).Warningf("接收到粗无的参数 debug[%v], 默认使用 logrus.InfoLevel", debugLevel)
	}
	logrus.SetLevel(debugLevel)
	if Config.Logger.Formatter == "text" {
		logrus.SetFormatter(&logrus.TextFormatter{})
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
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

	settingLogrus()
	StartServer(Config.Server.Addr)
}
