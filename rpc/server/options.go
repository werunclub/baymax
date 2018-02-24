package server

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pborman/uuid"
)

var (
	DefaultAddress     = ":0"
	DefaultName        = "go-server"
	DefaultProtocol    = "tcp"
	DefaultVersion     = "1.0.0"
	DefaultEtcdAddress = []string{"127.0.0.1:2379"}
)

type Option func(*Options)

type Options struct {
	Metadata map[string]string

	// 服务名称,一个RPC服务指定一个唯一名称,最好写到协议里
	Name string

	// 监听地址, 默认端口 :0 即随机端口(推荐)
	Address string

	// 协议:　tcp or http
	Protocol string

	Registry string

	// Etcd 地址用于注册服务
	EtcdAddress []string

	// 公开地址
	Advertise string
	ID        string
	Version   string

	RegisterTTL      time.Duration
	RegisterInterval time.Duration

	WriteTimeout time.Duration
	ReadTimeout  time.Duration

	// 健康检查开启
	CheckEnable bool

	StopWait int

	InfluxDBHost string
	InfluxDBDB   string
	InfluxDBUser string
	InfluxDBPass string
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

	if opts.Protocol != "tcp" && opts.Protocol != "http" {
		opts.Protocol = DefaultProtocol
	}

	if len(opts.ID) == 0 {
		opts.ID = uuid.NewUUID().String()
	}

	if len(opts.Version) == 0 {
		opts.Version = DefaultVersion
	}

	if len(opts.EtcdAddress) == 0 {
		opts.EtcdAddress = DefaultEtcdAddress
	}

	if opts.RegisterTTL == 0 {
		opts.RegisterTTL = time.Second * 30
	}

	if opts.RegisterInterval == 0 {
		opts.RegisterInterval = time.Second * 5
	}

	if opts.WriteTimeout <= 0 {
		opts.WriteTimeout = 5 * time.Second
	}

	if opts.ReadTimeout <= 0 {
		opts.ReadTimeout = 5 * time.Second
	}

	if opts.StopWait <= 0 {
		wait := os.Getenv("RPC_STOP_WAIT")
		opts.StopWait, _ = strconv.Atoi(wait)
		if opts.StopWait == 0 {
			opts.StopWait = 10
		}
	}

	if opts.InfluxDBHost == "" {
		opts.InfluxDBHost = os.Getenv("INFLUX_DB_HOST")
	}

	if opts.InfluxDBDB == "" {
		opts.InfluxDBDB = os.Getenv("INFLUX_DB_DB")
	}

	if opts.InfluxDBUser == "" {
		opts.InfluxDBUser = os.Getenv("INFLUX_DB_USER")
	}

	if opts.InfluxDBPass == "" {
		opts.InfluxDBPass = os.Getenv("INFLUX_DB_PASS")
	}

	envEtcdAddress := os.Getenv("REGISTRY_ETCD_ADDRESS")
	if envEtcdAddress != "" {
		addrs := strings.Split(envEtcdAddress, ",")
		opts.EtcdAddress = addrs
	}

	return opts
}

// Name 服务名称
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// ID 服务器唯一ID, 默认使用 uuid
func ID(id string) Option {
	return func(o *Options) {
		o.ID = id
	}
}

// Namespace 即将废除：名称空间
func Namespace(n string) Option {
	return func(o *Options) {}
}

// Version 服务版本
func Version(v string) Option {
	return func(o *Options) {
		o.Version = v
	}
}

// Address 监听地址
func Address(a string) Option {
	return func(o *Options) {
		o.Address = a
	}
}

// Protocol 使用协议 http or tcp
func Protocol(a string) Option {
	return func(o *Options) {
		o.Protocol = a
	}
}

// Registry 注册服务类型
// deprecated
func Registry(a string) Option {
	return func(o *Options) {
		o.Registry = a
	}
}

// ConsulAddress consulAddress
// deprecated
func ConsulAddress(addr string) Option {
	return func(o *Options) {
	}
}

// EtcdAddress 地址
func EtcdAddress(a []string) Option {
	return func(o *Options) {
		o.EtcdAddress = a
	}
}

// Advertise 公开地址,用于注册服务
func Advertise(a string) Option {
	return func(o *Options) {
		o.Advertise = a
	}
}

// RegisterTTL 服务有效时间
func RegisterTTL(t time.Duration) Option {
	return func(o *Options) {
		o.RegisterTTL = t
	}
}

// RegisterInterval 上报间隔
func RegisterInterval(t time.Duration) Option {
	return func(o *Options) {
		o.RegisterInterval = t
	}
}

// WriteTimeout 写入超时
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

// CheckEnable 开启健检查
func CheckEnable(enable bool) Option {
	return func(o *Options) {
		o.CheckEnable = enable
	}
}

// StopWait 关闭服务前等待时间
func StopWait(wait int) Option {
	return func(o *Options) {
		o.StopWait = wait
	}
}
