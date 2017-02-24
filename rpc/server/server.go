package server

import (
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"

	"baymax/rpc/helpers"
	"baymax/rpc/registry"
)

type Server struct {
	opts Options

	sync.RWMutex
	rpcServer *rpc.Server
	Registry  *registry.ConsulRegistry
	listener  net.Listener

	Handlers   map[string]*interface{}
	nodes      []*registry.Node
	registered bool

	ticker *time.Ticker

	Exit chan bool
}

func NewServer(opts ...Option) *Server {
	options := newOptions(opts...)

	server := &Server{
		opts:      options,
		rpcServer: rpc.NewServer(),
		Registry:  registry.NewConsulRegistry(),
		Handlers:  make(map[string]*interface{}),

		Exit: make(chan bool, 1),
	}

	server.Registry.ConsulAddress = options.ConsulAddress
	server.Registry.UpdateInterval = options.RegisterTTL
	return server
}

func (s *Server) Options() Options {
	s.RLock()
	opts := s.opts
	s.RUnlock()
	return opts
}

// 使用协程启动服务
func (s *Server) Start() error {
	if s.opts.RpcProtocol == "http" {
		return s.serveHttp()
	} else {
		return s.serveTcp()
	}
}

// via tcp
func (s *Server) serveTcp() error {

	ln, err := net.Listen("tcp", s.opts.Address)
	if err != nil {
		return err
	}

	s.listener = ln
	s.opts.Address = s.Address()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.WithError(err).Warnf("Error: accept rpc connection")
				continue
			}

			// 设置超时时间
			if s.opts.ReadTimeout > 0 {
				conn.SetReadDeadline(time.Now().Add(s.opts.ReadTimeout))
			}
			if s.opts.WriteTimeout > 0 {
				conn.SetWriteDeadline(time.Now().Add(s.opts.WriteTimeout))
			}

			go s.rpcServer.ServeCodec(jsonrpc.NewServerCodec(conn))

			//go func(conn net.Conn) {
			//	srv := jsonrpc.NewServerCodec(conn)
			//
			//	if err := s.rpcServer.ServeRequest(srv); err != nil {
			//		log.WithError(err).Errorf("Error: server rpc request")
			//	}
			//
			//	srv.Close()
			//}(conn)
		}
	}()

	return nil
}

// via http
func (s *Server) serveHttp() error {

	s.rpcServer.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

	ln, err := net.Listen("tcp", s.opts.Address)
	if err != nil {
		log.Fatal("listen error:", err)
	}

	s.listener = ln
	s.opts.Address = s.Address()

	go http.Serve(ln, nil)

	return nil
}

// 将服务注册到服务注册发现服务器
func (s *Server) Register() error {

	config := s.Options()
	var advt, host string
	var port int

	// 优先使用 Advertise 地址注册
	// Advertise 用于对外公布地址, 比如在 docker 中运行需要外部服务调用时需要指定
	if len(config.Advertise) > 0 {
		advt = config.Advertise
	} else {
		advt = config.Address
	}

	parts := strings.Split(advt, ":")
	if len(parts) > 1 {
		host = strings.Join(parts[:len(parts)-1], ":")
		port, _ = strconv.Atoi(parts[len(parts)-1])
	} else {
		host = parts[0]
	}

	addr, err := helpers.ExtractAddress(host)
	if err != nil {
		return err
	}

	s.RLock()

	for name, _ := range s.Handlers {
		node := &registry.Node{
			Id:       s.opts.Namespace + name + "@" + addr + ":" + strconv.Itoa(port),
			Name:     s.opts.Namespace + name,
			Address:  s.opts.RpcProtocol + "@" + addr + ":" + strconv.Itoa(port),
			Port:     port,
			Metadata: config.Metadata,
			Version:  "1",
		}

		// 注册服务
		if err := s.Registry.Register(node); err != nil {

			// 注销已注册服务
			s.Deregister()
			return err
		}

		// 注册成功添加到节点列表
		s.nodes = append(s.nodes, node)
	}

	s.registered = true
	s.RUnlock()

	// 按指定时间上报状态
	// fixme: 上报状态失败考虑重新注册
	s.ticker = time.NewTicker(s.opts.RegisterInterval)
	go func() {
		for range s.ticker.C {
			fails := 0
			for _, node := range s.nodes {
				if err := s.Registry.CheckPass(node.Id); err != nil {
					fails++
				}
			}
		}
	}()

	return nil
}

// 注销服务
func (s *Server) Deregister() error {
	if s.ticker != nil {
		s.ticker.Stop()
	}

	for _, node := range s.nodes {
		s.Registry.Unregister(node.Id)
	}

	return nil
}

// 使用名称注册处理器
func (s *Server) Handle(serviceName string, service interface{}) {
	s.rpcServer.RegisterName(serviceName, service)
	s.Handlers[serviceName] = &service
}

// alias func Handle
func (s *Server) RegisterName(serviceName string, service interface{}) {
	s.Handle(serviceName, service)
}

// 关闭连接
func (s *Server) Stop() error {
	return s.listener.Close()
}

func (s *Server) Address() string {
	return s.listener.Addr().String()
}

// 注册服务并运行
func (s *Server) RegisterAndRun() error {
	defer func() {
		s.Exit <- true
	}()

	// 启动注册服务
	if err := s.Registry.Init(); err != nil {
		log.Panicf("registry init error: %v", err)
		return err
	}

	if err := s.Start(); err != nil {
		log.Panicf("start error: %v", err)
		return err
	}

	// 注册服务
	if err := s.Register(); err != nil {
		log.Panicf("registr error: %v", err)
		return err
	}

	log.Printf("Running on %s", s.Address())

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	log.Printf("Received signal %s", <-ch)

	// 取消注册
	if err := s.Deregister(); err != nil {
		log.Errorf("rpc server deregister fail")
	}

	// 暂停10s
	if s.opts.StopWait > 0 {
		time.Sleep(time.Second * time.Duration(s.opts.StopWait))
	}

	s.Registry.Close()
	s.Stop()

	log.Printf("Rpc server exit.")

	return nil
}
