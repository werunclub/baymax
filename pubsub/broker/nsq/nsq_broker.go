package nsq

import (
	"encoding/json"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/werunclub/baymax/v2/pubsub/broker"
	"github.com/werunclub/go-nsq"
	"golang.org/x/net/context"
)

type contextKeyT string

var (
	concurrentHandlerKey = contextKeyT("github.com/werunclub/baymax/v2/broker/nsq/concurrentHandlers")
)

type nsqBroker struct {
	addrs  []string
	opts   broker.Options
	config *nsq.Config

	sync.Mutex
	running bool
	d       *nsq.Driver
	p       []*nsq.Producer
	c       []*nsqSubscriber
}

type nsqPublication struct {
	topic string
	m     *broker.Message
	nm    *nsq.Message
	opts  broker.PublishOptions
}

type nsqSubscriber struct {
	topic string
	opts  broker.SubscribeOptions

	c *nsq.Consumer

	// handler so we can resubcribe
	h nsq.HandlerFunc
	// concurrency
	n int
}

var (
	DefaultConcurrentHandlers = 1
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (n *nsqBroker) Init(opts ...broker.Option) error {
	for _, o := range opts {
		o(&n.opts)
	}
	return nil
}

func (n *nsqBroker) Options() broker.Options {
	return n.opts
}

func (n *nsqBroker) Address() string {
	return n.addrs[rand.Int()%len(n.addrs)]
}

func (n *nsqBroker) Connect() error {
	n.Lock()
	defer n.Unlock()

	if n.running {
		return nil
	}

	// create producers
	n.d = nsq.NewProducerDriver(n.config)
	if err := n.d.ConnectToNSQLookupds(n.addrs); err != nil {
		return err
	}

	// create consumers
	for _, c := range n.c {
		channel := c.opts.Queue
		if len(channel) == 0 {
			channel = uuid.NewString()
		}

		cm, err := nsq.NewConsumer(c.topic, channel, n.config)
		if err != nil {
			return err
		}

		cm.AddConcurrentHandlers(c.h, c.n)

		c.c = cm

		err = c.c.ConnectToNSQLookupds(n.addrs)
		//err = c.c.ConnectToNSQDs(n.addrs)
		if err != nil {
			return err
		}
	}

	n.running = true
	return nil
}

func (n *nsqBroker) Disconnect() error {
	n.Lock()
	defer n.Unlock()

	if !n.running {
		return nil
	}

	// stop the producers
	n.d.Stop()

	// stop the consumers
	for _, c := range n.c {
		c.c.Stop()

		// disconnect from all nsq brokers
		for _, addr := range n.addrs {
			//c.c.DisconnectFromNSQD(addr)
			c.c.DisconnectFromNSQLookupd(addr)
		}
	}

	n.p = nil
	n.running = false
	return nil
}

func (n *nsqBroker) Publish(topic string, message *broker.Message, opts ...broker.PublishOption) error {
	b, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return n.d.Publish(topic, b)
}

func (n *nsqBroker) Subscribe(topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	options := broker.SubscribeOptions{
		AutoAck: true,
	}

	for _, o := range opts {
		o(&options)
	}

	var concurrency int

	if options.Context != nil {
		var ok bool
		concurrency, ok = options.Context.Value(concurrentHandlerKey).(int)
		if !ok {
			concurrency = DefaultConcurrentHandlers
		}
	} else {
		concurrency = DefaultConcurrentHandlers
	}

	channel := options.Queue
	if len(channel) == 0 {
		channel = uuid.NewString()
	}

	c, err := nsq.NewConsumer(topic, channel, n.config)
	if err != nil {
		return nil, err
	}

	h := nsq.HandlerFunc(func(nm *nsq.Message) error {
		if !options.AutoAck {
			nm.DisableAutoResponse()
		}

		var m *broker.Message

		if err := json.Unmarshal(nm.Body, &m); err != nil {
			logrus.WithError(err).Errorf("can not parse message")
			return err
		}

		return handler(&nsqPublication{
			topic: topic,
			m:     m,
			nm:    nm,
		})
	})

	c.AddConcurrentHandlers(h, concurrency)

	//err = c.ConnectToNSQDs(n.addrs)
	err = c.ConnectToNSQLookupds(n.addrs)
	if err != nil {
		return nil, err
	}

	return &nsqSubscriber{
		topic: topic,
		c:     c,
		h:     h,
		n:     concurrency,
	}, nil
}

func (n *nsqBroker) String() string {
	return "nsq"
}

func (p *nsqPublication) Topic() string {
	return p.topic
}

func (p *nsqPublication) Message() *broker.Message {
	return p.m
}

func (p *nsqPublication) Ack() error {
	p.nm.Finish()
	return nil
}

func (p *nsqPublication) Error() error {
	return nil
}

func (s *nsqSubscriber) Options() broker.SubscribeOptions {
	return s.opts
}

func (s *nsqSubscriber) Topic() string {
	return s.topic
}

func (s *nsqSubscriber) Unsubscribe() error {
	s.c.Stop()
	return nil
}

func NewNsqBroker(opts ...broker.Option) broker.Broker {
	var options broker.Options
	for _, o := range opts {
		o(&options)
	}

	var cAddrs []string

	for _, addr := range options.Addrs {
		if len(addr) > 0 {
			cAddrs = append(cAddrs, addr)
		}
	}

	if len(cAddrs) == 0 {
		cAddrs = []string{"127.0.0.1:4161"}
	}

	config := nsq.NewConfig()
	config.MaxInFlight = 12
	config.LookupdPollInterval = time.Second * 5
	config.MaxAttempts = 10

	return &nsqBroker{
		addrs:  cAddrs,
		opts:   options,
		config: config,
	}
}

func ConcurrentHandlers(n int) broker.SubscribeOption {
	return func(o *broker.SubscribeOptions) {
		o.Context = context.WithValue(o.Context, concurrentHandlerKey, n)
	}
}
