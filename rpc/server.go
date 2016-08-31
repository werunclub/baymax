package rpc

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"baymax/broker"
)

type Server struct {
	rpcServer   *rpc.Server
	listener    net.Listener
	subscribers map[*subscriber][]broker.Subscriber
}

func NewServer() *Server {
	return &Server{
		rpcServer: rpc.NewServer(),
	}
}

// 启动服务
func (s *Server) Serve(network, address string) {

	ln, err := net.Listen(network, address)
	if err != nil {
		return
	}

	s.listener = ln
	for {
		c, err := ln.Accept()
		if err != nil {
			continue
		}
		go s.rpcServer.ServeCodec(jsonrpc.NewServerCodec(c))
	}
}

// 使用名称注册服务
func (s *Server) RegisterName(name string, service interface{}) {
	s.rpcServer.RegisterName(name, service)
}

// 关闭服务
func (s *Server) Close() error {
	return s.listener.Close()
}

func (s *Server) NewSubscriber(topic string, sb interface{}, opts ...SubscriberOption) Subscriber {
	return newSubscriber(topic, sb, opts...)
}

func (s *rpcServer) Subscribe(sb Subscriber) error {
	sub, ok := sb.(*subscriber)
	if !ok {
		return fmt.Errorf("invalid subscriber: expected *subscriber")
	}
	if len(sub.handlers) == 0 {
		return fmt.Errorf("invalid subscriber: no handler functions")
	}

	if err := validateSubscriber(sb); err != nil {
		return err
	}

	s.Lock()
	_, ok = s.subscribers[sub]
	if ok {
		return fmt.Errorf("subscriber %v already exists", s)
	}
	s.subscribers[sub] = nil
	s.Unlock()
	return nil
}
