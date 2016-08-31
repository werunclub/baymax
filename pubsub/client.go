package pubsub

import (
	"bytes"
	"sync"

	"baymax/pubsub/broker"
	"github.com/micro/go-micro/codec"
	mj "github.com/micro/go-micro/codec/jsonrpc"
	"github.com/micro/go-micro/metadata"
	"golang.org/x/net/context"
)

var (
	contentType = "application/json"
)

// Client represents a RPC client.
type Client struct {
	broker broker.Broker
	once   sync.Once
}

func NewClient(opts ...broker.Option) *Client {
	return &Client{
		broker: broker.NewBroker(opts...),
		once:   sync.Once{},
	}
}

func (c *Client) NewPublication(topic string, message interface{}) publication {
	return newPublication(topic, message)
}

func (c *Client) Publish(ctx context.Context, p publication) error {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = make(map[string]string)
	}

	// encode message body
	b := &buffer{bytes.NewBuffer(nil)}
	if err := mj.NewCodec(b).Write(&codec.Message{Type: codec.Publication}, p.Message()); err != nil {
		return err
	}
	c.once.Do(func() {
		c.broker.Connect()
	})

	return c.broker.Publish(p.Topic(), &broker.Message{
		Header: md,
		Body:   b.Bytes(),
	})
}
