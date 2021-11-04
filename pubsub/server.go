package pubsub

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"

	"github.com/werunclub/baymax/v2/log"
	"github.com/werunclub/baymax/v2/pubsub/broker"
	"github.com/werunclub/baymax/v2/pubsub/codec"
	mj "github.com/werunclub/baymax/v2/pubsub/codec/jsonrpc"
	"github.com/werunclub/baymax/v2/pubsub/metadata"
	"golang.org/x/net/context"

	"github.com/go-errors/errors"
)

// 新建订阅器
func NewSubscriber(topic string, sub interface{}, opts ...SubscriberOption) Subscriber {
	return newSubscriber(topic, sub, opts...)
}

type Server struct {
	broker     broker.Broker
	registered bool

	Exit chan bool

	sync.RWMutex
	subscribers map[*subscriber][]broker.Subscriber
}

func NewServer(addrs ...string) *Server {
	opts := broker.Addrs(addrs...)

	return &Server{
		broker:      broker.NewBroker(opts),
		subscribers: make(map[*subscriber][]broker.Subscriber),

		Exit: make(chan bool, 1),
	}
}

// 新建订阅器
// Deprecated:  使用 pubsub.NewSubscriber
func (s *Server) NewSubscriber(topic string, sub interface{}, opts ...SubscriberOption) Subscriber {
	return NewSubscriber(topic, sub, opts...)
}

func (s *Server) Subscribe(sb Subscriber) error {
	sub, ok := sb.(*subscriber)
	if !ok {
		log.SourcedLogrus().Error("invalid subscriber: expected *subscriber")
		return errors.New("invalid subscriber: expected *subscriber")
	}
	if len(sub.handlers) == 0 {
		log.SourcedLogrus().Error("invalid subscriber: no handler functions")
		return errors.New("invalid subscriber: no handler functions")
	}

	if err := validateSubscriber(sb); err != nil {
		log.SourcedLogrus().Errorf("Subscribe error: %s", err.Error())
		return err
	}

	s.Lock()
	_, ok = s.subscribers[sub]
	if ok {
		log.SourcedLogrus().Errorf("subscriber %v already exists", s)
		return errors.New(fmt.Sprintf("subscriber %v already exists", s))
	}
	s.subscribers[sub] = nil
	s.Unlock()
	return nil
}

func (s *Server) Register() error {
	// subscribe for all of the subscribers
	for sb := range s.subscribers {
		var opts []broker.SubscribeOption
		if queue := sb.Options().Queue; len(queue) > 0 {
			opts = append(opts, broker.Queue(queue))
		}

		if cx := sb.Options().Context; cx != nil {
			opts = append(opts, broker.SubscribeContext(cx))
		}

		if !sb.Options().AutoAck {
			opts = append(opts, broker.DisableAutoAck())
		}

		sub, err := s.broker.Subscribe(sb.Topic(), createSubHandler(sb), opts...)
		if err != nil {
			return err
		}
		log.SourcedLogrus().Infof("Subscribing to topic: %s", sub.Topic())
		s.subscribers[sb] = []broker.Subscriber{sub}
	}

	return nil
}

func (s *Server) Deregister() error {
	for sb, subs := range s.subscribers {
		for _, sub := range subs {
			log.SourcedLogrus().Infof("Unsubscribing from topic: %s", sub.Topic())
			sub.Unsubscribe()
		}
		s.subscribers[sb] = nil
	}
	return nil
}

func (s *Server) Connect() error {
	return s.broker.Connect()
}

func (s *Server) Stop() error {
	return s.broker.Disconnect()
}

// Run starts the default server and waits for a kill
// signal before exiting. Also registers/deregisters the server
func (s *Server) Run() error {
	defer func() {
		s.Exit <- true
	}()

	if err := s.Connect(); err != nil {
		log.SourcedLogrus().WithError(err).Errorf("pubsub connect fail")
		panic("pubsub connect fail")
	}

	if err := s.Register(); err != nil {
		log.SourcedLogrus().WithError(err).Errorf("pubsub register fail")
		panic("pubsub register fail")
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	log.SourcedLogrus().Printf("Received signal %s", <-ch)

	if err := s.Deregister(); err != nil {
		log.SourcedLogrus().Errorf("Deregister fail")
	}

	s.Stop()

	log.SourcedLogrus().Printf("exit.")

	return nil
}

func createSubHandler(sb *subscriber) broker.Handler {
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
					log.SourcedLogrus().WithField("topic", msg.topic).WithField("msg", msg).WithError(err.(error)).Errorf("msg handle fail")
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
