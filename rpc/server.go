package rpc

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
	"strings"
	"sync"

	"github.com/pborman/uuid"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	DefaultAddress       = ":0"
	DefaultName          = "go-server"
	DefaultVersion       = "1.0.0"
	DefaultId            = uuid.NewUUID().String()
	DefaultConsulAddress = "127.0.0.1:8500"
)

type Option func(*Options)

type Server struct {
	opts Options

	sync.RWMutex
	rpcServer *rpc.Server
	Registry  *ConsulRegistry
	listener  net.Listener

	Handlers   map[string]*interface{}
	nodes      []*Node
	registered bool

	ticker *time.Ticker
}

func NewServer(opts ...Option) *Server {
	options := newOptions(opts...)

	server := &Server{
		opts:      options,
		rpcServer: rpc.NewServer(),
		Registry:  NewConsulRegistry(),
		Handlers:  make(map[string]*interface{}),
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

// 启动服务
func (s *Server) Serve() {

	ln, err := net.Listen("tcp", s.opts.Address)
	if err != nil {
		return
	}

	s.listener = ln
	s.opts.Address = s.Address()

	for {
		c, err := ln.Accept()
		if err != nil {
			continue
		}
		go s.rpcServer.ServeCodec(jsonrpc.NewServerCodec(c))
	}
}

// 使用协程启动服务
func (s *Server) Start() error {

	ln, err := net.Listen("tcp", s.opts.Address)
	if err != nil {
		return err
	}

	s.listener = ln
	s.opts.Address = s.Address()

	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				continue
			}
			go s.rpcServer.ServeCodec(jsonrpc.NewServerCodec(c))
		}
	}()

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

	addr, err := extractAddress(host)
	if err != nil {
		return err
	}

	s.RLock()

	for name, _ := range s.Handlers {
		node := &Node{
			Id:       name + "-" + uuid.New(),
			Name:     name,
			Address:  addr + ":" + strconv.Itoa(port),
			Metadata: config.Metadata,
		}

		s.nodes = append(s.nodes, node)

		// 注册服务
		if err := s.Registry.Register(node); err != nil {
			return err
		}

		// 按指定时间上报状态
		s.ticker = time.NewTicker(s.opts.RegisterInterval)
		go func() {
			for range s.ticker.C {
				s.Registry.CheckPass(node.Id)
			}
		}()
	}

	s.RUnlock()

	return nil
}

func (s *Server) Deregister() error {
	s.ticker.Stop()

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

	// 启动注册服务
	if err := s.Registry.Init(); err != nil {
		log.Fatalf("registry init error: %v", err)
		return err
	}

	if err := s.Start(); err != nil {
		log.Fatalf("start error: %v", err)
		return err
	}

	// 注册服务
	if err := s.Register(); err != nil {
		log.Fatalf("registr error: %v", err)
		return err
	}

	log.Printf("Running on %s", s.Address())

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	log.Printf("Received signal %s", <-ch)

	// 取消注册
	if err := s.Deregister(); err != nil {
		return err
	}

	s.Registry.Close()
	return s.Stop()
}
