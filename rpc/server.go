package rpc

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Server struct {
	rpcServer  *rpc.Server
	listener   net.Listener
	registered bool
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

// 使用协程启动服务
func (s *Server) Start(network, address string) {

	ln, err := net.Listen(network, address)
	if err != nil {
		return
	}

	s.listener = ln
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				continue
			}
			go s.rpcServer.ServeCodec(jsonrpc.NewServerCodec(c))
		}
	}()
}

// 使用名称注册服务
func (s *Server) RegisterName(name string, service interface{}) {
	s.rpcServer.RegisterName(name, service)
}

// 关闭服务
func (s *Server) Close() error {
	return s.listener.Close()
}

func (s *Server) Address() string {
	return s.listener.Addr().String()
}
