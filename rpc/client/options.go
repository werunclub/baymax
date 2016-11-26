package client

import (
	"time"
)

var (
	DefaultNamespace     = "go-srv-"
	DefaultConsulAddress = "127.0.0.1:8500"
	DefaultPoolSize      = 100
)

type Option func(*Options)

type Options struct {

	// 名称空间
	Namespace string

	// Consul地址用于注册服务默认
	ConsulAddress string

	// 会话时长
	SessionTimeout time.Duration

	// 连接超时时长
	ConnTimeout time.Duration

	// 连接池大小
	PoolSize int

	//　连接有效时长
	PoolTTL time.Duration
}

func newOptions(opt ...Option) Options {

	opts := Options{}

	for _, o := range opt {
		o(&opts)
	}

	if len(opts.Namespace) == 0 {
		opts.Namespace = DefaultNamespace
	}

	if len(opts.ConsulAddress) == 0 {
		opts.ConsulAddress = DefaultConsulAddress
	}

	if opts.SessionTimeout == 0 {
		opts.SessionTimeout = time.Second * 10
	}

	if opts.ConnTimeout == 0 {
		opts.ConnTimeout = time.Second * 5
	}

	if opts.PoolSize == 0 {
		opts.PoolSize = DefaultPoolSize
	}

	if opts.PoolTTL == 0 {
		opts.PoolTTL = time.Minute * 20
	}

	return opts
}

// 名称空间
func Namespace(n string) Option {
	return func(o *Options) {
		o.Namespace = n + "-"
	}
}

// Consul 地址
func ConsulAddress(a string) Option {
	return func(o *Options) {
		o.ConsulAddress = a
	}
}

// 连接超时
func ConnTimeout(t time.Duration) Option {
	return func(o *Options) {
		o.ConnTimeout = t
	}
}

// 连接池大小
func PoolSize(size int) Option {
	return func(o *Options) {
		o.PoolSize = size
	}
}
