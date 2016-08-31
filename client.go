package main

import (
	"baymax/rpc"
	"baymax/storage_srv/protocol/storage"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func main() {
	var (
		port string
		err  error
	)

	client := rpc.NewClient("tcp", ":8080", time.Duration(24*365)*time.Hour)
	//defer client.Close()

	var (
		//log   = logrus.New()
		req   storage.StorePhotoArgs
		reply storage.StorePhotoReply
	)

	f, err1 := os.Open("test.jpg")
	if err1 != nil {
		log.Fatal(err1)
	}

	photo, err2 := ioutil.ReadAll(f)
	if err2 != nil {
		log.Fatal(err2)
	}

	req = storage.StorePhotoArgs{
		UserId:   100,
		FileType: "jpg",
		FileSize: 100000,
		Photo:    photo,
	}

	err := client.Call(storage.StorePhoto, &req, &reply)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(reply.Url, reply.Suffixes, reply.Width, reply.Height)
}
