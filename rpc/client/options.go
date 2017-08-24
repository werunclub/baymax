package client

import (
	"os"
	"time"
)

var (
	DefaultRegistry      = "consul" // consul or etcd
	DefaultEtcdAddress   = "http://127.0.0.1:2379"
	DefaultConsulAddress = "127.0.0.1:8500"
	DefaultPoolSize      = 100
)

type Option func(*Options)

type Options struct {
	Registry string

	// EtcdAddress 地址用于注册服务
	EtcdAddress []string

	// Consul 地址用于注册服务
	ConsulAddress string

	SessionTimeout time.Duration

	// 连接超时时长
	ConnTimeout time.Duration

	// 连接池大小
	PoolSize int

	//　连接有效时长
	PoolTTL time.Duration

	// 重试次数
	Retries int
}

func newOptions(opt ...Option) Options {

	opts := Options{}

	for _, o := range opt {
		o(&opts)
	}

	opts.Registry = os.Getenv("RPC_REGISTRY")
	if opts.Registry == "" {
		opts.Registry = DefaultRegistry
	}

	if len(opts.EtcdAddress) == 0 {
		opts.EtcdAddress = []string{DefaultEtcdAddress}
	}

	if opts.ConsulAddress == "" {
		opts.ConsulAddress = DefaultConsulAddress
	}

	if opts.SessionTimeout == 0 {
		opts.SessionTimeout = time.Second * 5
	}

	if opts.ConnTimeout == 0 {
		opts.ConnTimeout = time.Second * 5
	}

	if opts.PoolSize == 0 {
		opts.PoolSize = DefaultPoolSize
	}

	if opts.PoolTTL == 0 {
		opts.PoolTTL = time.Minute * 30
	}

	if opts.Retries == 0 {
		opts.Retries = 3
	}

	return opts
}

// Namespace 即将废除：名称空间
func Namespace(n string) Option {
	return func(o *Options) {}
}

// 注册服务类型
func Registry(a string) Option {
	return func(o *Options) {
		o.Registry = a
	}
}

// EtcdAddress 地址
func EtcdAddress(a []string) Option {
	return func(o *Options) {
		o.EtcdAddress = a
	}
}

// Consul 地址
func ConsulAddress(a string) Option {
	return func(o *Options) {
		o.ConsulAddress = a
	}
}

// ConnTimeout 连接超时
func ConnTimeout(t time.Duration) Option {
	return func(o *Options) {
		o.ConnTimeout = t
	}
}

// PoolSize 连接池大小
func PoolSize(size int) Option {
	return func(o *Options) {
		o.PoolSize = size
	}
}

// Retries 重试次数
func Retries(times int) Option {
	return func(o *Options) {
		o.Retries = times
	}
}
