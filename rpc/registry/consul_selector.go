package registry

import (
	"math/rand"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"

	"baymax/rpc/helpers"
)

// SelectMode defines the algorithm of selecting a services from cluster
type SelectMode int

const (
	RandomSelect SelectMode = iota
	RoundRobinSelect
	LeastActiveSelect
	ConsistentHashSeelct
)

var selectModeStrs = [...]string{
	"RandomSelect",
	"RoundRobin",
	"LeastActive",
	"ConsistentHash",
}

func (s SelectMode) String() string {
	return selectModeStrs[s]
}

// ConsulClientSelector is used to select a rpc server from consul.
//This registry is experimental and has not been test.
type ConsulClientSelector struct {
	ConsulAddress      string
	consulConfig       *api.Config
	client             *api.Client
	ticker             *time.Ticker
	sessionTimeout     time.Duration
	Servers            []*Node
	SelectMode         SelectMode
	timeout            time.Duration
	rnd                *rand.Rand
	currentServer      int
	len                int
	HashServiceAndArgs helpers.HashServiceAndArgs
	serviceName        string

	sync.RWMutex
}

// NewConsulClientSelector creates a ConsulClientSelector
func NewConsulClientSelector(consulAddress string, serviceName string, sessionTimeout time.Duration,
	sm SelectMode, timeout time.Duration) *ConsulClientSelector {

	selector := &ConsulClientSelector{
		ConsulAddress:  consulAddress,
		Servers:        make([]*Node, 1),
		sessionTimeout: sessionTimeout,
		SelectMode:     sm,
		timeout:        timeout,
		rnd:            rand.New(rand.NewSource(time.Now().UnixNano())),
		serviceName:    serviceName,
	}

	selector.init()
	return selector
}

func (s *ConsulClientSelector) SetSelectMode(sm SelectMode) {
	s.SelectMode = sm
}

func (s *ConsulClientSelector) init() {
	if s.consulConfig == nil {
		s.consulConfig = api.DefaultConfig()
		s.consulConfig.Address = s.ConsulAddress
	}
	s.client, _ = api.NewClient(s.consulConfig)

	s.pullServers()

	s.ticker = time.NewTicker(s.sessionTimeout)
	go func() {
		for range s.ticker.C {
			s.pullServers()
		}
	}()
}

func (s *ConsulClientSelector) pullServers() {
	rsp, _, err := s.client.Health().Service(s.serviceName, "", true, nil)
	if err != nil {
		return
	}

	var nodes []*Node

	s.RLock()
	for _, r := range rsp {
		if r.Service.Service != s.serviceName {
			continue
		}

		// 服务器ID
		id := r.Service.ID

		// 服务器地址
		address := r.Service.Address

		node := &Node{
			Id:      id,
			Name:    s.serviceName,
			Address: address,
			Port:    r.Service.Port,
		}
		nodes = append(nodes, node)
	}

	s.Servers = nodes
	s.len = len(s.Servers)

	s.RUnlock()
}

// CheckFail sets check fail
func (c *ConsulClientSelector) CheckFail(nodeId string) error {
	agent := c.client.Agent()
	id := nodeId
	return agent.UpdateTTL("service:"+id, "", api.HealthCritical)
}

// 从已有服务器列表中选择服务器
func (s *ConsulClientSelector) Select(options ...interface{}) (Next, error) {

	if s.len == 0 {
		s.pullServers()
	}

	if s.len == 0 {
		return nil, ErrNoneAvailable
	}

	return Random(s.Servers), nil
}

// fixme: 是否可靠
func (s *ConsulClientSelector) Mark(nodeId string, err error) {

	index := -1
	for i, server := range s.Servers {
		if server.Id == nodeId {
			index = i
		}
	}

	// 标记不可用
	s.CheckFail(nodeId)

	if index >= 0 {
		s.RLock()
		s.Servers = append(s.Servers[:index], s.Servers[index+1:]...)
		s.RUnlock()
	}
}
