package pubsub

type publication struct {
	topic   string
	message interface{}
}

func newPublication(topic string, message interface{}) publication {
	return publication{
		message: message,
		topic:   topic,
	}
}

func (r *publication) Topic() string {
	return r.topic
}

func (r *publication) Message() interface{} {
	return r.message
}
