package registry

var (
	DefaultConsulAddress = "127.0.0.1:8500"
	DefaultPoolSize      = 100
)

type Option func(*Options)

type Options struct {

	// Consul地址用于注册服务
	ConsulAddress string

	// 初始服务列表
	ServiceNames []string

	// 连接池大小
	PoolSize int
}

func newOptions(opt ...Option) Options {

	opts := Options{}

	for _, o := range opt {
		o(&opts)
	}

	if len(opts.ConsulAddress) == 0 {
		opts.ConsulAddress = DefaultConsulAddress
	}

	if opts.PoolSize <= 0 {
		opts.PoolSize = DefaultPoolSize
	}

	return opts
}

// Consul 地址
func ConsulAddress(a string) Option {
	return func(o *Options) {
		o.ConsulAddress = a
	}
}

// 连接池大小
func PoolSize(size int) Option {
	return func(o *Options) {
		o.PoolSize = size
	}
}
