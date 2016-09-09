package rpc

import (
	"errors"
	"math/rand"
	"time"

	"github.com/hashicorp/consul/api"
)

// SelectMode defines the algorithm of selecting a services from cluster
type SelectMode int

const (
	RandomSelect SelectMode = iota
	RoundRobin
	LeastActive
	ConsistentHash
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
	HashServiceAndArgs HashServiceAndArgs
	serviceName        string
}

// NewConsulClientSelector creates a ConsulClientSelector
func NewConsulClientSelector(consulAddress string, serviceName string, sessionTimeout time.Duration, sm SelectMode, timeout time.Duration) *ConsulClientSelector {
	selector := &ConsulClientSelector{
		ConsulAddress:  consulAddress,
		Servers:        make([]*Node, 1),
		sessionTimeout: sessionTimeout,
		SelectMode:     sm,
		timeout:        timeout,
		rnd:            rand.New(rand.NewSource(time.Now().UnixNano())),
		serviceName:    serviceName}

	selector.init()
	return selector
}

func (s *ConsulClientSelector) SetSelectMode(sm SelectMode) {
	s.SelectMode = sm
}

func (s *ConsulClientSelector) AllClients() []*Client {
	var clients []*Client

	for _, sv := range s.Servers {
		c := NewClient("tcp", sv.Address, s.timeout)
		if c != nil {
			clients = append(clients, c)
		}
	}

	return clients
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
		}
		nodes = append(nodes, node)
	}

	s.Servers = nodes
	s.len = len(s.Servers)
}

// 从已有服务器列表中选择服务器
func (s *ConsulClientSelector) Select(options ...interface{}) (*Node, error) {

	if s.len == 0 {
		return nil, errors.New("no valid server")
	}

	if s.SelectMode == RandomSelect {
		s.currentServer = s.rnd.Intn(s.len)
		server := s.Servers[s.currentServer]
		return server, nil

	} else if s.SelectMode == RandomSelect {
		s.currentServer = (s.currentServer + 1) % s.len
		server := s.Servers[s.currentServer]
		return server, nil

	} else if s.SelectMode == ConsistentHash {
		if s.HashServiceAndArgs == nil {
			s.HashServiceAndArgs = JumpConsistentHash
		}
		s.currentServer = s.HashServiceAndArgs(s.len, options)
		server := s.Servers[s.currentServer]
		return server, nil
	}

	return nil, errors.New("not supported SelectMode: " + s.SelectMode.String())
}

// TODO: 标记出错服务器
func (s *ConsulClientSelector) Mark(nodeId string, err error) {
}
