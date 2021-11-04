// Package nats provides a NATS broker
package broker

import (
	"context"
	"errors"
	"strings"
	"sync"

	nats "github.com/nats-io/nats.go"
	logger "github.com/sirupsen/logrus"

	"github.com/werunclub/baymax/v2/pubsub/codec"
	"github.com/werunclub/baymax/v2/pubsub/codec/json"
)

type natsBroker struct {
	sync.Once
	sync.RWMutex

	// indicate if we're connected
	connected bool

	addrs []string
	conn  *nats.Conn
	opts  Options
	nopts nats.Options

	codec codec.Marshaler

	// should we drain the connection
	drain   bool
	closeCh chan (error)
}

type natsSubscriber struct {
	s    *nats.Subscription
	opts SubscribeOptions
}

type natsPublication struct {
	t   string
	err error
	m   *Message
}

func (p *natsPublication) Topic() string {
	return p.t
}

func (p *natsPublication) Message() *Message {
	return p.m
}

func (p *natsPublication) Ack() error {
	// nats does not support acking
	return nil
}

func (p *natsPublication) Error() error {
	return p.err
}

func (s *natsSubscriber) Options() SubscribeOptions {
	return s.opts
}

func (s *natsSubscriber) Topic() string {
	return s.s.Subject
}

func (s *natsSubscriber) Unsubscribe() error {
	return s.s.Unsubscribe()
}

func (n *natsBroker) Address() string {
	if n.conn != nil && n.conn.IsConnected() {
		return n.conn.ConnectedUrl()
	}

	if len(n.addrs) > 0 {
		return n.addrs[0]
	}

	return ""
}

func (n *natsBroker) setAddrs(addrs []string) []string {
	//nolint:prealloc
	var cAddrs []string
	for _, addr := range addrs {
		if len(addr) == 0 {
			continue
		}
		if !strings.HasPrefix(addr, "nats://") {
			addr = "nats://" + addr
		}
		cAddrs = append(cAddrs, addr)
	}
	if len(cAddrs) == 0 {
		cAddrs = []string{nats.DefaultURL}
	}
	return cAddrs
}

func (n *natsBroker) Connect() error {
	n.Lock()
	defer n.Unlock()

	if n.connected {
		return nil
	}

	status := nats.CLOSED
	if n.conn != nil {
		status = n.conn.Status()
	}

	switch status {
	case nats.CONNECTED, nats.RECONNECTING, nats.CONNECTING:
		n.connected = true
		return nil
	default: // DISCONNECTED or CLOSED or DRAINING
		opts := n.nopts
		opts.Servers = n.addrs
		opts.Secure = n.opts.Secure
		opts.TLSConfig = n.opts.TLSConfig

		// secure might not be set
		if n.opts.TLSConfig != nil {
			opts.Secure = true
		}

		c, err := opts.Connect()
		if err != nil {
			return err
		}
		n.conn = c
		n.connected = true
		return nil
	}
}

func (n *natsBroker) Disconnect() error {
	n.Lock()
	defer n.Unlock()

	// drain the connection if specified
	if n.drain {
		n.conn.Drain()
		n.closeCh <- nil
	}

	// close the client connection
	n.conn.Close()

	// set not connected
	n.connected = false

	return nil
}

func (n *natsBroker) Init(opts ...Option) error {
	n.setOption(opts...)
	return nil
}

func (n *natsBroker) Options() Options {
	return n.opts
}

func (n *natsBroker) Publish(topic string, msg *Message, opts ...PublishOption) error {
	n.RLock()
	defer n.RUnlock()

	if n.conn == nil {
		return errors.New("not connected")
	}

	b, err := n.codec.Marshal(msg)
	if err != nil {
		return err
	}
	return n.conn.Publish(topic, b)
}

func (n *natsBroker) Subscribe(topic string, handler Handler, opts ...SubscribeOption) (Subscriber, error) {
	n.RLock()
	if n.conn == nil {
		n.RUnlock()
		return nil, errors.New("not connected")
	}
	n.RUnlock()

	opt := SubscribeOptions{
		AutoAck: true,
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&opt)
	}

	fn := func(msg *nats.Msg) {
		var m Message
		pub := &natsPublication{t: msg.Subject}
		errorHandler := n.opts.ErrorHandler
		err := n.codec.Unmarshal(msg.Data, &m)
		pub.err = err
		pub.m = &m
		if err != nil {
			m.Body = msg.Data
			logger.Error(err)

			if errorHandler != nil {
				errorHandler(pub)
			}
			return
		}
		if err := handler(pub); err != nil {
			pub.err = err
			logger.Error(err)
			if errorHandler != nil {
				errorHandler(pub)
			}
		}
	}

	var sub *nats.Subscription
	var err error

	n.RLock()
	if len(opt.Queue) > 0 {
		sub, err = n.conn.QueueSubscribe(topic, opt.Queue, fn)
	} else {
		sub, err = n.conn.Subscribe(topic, fn)
	}
	n.RUnlock()
	if err != nil {
		return nil, err
	}
	return &natsSubscriber{s: sub, opts: opt}, nil
}

func (n *natsBroker) String() string {
	return "nats"
}

func (n *natsBroker) setOption(opts ...Option) {
	for _, o := range opts {
		o(&n.opts)
	}

	n.Once.Do(func() {
		n.nopts = nats.GetDefaultOptions()
	})

	if nopts, ok := n.opts.Context.Value(optionsKey{}).(nats.Options); ok {
		n.nopts = nopts
	}

	// Options have higher priority than nats.Options
	// only if Addrs, Secure or TLSConfig were not set through a Option
	// we read them from nats.Option
	if len(n.opts.Addrs) == 0 {
		n.opts.Addrs = n.nopts.Servers
	}

	if !n.opts.Secure {
		n.opts.Secure = n.nopts.Secure
	}

	if n.opts.TLSConfig == nil {
		n.opts.TLSConfig = n.nopts.TLSConfig
	}
	n.addrs = n.setAddrs(n.opts.Addrs)

	if n.opts.Context.Value(drainConnectionKey{}) != nil {
		n.drain = true
		n.closeCh = make(chan error)
		n.nopts.ClosedCB = n.onClose
		n.nopts.AsyncErrorCB = n.onAsyncError
		n.nopts.DisconnectedErrCB = n.onDisconnectedError
	}
}

func (n *natsBroker) onClose(conn *nats.Conn) {
	n.closeCh <- nil
}

func (n *natsBroker) onAsyncError(conn *nats.Conn, sub *nats.Subscription, err error) {
	// There are kinds of different async error nats might callback, but we are interested
	// in ErrDrainTimeout only here.
	if err == nats.ErrDrainTimeout {
		n.closeCh <- err
	}
}

func (n *natsBroker) onDisconnectedError(conn *nats.Conn, err error) {
	n.closeCh <- err
}

func NewNatsBroker(opts ...Option) Broker {
	options := Options{
		// addrs: opts,
		// Default codec
		Context: context.Background(),
	}

	n := &natsBroker{
		opts:  options,
		codec: json.Marshaler{},
	}
	n.setOption(opts...)

	return n
}
