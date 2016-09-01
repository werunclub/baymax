package broker

import (
	"encoding/json"
	"math/rand"
	"sync"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"
)

var (
	concurrentHandlerKey = contextKeyT("baymax/broker/nsq/concurrentHandlers")
)

type nsqBroker struct {
	addrs  []string
	opts   Options
	config *nsq.Config

	sync.Mutex
	running bool
	p       []*nsq.Producer
	c       []*nsqSubscriber
}

type nsqPublication struct {
	topic string
	m     *Message
	nm    *nsq.Message
	opts  PublishOptions
}

type nsqSubscriber struct {
	topic string
	opts  SubscribeOptions

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

func (n *nsqBroker) Init(opts ...Option) error {
	for _, o := range opts {
		o(&n.opts)
	}
	return nil
}

func (n *nsqBroker) Options() Options {
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

	var producers []*nsq.Producer

	// create producers
	for _, addr := range n.addrs {
		p, err := nsq.NewProducer(addr, n.config)
		if err != nil {
			return err
		}

		producers = append(producers, p)
	}

	// create consumers
	for _, c := range n.c {
		channel := c.opts.Queue
		if len(channel) == 0 {
			channel = uuid.NewUUID().String()
		}

		cm, err := nsq.NewConsumer(c.topic, channel, n.config)
		if err != nil {
			return err
		}

		cm.AddConcurrentHandlers(c.h, c.n)

		c.c = cm

		//err = c.c.ConnectToNSQLookupds(n.addrs)
		err = c.c.ConnectToNSQDs(n.addrs)
		if err != nil {
			return err
		}
	}

	n.p = producers
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
	for _, p := range n.p {
		p.Stop()
	}

	// stop the consumers
	for _, c := range n.c {
		c.c.Stop()

		// disconnect from all nsq brokers
		for _, addr := range n.addrs {
			c.c.DisconnectFromNSQD(addr)
			//c.c.DisconnectFromNSQLookupd(addr)
		}
	}

	n.p = nil
	n.running = false
	return nil
}

func (n *nsqBroker) Publish(topic string, message *Message, opts ...PublishOption) error {
	p := n.p[rand.Int()%len(n.p)]

	b, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return p.Publish(topic, b)
}

func (n *nsqBroker) Subscribe(topic string, handler Handler, opts ...SubscribeOption) (Subscriber, error) {
	options := SubscribeOptions{
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
		channel = uuid.NewUUID().String()
	}

	c, err := nsq.NewConsumer(topic, channel, n.config)
	if err != nil {
		return nil, err
	}

	h := nsq.HandlerFunc(func(nm *nsq.Message) error {
		if !options.AutoAck {
			nm.DisableAutoResponse()
		}

		var m *Message

		if err := json.Unmarshal(nm.Body, &m); err != nil {
			return err
		}

		return handler(&nsqPublication{
			topic: topic,
			m:     m,
			nm:    nm,
		})

	})

	c.AddConcurrentHandlers(h, concurrency)

	err = c.ConnectToNSQDs(n.addrs)
	//err = c.ConnectToNSQLookupds(n.addrs)
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

func (p *nsqPublication) Message() *Message {
	return p.m
}

func (p *nsqPublication) Ack() error {
	p.nm.Finish()
	return nil
}

func (s *nsqSubscriber) Options() SubscribeOptions {
	return s.opts
}

func (s *nsqSubscriber) Topic() string {
	return s.topic
}

func (s *nsqSubscriber) Unsubscribe() error {
	s.c.Stop()
	return nil
}

func NewNsqBroker(opts ...Option) Broker {
	var options Options
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
		cAddrs = []string{"127.0.0.1:4150"}
	}

	return &nsqBroker{
		addrs:  cAddrs,
		opts:   options,
		config: nsq.NewConfig(),
	}
}

func ConcurrentHandlers(n int) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Context = context.WithValue(o.Context, concurrentHandlerKey, n)
	}
}
