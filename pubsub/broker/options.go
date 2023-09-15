package broker

import (
	"context"
	"crypto/tls"
)

type Options struct {
	Name      string
	Addrs     []string
	Secure    bool
	TLSConfig *tls.Config

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context

	// Handler executed when error happens in broker mesage
	// processing
	ErrorHandler Handler
}

type PublishOptions struct {
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type SubscribeOptions struct {
	// AutoAck defaults to true. When a handler returns
	// with a nil error the message is acked.
	AutoAck bool
	// Subscribers with the same queue name
	// will create a shared subscription where each
	// receives a subset of messages.
	Queue string

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type Option func(*Options)

type PublishOption func(*PublishOptions)

type SubscribeOption func(*SubscribeOptions)

type contextKeyT string

func newSubscribeOptions(opts ...SubscribeOption) SubscribeOptions {
	opt := SubscribeOptions{
		AutoAck: true,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// name  broker
func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

// Addrs sets the host addresses to be used by the broker
func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

// DisableAutoAck will disable auto acking of messages
// after they have been handled.
func DisableAutoAck() SubscribeOption {
	return func(o *SubscribeOptions) {
		o.AutoAck = false
	}
}

// Queue sets the name of the queue to share messages on
func Queue(name string) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Queue = name
	}
}

// Secure communication with the broker
func Secure(b bool) Option {
	return func(o *Options) {
		o.Secure = b
	}
}

// Specify TLS Config
func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}

func NewOptions(opts ...Option) *Options {
	options := Options{
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	return &options
}
