package pubsub

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"baymax/pubsub/broker"
	log "github.com/Sirupsen/logrus"
	"github.com/go-errors/errors"
)

type Server struct {
	broker broker.Broker
	exit   chan chan error

	sync.RWMutex
	subscribers map[*subscriber][]broker.Subscriber
}

func NewServer(addrs ...string) *Server {
	opt := broker.Addrs(addrs...)

	return &Server{
		broker:      broker.NewBroker(opt),
		subscribers: make(map[*subscriber][]broker.Subscriber),
		exit:        make(chan chan error),
	}
}

// 新建订阅器
func (s *Server) NewSubscriber(topic string, sb interface{}, opts ...SubscriberOption) Subscriber {
	return newSubscriber(topic, sb, opts...)
}

func (s *Server) Subscribe(sb Subscriber) error {
	sub, ok := sb.(*subscriber)
	if !ok {
		log.Error("invalid subscriber: expected *subscriber")
		return errors.New("invalid subscriber: expected *subscriber")
	}
	if len(sub.handlers) == 0 {
		log.Error("invalid subscriber: no handler functions")
		return errors.New("invalid subscriber: no handler functions")
	}

	if err := validateSubscriber(sb); err != nil {
		log.Errorf("Subscribe error:%f", err.Error())
		return err
	}

	s.Lock()
	_, ok = s.subscribers[sub]
	if ok {
		log.Errorf("subscriber %v already exists", s)
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
			log.Infof("Unsubscribing from topic: %s", sub.Topic())
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
	if err := s.Start(); err != nil {
		return err
	}

	if err := s.Register(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	log.Printf("Received signal %s", <-ch)

	if err := s.Deregister(); err != nil {
		return err
	}

	return s.Stop()
}
