package main

import (
	"log"
	"time"

	"golang.org/x/net/context"

	"baymax/pubsub"
	"github.com/satori/go.uuid"
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

	server := pubsub.NewServer("127.0.0.1:4150")

	queueName := uuid.NewV4().String()
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

	select {
	case <-server.Exit:
	}

}
