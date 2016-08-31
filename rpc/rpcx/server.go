package rpcx

import (
	"github.com/smallnest/rpcx"
	"net/rpc/jsonrpc"
)

func init() {
}

type Server struct {
	rpcxServer *rpcx.Server
}

func NewServer() *Server {

	server := &Server{
		rpcxServer: rpcx.NewServer(),
	}
	server.rpcxServer.ServerCodecFunc = jsonrpc.NewServerCodec

	//plugin := &plugin.EtcdRegisterPlugin{
	//	ServiceAddress: "tcp@127.0.0.1:8972",
	//	EtcdServers:    []string{"http://127.0.0.1:2379"},
	//	BasePath:       "/rpcx",
	//	Metrics:        metrics.NewRegistry(),
	//	Services:       make([]string, 1),
	//	UpdateInterval: time.Minute,
	//}
	//
	//err := plugin.Start()
	//if err != nil {
	//	panic("rpc start fail")
	//}
	//
	//server.rpcxServer.PluginContainer.Add(plugin)

	return server
}

// 启动服务
func (s *Server) Serve(network, address string) {
	s.rpcxServer.Serve(network, address)
}

// 使用名称注册服务
func (s *Server) RegisterName(name string, service interface{}) {
	s.rpcxServer.RegisterName(name, service)
}

// 关闭服务
func (s *Server) Close() error {
	return s.rpcxServer.Close()
}
