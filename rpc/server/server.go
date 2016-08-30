package server

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Server struct {
	rpcServer *rpc.Server
	listener  net.Listener
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

// 注册服务
func (s *Server) Register(service interface{}) {
	s.rpcServer.Register(service)
}

// 使用名称注册服务
func (s *Server) RegisterName(name string, service interface{}) {
	s.rpcServer.RegisterName(name, service)
}

// 关闭服务
func (s *Server) Close() error {
	return s.listener.Close()
}
