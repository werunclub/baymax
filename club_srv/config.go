package main

var Config = struct {
	ServiceName string `default:"club_srv"`

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
