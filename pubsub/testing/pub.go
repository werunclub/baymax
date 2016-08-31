package main

import (
	"fmt"

	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/metadata"
	"golang.org/x/net/context"

	"baymax/pubsub"
)

type Message struct {
	Say string `protobuf:"bytes,1,opt,name=say" json:"say,omitempty"`
}

// publishes a message
func pub2(i int) {

	client := pubsub.NewClient()

	msg := client.NewPublication("go.testing.topic.good", Message{
		Say: fmt.Sprintf("This is a publication %d", i),
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
	cmd.Init()
	fmt.Println("\n--- Publisher example ---\n")
	for i := 0; i < 10; i++ {
		pub2(i)
	}
}
