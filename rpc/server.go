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
	opts   Options
	nodeId string

	sync.RWMutex
	rpcServer  *rpc.Server
	Registry   *ConsulRegistry
	listener   net.Listener
	registered bool

	ticker *time.Ticker
}

func NewServer(opts ...Option) *Server {
	options := newOptions(opts...)

	server := &Server{
		opts:      options,
		rpcServer: rpc.NewServer(),
		Registry:  NewConsulRegistry(),
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

	// 注册服务
	node := &Node{
		Id:       config.Id,
		Name:     config.Name,
		Address:  addr + ":" + strconv.Itoa(port),
		Metadata: config.Metadata,
	}

	s.Registry.Register(node)
	s.nodeId = node.Id

	// 按指定时间上报状态
	s.ticker = time.NewTicker(s.opts.RegisterInterval)
	go func() {
		for range s.ticker.C {
			s.Registry.CheckPass(s.nodeId)
		}
	}()

	return nil
}

func (s *Server) Deregister() error {
	s.ticker.Stop()
	s.Registry.Unregister(s.nodeId)
	return nil
}

// 使用名称注册处理器
func (s *Server) Handle(name string, service interface{}) {
	s.rpcServer.RegisterName(name, service)
}

// alias func Handle
func (s *Server) RegisterName(name string, service interface{}) {
	s.Handle(name, service)
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
		return err
	}

	if err := s.Start(); err != nil {
		return err
	}

	// 注册服务
	if err := s.Register(); err != nil {
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
