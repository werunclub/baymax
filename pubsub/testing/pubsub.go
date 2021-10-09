package main

import (
	"fmt"
	"log"
	"time"

	"github.com/werunclub/baymax/v2/pubsub/broker"
)

var (
	topic = "go.baymax.topic.testing"
)

func pub() {
	tick := time.NewTicker(time.Second)
	i := 0
	for _ = range tick.C {
		msg := &broker.Message{
			Header: map[string]string{
				"id": fmt.Sprintf("%d", i),
			},
			Body: []byte(fmt.Sprintf("%d: %s", i, time.Now().String())),
		}
		if err := broker.Publish(topic, msg); err != nil {
			log.Printf("[pub] failed: %v", err)
		} else {
			fmt.Println("[pub] pubbed message:", string(msg.Body))
		}
		i++
	}
}

func sub() {
	_, err := broker.Subscribe(topic, func(p broker.Publication) error {
		fmt.Println("[sub] received message:", string(p.Message().Body), "header", p.Message().Header)
		return nil
	}, broker.Queue("testing"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	err := broker.Init(broker.Addrs("127.0.0.1:4150"))

	if err != nil {
		log.Fatalf("Broker Init error: %v", err)
	}
	if err := broker.Connect(); err != nil {
		log.Fatalf("Broker Connect error: %v", err)
	}

	go pub()
	go sub()
	go sub()

	<-time.After(time.Second * 30)
}
