package pubsub

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/werunclub/baymax/v2/log"
	"github.com/werunclub/baymax/v2/pubsub/broker"
	"github.com/werunclub/baymax/v2/pubsub/broker/nats"
	"github.com/werunclub/baymax/v2/pubsub/broker/nsq"
)

type Server struct {
	broker broker.Broker

	Exit chan bool

	sync.RWMutex
	subscribers map[*subscriber][]broker.Subscriber
}

func NewServer(addrs ...string) *Server {
	opts := broker.Addrs(addrs...)

	var b broker.Broker
	if brokerType == "nats" {
		b = nats.NewNatsBroker(opts)
	} else {
		b = nsq.NewNsqBroker(opts)
	}

	return &Server{
		broker:      b,
		subscribers: make(map[*subscriber][]broker.Subscriber),

		Exit: make(chan bool, 1),
	}
}

// 新建订阅器
func (s *Server) NewSubscriber(topic string, sb interface{}, opts ...SubscriberOption) Subscriber {
	return newSubscriber(topic, sb, opts...)
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
	for sb, _ := range s.subscribers {
		handler := s.createSubHandler(sb)
		var opts []broker.SubscribeOption

		opts = append(opts, broker.DisableAutoAck())
		if queue := sb.Options().Queue; len(queue) > 0 {
			opts = append(opts, broker.Queue(queue))
		}
		sub, err := s.broker.Subscribe(sb.Topic(), handler, opts...)
		if err != nil {
			return err
		}
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

func (s *Server) Start() error {
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

	if err := s.Start(); err != nil {
		log.SourcedLogrus().WithError(err).Errorf("pubsub start fail")
		panic("pubsub start fail")
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
