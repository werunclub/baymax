package server

import (
	"os"
	"strconv"
	"time"

	"github.com/pborman/uuid"
)

var (
	DefaultAddress       = ":0"
	DefaultName          = "go-server"
	DefaultNamespace     = "go-srv-"
	DefaultProtocol      = "tcp"
	DefaultVersion       = "1.0.0"
	DefaultConsulAddress = "127.0.0.1:8500"
)

type Option func(*Options)

type Options struct {
	Metadata map[string]string

	// 服务名称,一个RPC服务指定一个唯一名称,最好写到协议里
	Name string

	// 名称空间
	Namespace string

	// 监听地址, 默认端口 :0 即随机端口(推荐)
	Address string

	// 协议:　tcp or http
	RpcProtocol string

	// Consul地址用于注册服务默认
	ConsulAddress string

	// 公开地址
	Advertise string
	Id        string
	Version   string

	RegisterTTL      time.Duration
	RegisterInterval time.Duration

	// 健康检查开启
	CheckEnable bool

	StopWait int

	WriteTimeout time.Duration
	ReadTimeout  time.Duration
}

func newOptions(opt ...Option) Options {
	opts := Options{
		Metadata: map[string]string{},
	}

	for _, o := range opt {
		o(&opts)
	}

	if len(opts.Address) == 0 {
		opts.Address = DefaultAddress
	}

	if len(opts.Name) == 0 {
		opts.Name = DefaultName
	}

	if len(opts.Namespace) == 0 {
		opts.Namespace = DefaultNamespace
	}

	if opts.RpcProtocol != "tcp" && opts.RpcProtocol != "http" {
		opts.RpcProtocol = DefaultProtocol
	}

	if len(opts.Id) == 0 {
		opts.Id = uuid.NewUUID().String()

	}

	if len(opts.Version) == 0 {
		opts.Version = DefaultVersion
	}

	if len(opts.ConsulAddress) == 0 {
		opts.ConsulAddress = DefaultConsulAddress
	}

	if opts.RegisterTTL == 0 {
		opts.RegisterTTL = time.Second * 30
	}

	if opts.RegisterInterval == 0 {
		opts.RegisterInterval = time.Second * 10
	}

	if opts.StopWait <= 0 {
		wait := os.Getenv("RPC_STOP_WAIT")
		opts.StopWait, _ = strconv.Atoi(wait)
	}

	if opts.WriteTimeout <= 0 {
		opts.WriteTimeout = 5 * time.Second
	}

	if opts.ReadTimeout <= 0 {
		opts.ReadTimeout = 5 * time.Second
	}

	return opts
}

// 服务名称
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// 名称空间
func Namespace(n string) Option {
	return func(o *Options) {
		o.Namespace = n + "-"
	}
}

// 服务器唯一ID, 默认使用 uuid
func Id(id string) Option {
	return func(o *Options) {
		o.Id = id
	}
}

// 服务版本
func Version(v string) Option {
	return func(o *Options) {
		o.Version = v
	}
}

// 监听地址
func Address(a string) Option {
	return func(o *Options) {
		o.Address = a
	}
}

// 使用协议
func Protocol(a string) Option {
	return func(o *Options) {
		o.RpcProtocol = a
	}
}

// Consul 地址
func ConsulAddress(a string) Option {
	return func(o *Options) {
		o.ConsulAddress = a
	}
}

// 公开地址,用于注册服务
func Advertise(a string) Option {
	return func(o *Options) {
		o.Advertise = a
	}
}

// 服务有效时间
func RegisterTTL(t time.Duration) Option {
	return func(o *Options) {
		o.RegisterTTL = t
	}
}

// 上报间隔
func RegisterInterval(t time.Duration) Option {
	return func(o *Options) {
		o.RegisterInterval = t
	}
}

// 开启健检查
func CheckEnable(enable bool) Option {
	return func(o *Options) {
		o.CheckEnable = enable
	}
}

// 关闭服务前等待时间
func StopWait(wait int) Option {
	return func(o *Options) {
		o.StopWait = wait
	}
}

func WriteTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.WriteTimeout = timeout
	}
}

func ReadTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.ReadTimeout = timeout
	}
}
