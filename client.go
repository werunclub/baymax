package main

import (
	clubProto "baymax/club_srv/protocol/club"
	"baymax/rpc/client"
	"fmt"
	"log"
	"time"
	"flag"
)

func testClubGet(c *client.Client) error {
	var reply clubProto.GetResponse

	err := c.Call(clubProto.ServiceGetClub, &clubProto.GetRequest{ClubId: 1024}, &reply)
	if err != nil {
		return err
	} else {
		fmt.Println(reply)
		return nil
	}
}

func testClubCreate(_ *client.Client) error {
	return nil
}

func testClubGetAllCreate(_ *client.Client) error {
	return nil
}

func testClubDelete(_ *client.Client) error {
	return nil
}

func main() {
	var port string

	flag.StringVar(&port, "port", "8080", "server port")
	flag.Parse()

	c := client.NewClient("tcp", ":" + port)
	defer c.Close()

	err := testClubGet(c)
	if err != nil {
		log.Fatal(err)
	}
}
