package main

var Config = struct {
	Debug       bool   `default:"false"`
	ServiceName string `default:"user_srv"`

	Server struct {
		Addr string `default:":8080"`
	}

	Database struct {
		Address string
	}

	Registry struct {
		Type    string
		Address string
	}
}{}
