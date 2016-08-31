package main

import (
	"log"

	"golang.org/x/net/context"

	"baymax/pubsub"
)

type Message struct {
	Say string `protobuf:"bytes,1,opt,name=say" json:"say,omitempty"`
}

func Handler(ctx context.Context, msg *Message) error {
	log.Print("good: ", msg.Say)
	return nil
}

func BadHandler(ctx context.Context, msg *Message) error {
	log.Print("bad: ", msg.Say)
	return nil
}

func main() {

	server := pubsub.NewServer()

	var err error

	err = server.Subscribe(pubsub.NewSubscriber(
		"go.testing.topic.good",
		Handler,
	))

	if err != nil {
		log.Fatalf("fail")
	}

	err = server.Subscribe(pubsub.NewSubscriber(
		"go.testing.topic.bad",
		BadHandler,
	))

	if err != nil {
		log.Fatalf("fail")
	}

	server.Run()
}
