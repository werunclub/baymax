package main

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/werunclub/baymax/v2/pubsub"
	"github.com/werunclub/baymax/v2/pubsub/metadata"
)

type Message struct {
	Say string `protobuf:"bytes,1,opt,name=say" json:"say,omitempty"`
}

// publishes a message
func pub(i int) {
	client := pubsub.NewClient("127.0.0.1:4161")

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
	fmt.Printf("\n--- Publisher example ---\n")
	for i := 0; i < 1; i++ {
		pub(i)
	}
}
