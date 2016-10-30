package registry

import (
	"time"

	"github.com/go-errors/errors"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrNoneAvailable = errors.New("none available")
	ErrConnectIsLost = errors.New("connect is lost")
)

// Next is a function that returns the next node
// based on the selector's strategy
type Next func() (*Node, error)

// Strategy is a selection strategy e.g random, round robin
type Strategy func([]*Node) Next

type Selector struct {
	opts Options

	sessionTimeout time.Duration
	SelectMode     SelectMode
	timeout        time.Duration
	connTimeout    time.Duration

	selectors map[string]*ConsulClientSelector
}

// 新建选择器
func NewSelector(opts ...Option) *Selector {

	options := newOptions(opts...)

	selector := &Selector{
		opts:           options,
		sessionTimeout: time.Second * 20,
		timeout:        time.Second * 20,
		connTimeout:    time.Second * 5,

		// 随机选择模型
		SelectMode: RandomSelect,

		// 为每个服务建一个选择器
		selectors: make(map[string]*ConsulClientSelector),
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
	if _, ok := s.selectors[serviceName]; !ok {
		s.selectors[serviceName] = NewConsulClientSelector(
			s.opts.ConsulAddress,
			serviceName,
			s.sessionTimeout,
			s.SelectMode,
			s.timeout,
		)
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

// 选择一个服务器,并创建客户端
func (s *Selector) Select(serviceName string) (Next, error) {

	selector, err := s.getSelector(serviceName)
	if err != nil {
		return nil, err
	}

	return selector.Select()
}

// fixme: 是否可靠
func (s *Selector) Mark(serviceName string, nodeId string, err error) error {
	selector, gerr := s.getSelector(serviceName)
	if gerr != nil {
		return gerr
	}

	selector.Mark(nodeId, err)
	return nil
}
