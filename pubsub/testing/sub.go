package main

import (
	"log"

	"golang.org/x/net/context"

	"baymax/pubsub"
	"time"
)

type Message struct {
	Say string `protobuf:"bytes,1,opt,name=say" json:"say,omitempty"`
}

func Handler(ctx context.Context, msg *Message) error {
	//log.Print("rev:", msg.Say)
	time.Sleep(time.Second * 10)
	log.Print("done:", msg.Say)
	return nil
}

func BadHandler(ctx context.Context, msg *Message) error {
	log.Print("bad: ", msg.Say)
	return nil
}

func main() {

	server := pubsub.NewServer("127.0.0.1:4150")

	var err error

	err = server.Subscribe(server.NewSubscriber(
		"go.testing.topic.good",
		Handler,
		pubsub.SubscriberQueue("testing"),
	))

	if err != nil {
		log.Fatalf("fail")
	}

	err = server.Subscribe(server.NewSubscriber(
		"go.testing.topic.bad",
		BadHandler,
		pubsub.SubscriberQueue("testing"),
	))

	if err != nil {
		log.Fatalf("fail")
	}

	server.Run()
}
