package client

import (
	"context"

	"baymax/errors"

	"github.com/Sirupsen/logrus"
	rpcxClient "github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
)

const (
	BasePath = "/rpcx"
)

// Client represents a RPC client.
type Client struct {
	opts Options

	ServiceName string

	rpcClient *rpcxClient.XClient
	discovery *rpcxClient.ServiceDiscovery

	//重试次数
	Retries int
}

// NewClient 初始化客户端
func NewClient(serviceName string, opts ...Option) *Client {

	options := newOptions(opts...)

	client := Client{
		opts:        options,
		ServiceName: serviceName,
		Retries:     options.Retries,
	}

	rpcxOption := rpcxClient.DefaultOption
	rpcxOption.ConnTimeout = options.ConnTimeout
	rpcxOption.SerializeType = protocol.JSON

	if options.Registry == "etcd" {
		client.discovery = rpcxClient.NewEtcdDiscovery(BasePath, serviceName, []string{options.EtcdAddress}, rpcxOption)

	} else {
		client.discovery = rpcxClient.NewConsulDiscovery(BasePath, serviceName, []string{options.ConsulAddress}, rpcxOption)
	}

	client.rpcClient = rpcxClient.NewXClient(
		serviceName,
		rpcxClient.Failtry,
		rpcxClient.RandomSelect,
		client.discovery,
		rpcxOption,
	)

	return &client
}

func (c *Client) getServiceName() string {
	return c.ServiceName
}

// Call 同步执行
func (c *Client) Call(serviceMethod string, args interface{}, reply interface{}) *errors.Error {
	return c.CallContext(context.Background(), serviceMethod, args, reply)
}

// Go 异步执行
func (c *Client) Go(serviceMethod string, args interface{}, reply interface{}, done chan *rpcxClient.Call) *rpcxClient.Call {
	return c.GoContext(context.Background(), serviceMethod, args, reply, done)
}

// CallContext 使用上下文同步执行
func (c *Client) CallContext(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) *errors.Error {
	err := c.rpcClient.Call(ctx, serviceMethod, args, reply)
	if err != nil {
		logrus.WithField("serviceMethod", serviceMethod).WithError(err).Errorf("rpc call fail")
		return errors.Parse(err.Error())
	}
	return nil
}

// GoContext 使用上下文异步执行
func (c *Client) GoContext(ctx context.Context, serviceMethod string, args interface{},
	reply interface{}, done chan *rpcxClient.Call) *rpcxClient.Call {
	return c.rpcClient.Go(ctx, serviceMethod, args, reply, done)
}
