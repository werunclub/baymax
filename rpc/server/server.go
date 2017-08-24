package server

import (
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx"
	"github.com/smallnest/rpcx/codec"
	"github.com/smallnest/rpcx/plugin"

	"baymax/log"
	"baymax/rpc/helpers"
)

type Server struct {
	opts      Options
	rpcServer *rpcx.Server

	registryConsul *plugin.ConsulRegisterPlugin
	registryEtcd   *plugin.EtcdV3RegisterPlugin

	Handlers map[string]interface{}

	listener net.Listener
	Exit     chan bool

	sync.RWMutex
	registered bool
	ticker     *time.Ticker
}

// NewServer 初始化rpc服务
func NewServer(opts ...Option) *Server {
	options := newOptions(opts...)

	server := &Server{
		opts:      options,
		rpcServer: rpcx.NewServer(),
		Handlers:  make(map[string]interface{}),

		Exit: make(chan bool, 1),
	}

	// 使用 JSON 编码
	server.rpcServer.ServerCodecFunc = codec.NewJSONRPCServerCodec

	if options.Registry == "etcd" {
		server.registryEtcd = &plugin.EtcdV3RegisterPlugin{
			EtcdServers:         options.EtcdAddress,
			BasePath:            "/rpcx",
			Metrics:             metrics.NewRegistry(),
			Services:            make([]string, 0),
			UpdateIntervalInSec: int64(server.opts.RegisterInterval.Seconds()),
		}

		server.rpcServer.PluginContainer.Add(server.registryEtcd)

	} else {
		server.registryConsul = &plugin.ConsulRegisterPlugin{
			ConsulAddress:  options.ConsulAddress,
			UpdateInterval: server.opts.RegisterTTL,
		}

		server.rpcServer.PluginContainer.Add(server.registryConsul)
	}

	if options.InfluxDBHost != "" {
		metrics := plugin.NewMetricsPlugin()
		metrics.InfluxDB(10e9, options.InfluxDBHost, options.InfluxDBDB,
			options.InfluxDBUser, options.InfluxDBPass)

		server.rpcServer.PluginContainer.Add(metrics)
	}

	return server
}

func (s *Server) setServiceAddress(addr string) {
	if s.opts.Registry == "etcd" {
		s.registryEtcd.ServiceAddress = addr
	} else {
		s.registryConsul.ServiceAddress = addr
	}
}

// Address 服务地址
func (s *Server) Address() string {
	return s.rpcServer.Address()
}

// Handle 注册服务
func (s *Server) Handle(serviceName string, service interface{}) {
	s.RegisterName(serviceName, service)
}

// RegisterName 注册服务
func (s *Server) RegisterName(serviceName string, service interface{}) {
	s.Handlers[serviceName] = service
}

// Register 将服务注册到服务注册发现服务器
func (s *Server) Register() error {

	s.opts.Address = s.Address()

	var advt, host string
	var port int

	// 优先使用 Advertise 地址注册
	// Advertise 用于对外公布地址, 比如在 docker 中运行需要外部服务调用时需要指定
	if len(s.opts.Advertise) > 0 {
		advt = s.opts.Advertise
	} else {
		advt = s.opts.Address
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

	// 设置服务地址
	s.setServiceAddress(s.opts.Protocol + "@" + addr + ":" + strconv.Itoa(port))

	for name, service := range s.Handlers {
		s.rpcServer.RegisterName(name, service)
	}

	s.registered = true
	s.RUnlock()

	if s.opts.Registry == "consul" {
		// 按指定时间上报状态
		// fixme: 上报状态失败考虑重新注册
		s.ticker = time.NewTicker(s.opts.RegisterInterval)
		go func() {
			for range s.ticker.C {
				for name := range s.Handlers {
					s.registryConsul.CheckPass(name)
				}
			}
		}()
	}

	return nil
}

// Deregister 注销服务
func (s *Server) Deregister() {
	if s.opts.Registry == "etcd" {
		for name := range s.Handlers {
			s.registryEtcd.Unregister(name)
		}
	} else {
		for name := range s.Handlers {
			s.registryConsul.Unregister(name)
		}
	}
}

//　开始服务
func (s *Server) start() error {

	if s.opts.Registry == "etcd" {
		if err := s.registryEtcd.Start(); err != nil {
			return err
		}
	} else {
		if err := s.registryConsul.Start(); err != nil {
			return err
		}
	}

	if s.opts.Protocol == "http" {
		ln, err := net.Listen("tcp", s.opts.Address)
		if err != nil {
			return err
		}
		s.listener = ln
		go s.rpcServer.ServeByHTTP(s.listener, rpcx.DefaultRPCPath)
	} else {
		s.rpcServer.Start("tcp", s.opts.Address)
	}

	return nil
}

// Stop 关闭连接
func (s *Server) Stop() error {
	return s.rpcServer.Close()
}

// RegisterAndRun 注册服务并运行
func (s *Server) RegisterAndRun() error {
	defer func() {
		log.SourcedLogrus().Printf("Rpc server exit.")
		s.Exit <- true
	}()

	if err := s.start(); err != nil {
		return err
	}
	s.Register()

	log.SourcedLogrus().Printf("Running on %s", s.Address())

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	log.SourcedLogrus().Printf("Received signal %s", <-ch)

	if s.ticker != nil {
		s.ticker.Stop()
	}

	// 注销服务
	s.Deregister()

	// 暂停10s
	if s.opts.StopWait > 0 {
		time.Sleep(time.Second * time.Duration(s.opts.StopWait))
	}

	// 关闭连接
	s.Stop()

	return nil
}
