package main

import (
	"fmt"
	"log"
	"time"

	"github.com/werunclub/baymax/v2/pubsub/broker"
	"github.com/werunclub/baymax/v2/pubsub/broker/nats"
)

var (
	topic = "go.baymax.topic.testing"
)

func doPub(mq broker.Broker) {
	tick := time.NewTicker(time.Second)
	i := 0
	for _ = range tick.C {
		msg := &broker.Message{
			Header: map[string]string{
				"id": fmt.Sprintf("%d", i),
			},
			Body: []byte(fmt.Sprintf("%d: %s", i, time.Now().String())),
		}
		if err := mq.Publish(topic, msg); err != nil {
			log.Printf("[pub] failed: %v", err)
		} else {
			fmt.Println("[pub] pubbed message:", string(msg.Body))
		}
		i++
	}
}

func sub(mq broker.Broker) {
	_, err := mq.Subscribe(topic, func(p broker.Publication) error {
		fmt.Println("[sub] received message:", string(p.Message().Body), "header", p.Message().Header)
		return nil
	}, broker.Queue("testing"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	mq := nats.NewNatsBroker(broker.Addrs("127.0.0.1:4222"))

	if err := mq.Connect(); err != nil {
		log.Fatalf("Broker Connect error: %v", err)
	}

	go doPub(mq)
	go sub(mq)
	go sub(mq)

	<-time.After(time.Second * 30)
}
