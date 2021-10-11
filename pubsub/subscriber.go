package pubsub

import (
	"bytes"
	"fmt"
	"reflect"

	"golang.org/x/net/context"

	"github.com/werunclub/baymax/v2/log"
	"github.com/werunclub/baymax/v2/pubsub/broker"
	"github.com/werunclub/baymax/v2/pubsub/codec"
	mj "github.com/werunclub/baymax/v2/pubsub/codec/jsonrpc"
	"github.com/werunclub/baymax/v2/pubsub/metadata"
)

const (
	subSig = "func(context.Context, interface{}) error"
)

type handler struct {
	method  reflect.Value
	reqType reflect.Type
	ctxType reflect.Type
}

type subscriber struct {
	topic      string
	rcvr       reflect.Value
	typ        reflect.Type
	subscriber interface{}
	handlers   []*handler
	opts       SubscriberOptions
}

func newSubscriber(topic string, sub interface{}, opts ...SubscriberOption) Subscriber {
	options := SubscriberOptions{
		AutoAck: true,
	}

	for _, o := range opts {
		o(&options)
	}

	var handlers []*handler

	if typ := reflect.TypeOf(sub); typ.Kind() == reflect.Func {
		h := &handler{
			method: reflect.ValueOf(sub),
		}

		switch typ.NumIn() {
		case 1:
			h.reqType = typ.In(0)
		case 2:
			h.ctxType = typ.In(0)
			h.reqType = typ.In(1)
		}

		handlers = append(handlers, h)

	} else {
		for m := 0; m < typ.NumMethod(); m++ {
			method := typ.Method(m)
			h := &handler{
				method: method.Func,
			}

			switch method.Type.NumIn() {
			case 2:
				h.reqType = method.Type.In(1)
			case 3:
				h.ctxType = method.Type.In(1)
				h.reqType = method.Type.In(2)
			}

			handlers = append(handlers, h)
		}
	}

	return &subscriber{
		rcvr:       reflect.ValueOf(sub),
		typ:        reflect.TypeOf(sub),
		topic:      topic,
		subscriber: sub,
		handlers:   handlers,
		opts:       options,
	}
}

func validateSubscriber(sub Subscriber) error {
	typ := reflect.TypeOf(sub.Subscriber())
	var argType reflect.Type

	if typ.Kind() == reflect.Func {
		name := "Func"
		switch typ.NumIn() {
		case 2:
			argType = typ.In(1)
		default:
			return fmt.Errorf("subscriber %v takes wrong number of args: %v required signature %s", name, typ.NumIn(), subSig)
		}
		if !isExportedOrBuiltinType(argType) {
			return fmt.Errorf("subscriber %v argument type not exported: %v", name, argType)
		}
		if typ.NumOut() != 1 {
			return fmt.Errorf("subscriber %v has wrong number of outs: %v require signature %s",
				name, typ.NumOut(), subSig)
		}
		if returnType := typ.Out(0); returnType != typeOfError {
			return fmt.Errorf("subscriber %v returns %v not error", name, returnType.String())
		}
	} else {
		hdlr := reflect.ValueOf(sub.Subscriber())
		name := reflect.Indirect(hdlr).Type().Name()

		for m := 0; m < typ.NumMethod(); m++ {
			method := typ.Method(m)

			switch method.Type.NumIn() {
			case 3:
				argType = method.Type.In(2)
			default:
				return fmt.Errorf("subscriber %v.%v takes wrong number of args: %v required signature %s",
					name, method.Name, method.Type.NumIn(), subSig)
			}

			if !isExportedOrBuiltinType(argType) {
				return fmt.Errorf("%v argument type not exported: %v", name, argType)
			}
			if method.Type.NumOut() != 1 {
				return fmt.Errorf(
					"subscriber %v.%v has wrong number of outs: %v require signature %s",
					name, method.Name, method.Type.NumOut(), subSig)
			}
			if returnType := method.Type.Out(0); returnType != typeOfError {
				return fmt.Errorf("subscriber %v.%v returns %v not error", name, method.Name, returnType.String())
			}
		}
	}

	return nil
}

func (s *Server) createSubHandler(sb *subscriber) broker.Handler {
	return func(p broker.Publication) error {
		msg := p.Message()

		hdr := make(map[string]string)
		for k, v := range msg.Header {
			hdr[k] = v
		}
		ctx := metadata.NewContext(context.Background(), hdr)

		done := make(chan bool, len(sb.handlers))

		for i := 0; i < len(sb.handlers); i++ {
			handler := sb.handlers[i]

			var isVal bool
			var req reflect.Value

			if handler.reqType.Kind() == reflect.Ptr {
				req = reflect.New(handler.reqType.Elem())
			} else {
				req = reflect.New(handler.reqType)
				isVal = true
			}
			if isVal {
				req = req.Elem()
			}

			b := &buffer{bytes.NewBuffer(msg.Body)}
			co := mj.NewCodec(b)
			defer co.Close()

			if err := co.ReadHeader(&codec.Message{}, codec.Publication); err != nil {
				return err
			}

			if err := co.ReadBody(req.Interface()); err != nil {
				return err
			}

			fn := func(ctx context.Context, msg *publication, done chan bool) error {
				var vals []reflect.Value
				if sb.typ.Kind() != reflect.Func {
					vals = append(vals, sb.rcvr)
				}
				if handler.ctxType != nil {
					vals = append(vals, reflect.ValueOf(ctx))
				}
				vals = append(vals, reflect.ValueOf(msg.Message()))

				returnValues := handler.method.Call(vals)
				if err := returnValues[0].Interface(); err != nil {
					log.SourcedLogrus().WithField("topic", msg.topic).
						WithField("msg", msg).
						WithError(err.(error)).
						Errorf("call msg handler fail")
					done <- false
					return err.(error)
				}
				done <- true
				return nil
			}

			go fn(ctx, &publication{
				topic:   sb.topic,
				message: req.Interface(),
			}, done)
		}

		var (
			finished int
			failures int
		)

		for {
			success := <-done
			finished++
			if !success {
				failures++
			}

			if finished == len(sb.handlers) {
				break
			}
		}

		if failures == 0 && sb.opts.AutoAck {
			return p.Ack()
		}

		return nil
	}
}

func (s *subscriber) Topic() string {
	return s.topic
}

func (s *subscriber) Subscriber() interface{} {
	return s.subscriber
}

func (s *subscriber) Options() SubscriberOptions {
	return s.opts
}
