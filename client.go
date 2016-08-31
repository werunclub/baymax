package main

import (
	"baymax/club_srv/protocol/club"
	"baymax/rpcx"
	"fmt"
	"log"
	"time"
)

func main() {
	client := rpcx.NewClient("tcp", ":8085", time.Duration(24*365)*time.Hour)
	defer client.Close()

	var reply club.GetResponse

	err := client.Call(club.ServiceCreateClub, &club.GetRequest{ClubId: 100}, &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)
}
