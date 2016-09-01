package main

import (
	clubProto "baymax/club_srv/protocol/club"
	"baymax/rpc/client"
	"flag"
	"log"
	"math/rand"
)

// 测试获取单条俱乐部信息
func testClubGetOne(c *client.Client) error {
	var reply clubProto.GetResponse

	err := c.Call(clubProto.ServiceGetClub, &clubProto.GetRequest{ClubID: 45}, &reply)
	if err != nil {
		return err
	} else {
		log.Print(reply)
		return nil
	}
}

// 测试创建新的俱乐部
func testClubCreate(c *client.Client) error {
	var (
		req clubProto.CreateRequest
		res clubProto.CreateResponse
		err error
	)

	req.Name = "2048"

	err = c.Call(clubProto.ServiceCreateClub, &req, &res)
	if err != nil {
		return err
	} else {
		log.Print(res)
		return nil
	}
}

func testClubUpdate(c *client.Client) error {

	var (
		req clubProto.UpdateRequest
		res clubProto.UpdateResponse
		err error
	)

	req.Name = (string)(rand.Intn(1024))
	req.ID = 316

	err = c.Call(clubProto.ServiceUpdateClub, &req, &res)
	if err != nil {
		return err
	} else {
		log.Print(res)
	}

	return nil
}

func testClubDelete(c *client.Client) error {

	var (
		req clubProto.DeleteRequest
		res clubProto.DeleteResponse
		err error
	)

	req.ClubID = 317
	err = c.Call(clubProto.ServiceDeleteClub, &req, &res)
	if err != nil {
		return err
	} else {
		log.Print(err)
		return nil
	}
}

func testClubGetMany(c *client.Client) error {
	return nil
}

func main() {
	var (
		port string
		err error
	)

	flag.StringVar(&port, "port", "8080", "server port")
	flag.Parse()

	c := client.NewClient("tcp", ":" + port)
	defer c.Close()

	// 测试 club 部分的 rpc

	err = testClubGetOne(c)
	if err != nil {
		log.Fatal(err)
	}

	err = testClubCreate(c)
	if err != nil {
		log.Fatal(err)
	}

	err = testClubDelete(c)
	if err != nil {
		log.Fatal(err)
	}

	err = testClubUpdate(c)
	if err != nil {
		log.Fatal(err)
	}

	err = testClubGetMany(c)
	if err != nil {
		log.Fatal(err)
	}
}
