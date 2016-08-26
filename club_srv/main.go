package main

import (
	"club-backend/club_srv/handler"
	"club-backend/club_srv/db"
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func run() {
	handler := new(handler.ClubHandler)
	server := rpc.NewServer()
	server.Register(handler)

	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
			log.Fatal(err)
		}

		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

func init() {
	const usage = "club_srv [-c config_file][-p cpupro file][-m mempro file]"
	c := flag.String("c", "", usage)
	log.Print(c)
}

func main() {

	flag.Parse()

	db.Init()

	run()
}
