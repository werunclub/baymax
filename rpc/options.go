package rpc

import (
	"time"
)

type Options struct {
	Metadata map[string]string

	// 服务名称,一个RPC服务指定一个唯一名称,最好写到协议里
	Name string

	// 监听地址, 默认端口 :0 即随机端口(推荐)
	Address string

	// Consul地址用于注册服务默认
	ConsulAddress string

	// 公开地址
	Advertise string
	Id        string
	Version   string

	RegisterTTL      time.Duration
	RegisterInterval time.Duration
	timeout          time.Duration

	ServiceNames []string

	// 连接池大小
	PoolSize int
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

	if len(opts.Id) == 0 {
		opts.Id = DefaultId
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

	if opts.PoolSize == 0 {
		opts.PoolSize = DefaultPoolSize
	}

	return opts
}

// 服务名称
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
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

// 登记服务名称列表
func ServiceNames(ns []string) Option {
	return func(o *Options) {
		o.ServiceNames = ns
	}
}

// 连接池大小
func PoolSize(size int) Option {
	return func(o *Options) {
		o.PoolSize = size
	}
}
