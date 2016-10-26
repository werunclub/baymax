package registry

import (
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/pborman/uuid"
	"strconv"
)

type Node struct {
	Id       string
	Name     string
	Address  string
	Port     int
	Version  string
	Metadata map[string]string `json:"metadata"`
}

func NewNode(name, address, version string) *Node {
	return &Node{
		Id:      uuid.New(),
		Name:    name,
		Address: address,
		Version: version,
	}
}

//ConsulRegisterPlugin a register plugin which can register services into consul for cluster
//This registry is experimental and has not been test.
type ConsulRegistry struct {
	ConsulAddress  string
	consulConfig   *api.Config
	client         *api.Client
	Services       []string
	UpdateInterval time.Duration

	CheckEnable bool
}

func NewConsulRegistry() *ConsulRegistry {
	return &ConsulRegistry{}
}

func (c *ConsulRegistry) Init() (err error) {
	if c.consulConfig == nil {
		c.consulConfig = api.DefaultConfig()
		c.consulConfig.Address = c.ConsulAddress
	}
	c.client, err = api.NewClient(c.consulConfig)

	if err != nil {
		return err
	}

	return
}

//Close closes this plugin
func (c *ConsulRegistry) Close() {
}

// Register handles registering event.
func (c *ConsulRegistry) Register(node *Node) (err error) {

	check := api.AgentServiceCheck{}

	if c.CheckEnable {
		check = api.AgentServiceCheck{
			TTL:    strconv.Itoa(int(c.UpdateInterval.Seconds())) + "s",
			Status: api.HealthPassing,
			TCP:    node.Address,
		}
	}

	service := &api.AgentServiceRegistration{
		ID:      node.Id,
		Name:    node.Name,
		Address: node.Address,
		Port:    node.Port,
		Tags:    []string{node.Version},
		Check:   &check,
	}
	agent := c.client.Agent()
	err = agent.ServiceRegister(service)
	return
}

// Unregister a service from consul but this service still exists in this node.
func (c *ConsulRegistry) Unregister(nodeId string) {
	agent := c.client.Agent()
	id := nodeId
	agent.ServiceDeregister(id)
}

// CheckPass sets check pass
func (c *ConsulRegistry) CheckPass(nodeId string) error {
	agent := c.client.Agent()
	id := nodeId
	return agent.UpdateTTL("service:"+id, "", api.HealthPassing)
}

// CheckFail sets check fail
func (c *ConsulRegistry) CheckFail(nodeId string) error {
	agent := c.client.Agent()
	id := nodeId
	return agent.UpdateTTL("service:"+id, "", api.HealthCritical)
}

func (c *ConsulRegistry) GetService(name string) ([]*Node, error) {
	rsp, _, err := c.client.Health().Service(name, "", true, nil)
	if err != nil {
		return nil, err
	}

	var nodes []*Node

	for _, s := range rsp {
		if s.Service.Service != name {
			continue
		}

		// service ID is now the node id
		id := s.Service.ID

		// address is service address
		address := s.Service.Address

		node := &Node{
			Id:      id,
			Name:    name,
			Address: address,
			Port:    s.Service.Port,
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

// Name return name of this plugin.
func (c *ConsulRegistry) Name() string {
	return "ConsulRegistry"
}

// Description return description of this plugin.
func (c *ConsulRegistry) Description() string {
	return "a register plugin which can register services into etcd for cluster"
}
