package pubsub

import "context"

type SubscriberOption func(*SubscriberOptions)

// Subscriber interface represents a subscription to a given topic using
// a specific subscriber function or object with methods.
type Subscriber interface {
	Topic() string
	Subscriber() interface{}
	Options() SubscriberOptions
}

type SubscriberOptions struct {
	Queue   string
	AutoAck bool

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// Shared queue name distributed messages across subscribers
func SubscriberQueue(n string) SubscriberOption {
	return func(o *SubscriberOptions) {
		o.Queue = n
	}
}

func DisableAutoAck() SubscriberOption {
	return func(o *SubscriberOptions) {
		o.AutoAck = false
	}
}
