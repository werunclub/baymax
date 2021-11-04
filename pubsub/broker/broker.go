package broker

// Broker is an interface used for asynchronous messaging.
// Its an abstraction over various message brokers
// {NATS, RabbitMQ, Kafka, ...}
type Broker interface {
	Init(...Option) error
	Options() Options
	Address() string
	Connect() error
	Disconnect() error

	Publish(topic string, m *Message, opts ...PublishOption) error
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
// 	DefaultBroker Broker = NewNatsBroker()
// )

func NewBroker(opts ...Option) Broker {
	return NewNatsBroker(opts...)
}

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
