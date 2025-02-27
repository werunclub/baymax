package main

import (
	"log"

	"github.com/google/uuid"
	"github.com/werunclub/baymax/v2/pubsub"
	"golang.org/x/net/context"
)

type Message struct {
	Say string `protobuf:"bytes,1,opt,name=say" json:"say,omitempty"`
}

func Handler(ctx context.Context, msg *Message) error {
	log.Print("done:", msg.Say)
	return nil
}

type SubHandler struct {
}

func (SubHandler) BadHandler(ctx context.Context, msg *Message) error {
	log.Print("bad1: ", msg.Say)
	return nil
}

func (SubHandler) Bad2Handler(ctx context.Context, msg *Message) error {
	log.Print("bad2: ", msg.Say)
	return nil
}

func main() {
	server := pubsub.NewServer("127.0.0.1:4161")

	queueName := uuid.NewString()
	var err error

	err = server.Subscribe(server.NewSubscriber(
		"go.testing.topic.good",
		Handler,
		pubsub.SubscriberQueue(queueName),
	))

	if err != nil {
		log.Fatal("fail")
	}

	err = server.Subscribe(server.NewSubscriber(
		"go.testing.topic.bad",
		new(SubHandler),
		pubsub.SubscriberQueue("testing"),
	))

	if err != nil {
		log.Fatal("fail")
	}

	go server.Run()

	//ch := make(chan os.Signal, 1)
	//signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	//fmt.Printf("Received signal %s", <-ch)

	<-server.Exit
}
