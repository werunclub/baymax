package pubsub

import (
	"bytes"
	"sync"

	"github.com/werunclub/baymax/v2/pubsub/broker"
	"github.com/werunclub/baymax/v2/pubsub/broker/nats"
	"github.com/werunclub/baymax/v2/pubsub/broker/nsq"
	"github.com/werunclub/baymax/v2/pubsub/codec"
	mj "github.com/werunclub/baymax/v2/pubsub/codec/jsonrpc"
	"github.com/werunclub/baymax/v2/pubsub/metadata"

	log "github.com/sirupsen/logrus"
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

func NewClient(addrs ...string) *Client {
	return NewClientWithName(getBrokerName(), addrs...)
}

func NewClientWithName(brokerName string, addrs ...string) *Client {
	opts := broker.Addrs(addrs...)

	var b broker.Broker
	if brokerName == "nats" {
		b = nats.NewNatsBroker(opts)
	} else {
		b = nsq.NewNsqBroker(opts)
	}

	return &Client{
		broker: b,
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
		log.Errorf("encode fail:%s", err.Error())
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

func (c *Client) Close() error {
	return c.broker.Disconnect()
}
