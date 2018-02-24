package client

import (
	"time"
)

var (
	DefaultEtcdAddress = "127.0.0.1:2379"
	DefaultPoolSize    = 100
)

type Option func(*Options)

type Options struct {
	Registry string

	// EtcdAddress 地址用于注册服务
	EtcdAddress []string

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

	if len(opts.EtcdAddress) == 0 {
		opts.EtcdAddress = []string{DefaultEtcdAddress}
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

// ConsulAddress consulAddress
// deprecated
func ConsulAddress(addr string) Option {
	return func(o *Options) {
	}
}
