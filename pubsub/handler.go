package pubsub

type SubscriberOption func(*SubscriberOptions)

// Subscriber interface represents a subscription to a given topic using
// a specific subscriber function or object with methods.
type Subscriber interface {
	Topic() string
	Subscriber() interface{}
	Options() SubscriberOptions
}

type HandlerOptions struct {
	Internal bool
	Metadata map[string]map[string]string
}

type SubscriberOptions struct {
	Queue    string
	Internal bool
}

// Internal Subscriber options specifies that a subscriber is not advertised
// to the discovery system.
func InternalSubscriber(b bool) SubscriberOption {
	return func(o *SubscriberOptions) {
		o.Internal = b
	}
}

// Shared queue name distributed messages across subscribers
func SubscriberQueue(n string) SubscriberOption {
	return func(o *SubscriberOptions) {
		o.Queue = n
	}
}
