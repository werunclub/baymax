package pubsub

type rpcRequest struct {
	service     string
	method      string
	contentType string
	request     interface{}
	stream      bool
}

type rpcPublication struct {
	topic   string
	message interface{}
}

func (r *rpcRequest) Service() string {
	return r.service
}

func (r *rpcRequest) Method() string {
	return r.method
}

func (r *rpcRequest) Request() interface{} {
	return r.request
}

func (r *rpcRequest) Stream() bool {
	return r.stream
}

func (r *rpcPublication) Topic() string {
	return r.topic
}

func (r *rpcPublication) Message() interface{} {
	return r.message
}
