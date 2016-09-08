package main

var Config = struct {
	Debug   bool   `default:"false"`
	APPName string `default:"api"`

	Server struct {
		Addr string `default:":8080"`
	}

	Registry struct {
		Type    string
		Address string
	}
	Logger struct {
		Level     string `default:"info"`
		Formatter string `default:"json"`
	}
}{}
