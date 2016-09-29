package rpc

import (
	"github.com/go-errors/errors"
	"time"
)

type Selector struct {
	opts          Options
	ConsulAddress string

	sessionTimeout time.Duration
	SelectMode     SelectMode
	timeout        time.Duration
	connTimeout    time.Duration

	selectors map[string]*ConsulClientSelector
	clients   map[string]*Client
}

// 新建选择器
func NewSelector(opts ...Option) *Selector {
	options := newOptions(opts...)

	selector := &Selector{
		opts:           options,
		ConsulAddress:  options.ConsulAddress,
		sessionTimeout: time.Second * 30,
		timeout:        time.Second * 20,
		connTimeout:    time.Second * 5,

		// 随机选择模型
		SelectMode: RandomSelect,

		// 为每个服务建一个选择器
		selectors: make(map[string]*ConsulClientSelector),

		// 缓存客户端
		clients: make(map[string]*Client),
	}
	selector.AddServices(options.ServiceNames...)

	return selector
}

func (s *Selector) AddServices(serviceNames ...string) {
	for _, name := range serviceNames {
		s.addService(name)
	}
}

// 添加一个服务
// 选择器会定时自动从注册服务器获取可用服务器列表
func (s *Selector) addService(serviceName string) {
	_, ok := s.selectors[serviceName]
	if !ok {
		s.selectors[serviceName] = NewConsulClientSelector(s.ConsulAddress,
			serviceName,
			s.sessionTimeout,
			s.SelectMode,
			s.timeout)
	}
}

// 获取或新建一个选择器
func (s *Selector) getSelector(serviceName string) (*ConsulClientSelector, error) {

	selector, ok := s.selectors[serviceName]

	if !ok {
		s.addService(serviceName)
		selector, ok = s.selectors[serviceName]

		if !ok {
			return nil, errors.New("add service fail")
		}
	}

	return selector, nil
}

// 获取或新建一个客户端
func (s *Selector) getClient(address string) (*Client, error) {

	client, ok := s.clients[address]

	if !ok || client == nil {

		client = NewClient("tcp", address, s.connTimeout)
		client.SetPoolSize(s.opts.PoolSize)

		s.clients[address] = client
	}

	return client, nil
}

// 选择一个服务器,并创建客户端
func (s *Selector) Select(serviceName string) (*Client, error) {

	var (
		err      error
		node     *Node
		selector *ConsulClientSelector
	)

	selector, err = s.getSelector(serviceName)
	if err != nil {
		return nil, err
	}

	node, err = selector.Select()
	if err != nil {
		return nil, err
	}

	return s.getClient(node.Address)
}

// TODO: 标记服务器不可用
// TODO: 自动移除不可用服务器客户端
func (s *Selector) Mark(serviceName string, address string, err error) {
	delete(s.clients, address)
}
