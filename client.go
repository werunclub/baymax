package main

import (
	"baymax/club_srv/protocol/club"
	"baymax/rpc/client"
	"fmt"
)

func main() {
	c := client.NewClient("tcp", ":8081")
	defer c.Close()

	var reply club.GetResponse

	err := c.Call(club.ServiceCreateClub, &club.GetRequest{ClubId: 100}, &reply)
	if err != nil {
		panic(err)
	}

	fmt.Println(reply)
}
