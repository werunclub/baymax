// Package redis provides a Redis broker
package broker

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/werunclub/baymax/v2/pubsub/codec"
	"github.com/werunclub/baymax/v2/pubsub/codec/json"
)

var (
	DefaultMaxActive      = 0
	DefaultMaxIdle        = 5
	DefaultIdleTimeout    = 2 * time.Minute
	DefaultConnectTimeout = 5 * time.Second
	DefaultReadTimeout    = 5 * time.Second
	DefaultWriteTimeout   = 5 * time.Second
)

// options contain additional options for the broker.
type redisBrokerOptions struct {
	maxIdle        int
	maxActive      int
	idleTimeout    time.Duration
	connectTimeout time.Duration
	readTimeout    time.Duration
	writeTimeout   time.Duration
}

// publication is an internal publication for the Redis
type redisPublication struct {
	topic   string
	message *Message
	err     error
}

// Topic returns the topic this publication applies to.
func (p *redisPublication) Topic() string {
	return p.topic
}

// Message returns the broker message of the publication.
func (p *redisPublication) Message() *Message {
	return p.message
}

// Ack sends an acknowledgement to the  However this is not supported
// is Redis and therefore this is a no-op.
func (p *redisPublication) Ack() error {
	return nil
}

func (p *redisPublication) Error() error {
	return p.err
}

// subscriber proxies and handles Redis messages as broker publications.
type redisSubscriber struct {
	codec  codec.Marshaler
	conn   *redis.PubSubConn
	topic  string
	handle Handler
	opts   SubscribeOptions
}

// recv loops to receive new messages from Redis and handle them
// as publications.
func (s *redisSubscriber) recv() {
	// Close the connection once the subscriber stops receiving.
	defer s.conn.Close()

	for {
		switch x := s.conn.Receive().(type) {
		case redis.Message:
			var m Message

			// Handle error? Only a log would be necessary since this type
			// of issue cannot be fixed.
			if err := s.codec.Unmarshal(x.Data, &m); err != nil {
				break
			}

			p := redisPublication{
				topic:   x.Channel,
				message: &m,
			}

			// Handle error? Retry?
			if p.err = s.handle(&p); p.err != nil {
				break
			}

			// Added for posterity, however Ack is a no-op.
			if s.opts.AutoAck {
				if err := p.Ack(); err != nil {
					break
				}
			}

		case redis.Subscription:
			if x.Count == 0 {
				return
			}

		case error:
			return
		}
	}
}

// Options returns the subscriber options.
func (s *redisSubscriber) Options() SubscribeOptions {
	return s.opts
}

// Topic returns the topic of the subscriber.
func (s *redisSubscriber) Topic() string {
	return s.topic
}

// Unsubscribe unsubscribes the subscriber and frees the connection.
func (s *redisSubscriber) Unsubscribe() error {
	return s.conn.Unsubscribe()
}

// broker implementation for Redis.
type redisBroker struct {
	addr  string
	pool  *redis.Pool
	opts  Options
	bopts *redisBrokerOptions

	codec codec.Marshaler
}

// String returns the name of the broker implementation.
func (b *redisBroker) String() string {
	return "redis"
}

// Options returns the options defined for the
func (b *redisBroker) Options() Options {
	return b.opts
}

// Address returns the address the broker will use to create new connections.
// This will be set only after Connect is called.
func (b *redisBroker) Address() string {
	return b.addr
}

// Init sets or overrides broker options.
func (b *redisBroker) Init(opts ...Option) error {
	if b.pool != nil {
		return errors.New("redis: cannot init while connected")
	}

	for _, o := range opts {
		o(&b.opts)
	}

	return nil
}

// Connect establishes a connection to Redis which provides the
// pub/sub implementation.
func (b *redisBroker) Connect() error {
	if b.pool != nil {
		return nil
	}

	var addr string

	if len(b.opts.Addrs) == 0 || b.opts.Addrs[0] == "" {
		addr = "redis://127.0.0.1:6379"
	} else {
		addr = b.opts.Addrs[0]

		if !strings.HasPrefix("redis://", addr) {
			addr = "redis://" + addr
		}
	}

	b.addr = addr

	b.pool = &redis.Pool{
		MaxIdle:     b.bopts.maxIdle,
		MaxActive:   b.bopts.maxActive,
		IdleTimeout: b.bopts.idleTimeout,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(
				b.addr,
				redis.DialConnectTimeout(b.bopts.connectTimeout),
				redis.DialReadTimeout(b.bopts.readTimeout),
				redis.DialWriteTimeout(b.bopts.writeTimeout),
			)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return nil
}

// Disconnect closes the connection pool.
func (b *redisBroker) Disconnect() error {
	err := b.pool.Close()
	b.pool = nil
	b.addr = ""
	return err
}

// Publish publishes a message.
func (b *redisBroker) Publish(topic string, msg *Message, opts ...PublishOption) error {
	v, err := b.codec.Marshal(msg)
	if err != nil {
		return err
	}

	conn := b.pool.Get()
	_, err = redis.Int(conn.Do("PUBLISH", topic, v))
	conn.Close()

	return err
}

// Subscribe returns a subscriber for the topic and handler.
func (b *redisBroker) Subscribe(topic string, handler Handler, opts ...SubscribeOption) (Subscriber, error) {
	var options SubscribeOptions
	for _, o := range opts {
		o(&options)
	}

	s := redisSubscriber{
		codec:  b.codec,
		conn:   &redis.PubSubConn{Conn: b.pool.Get()},
		topic:  topic,
		handle: handler,
		opts:   options,
	}

	// Run the receiver routine.
	go s.recv()

	if err := s.conn.Subscribe(s.topic); err != nil {
		return nil, err
	}

	return &s, nil
}

// NewRedisBroker returns a new broker implemented using the Redis pub/sub
// protocol. The connection address may be a fully qualified IANA address such
// as: redis://user:secret@localhost:6379/0?foo=bar&qux=baz
func NewRedisBroker(opts ...Option) Broker {
	// Default options.
	bopts := &redisBrokerOptions{
		maxIdle:        DefaultMaxIdle,
		maxActive:      DefaultMaxActive,
		idleTimeout:    DefaultIdleTimeout,
		connectTimeout: DefaultConnectTimeout,
		readTimeout:    DefaultReadTimeout,
		writeTimeout:   DefaultWriteTimeout,
	}

	// Initialize with empty broker options.
	options := Options{
		Context: context.WithValue(context.Background(), optionsKey{}, bopts),
	}

	for _, o := range opts {
		o(&options)
	}

	return &redisBroker{
		opts:  options,
		bopts: bopts,
		codec: json.Marshaler{},
	}
}
