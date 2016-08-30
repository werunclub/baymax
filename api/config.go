package main

var Config = struct {
	APPName string `default:"api"`

	Server struct {
		Addr string `default:":8080"`
	}

	Registry struct {
		Type    string
		Address string
	}
}{}
