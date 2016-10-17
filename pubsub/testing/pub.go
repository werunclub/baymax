package main

import (
	"fmt"

	"github.com/micro/go-micro/metadata"
	"golang.org/x/net/context"

	"baymax/pubsub"
	"time"
)

type Message struct {
	Say string `protobuf:"bytes,1,opt,name=say" json:"say,omitempty"`
}

// publishes a message
func pub(i int) {

	client := pubsub.NewClient("127.0.0.1:4150")

	msg := client.NewPublication("go.testing.topic.good", Message{
		Say: fmt.Sprintf("%d", i),
	})

	// create context with metadata
	ctx := metadata.NewContext(context.Background(), map[string]string{
		"X-User-Id": "john",
		"X-From-Id": "script",
	})

	// publish message
	if err := client.Publish(ctx, msg); err != nil {
		fmt.Println("pub err: ", err)
		return
	}

	fmt.Printf("Published %d: %v\n", i, msg)
}

func main() {
	fmt.Println("\n--- Publisher example ---\n")
	for i := 0; i < 1; i++ {
		pub(time.Now().Second())
	}
}
