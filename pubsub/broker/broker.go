package broker

// Broker is an interface used for asynchronous messaging.
// Its an abstraction over various message brokers
// {NATS, RabbitMQ, Kafka, ...}
type Broker interface {
	Options() Options
	Address() string
	Connect() error
	Disconnect() error
	Init(...Option) error
	Publish(string, *Message, ...PublishOption) error
	Subscribe(string, Handler, ...SubscribeOption) (Subscriber, error)
	String() string
}

// Handler is used to process messages via a subscription of a topic.
// The handler is passed a publication interface which contains the
// message and optional Ack method to acknowledge receipt of the message.
type Handler func(Publication) error

type Message struct {
	Header map[string]string
	Body   []byte
}

// Publication is given to a subscription handler for processing
type Publication interface {
	Topic() string
	Message() *Message
	Ack() error
	Error() error
}

type Subscriber interface {
	Options() SubscribeOptions
	Topic() string
	Unsubscribe() error
}

// var (
// 	DefaultBroker Broker = NewBroker()
// )

// func NewBroker(opts ...Option) Broker {
// 	var options Options
// 	for _, o := range opts {
// 		o(&options)
// 	}

// 	if options.Name == "nats" {
// 		return nats.NewNatsBroker(opts...)
// 	}

// 	return nsq.NewNsqBroker(opts...)
// }

// func Init(opts ...Option) error {
// 	return DefaultBroker.Init(opts...)
// }

// func Connect() error {
// 	return DefaultBroker.Connect()
// }

// func Disconnect() error {
// 	return DefaultBroker.Disconnect()
// }

// func Publish(topic string, msg *Message, opts ...PublishOption) error {
// 	return DefaultBroker.Publish(topic, msg, opts...)
// }

// func Subscribe(topic string, handler Handler, opts ...SubscribeOption) (Subscriber, error) {
// 	return DefaultBroker.Subscribe(topic, handler, opts...)
// }

// func String() string {
// 	return DefaultBroker.String()
// }
