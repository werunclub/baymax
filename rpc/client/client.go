package client

import (
	"context"
	"strings"

	"baymax/errors"
	"baymax/rpc/helpers"

	"github.com/Sirupsen/logrus"
	rpcxClient "github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
)

// Client represents a RPC client.
type Client struct {
	opts Options

	ServiceName string

	rpcClient rpcxClient.XClient
	discovery rpcxClient.ServiceDiscovery

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
	rpcxOption.ConnectTimeout = options.ConnTimeout
	rpcxOption.SerializeType = protocol.MsgPack

	client.discovery = rpcxClient.NewEtcdDiscovery(helpers.RPCRath, serviceName,
		options.EtcdAddress, nil)

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

func (c *Client) cleanServiceMethod(serviceMethod string) string {
	parts := strings.Split(serviceMethod, ".")

	if len(parts) >= 2 {
		return parts[len(parts)-1]
	}
	return serviceMethod
}

// Call 同步执行
func (c *Client) Call(serviceMethod string, args interface{}, reply interface{}) *errors.Error {
	return c.CallContext(context.Background(), serviceMethod, args, reply)
}

// Go 异步执行
func (c *Client) Go(serviceMethod string, args interface{}, reply interface{}, done chan *rpcxClient.Call) (*rpcxClient.Call, error) {
	return c.GoContext(context.Background(), serviceMethod, args, reply, done)
}

// CallContext 使用上下文同步执行
func (c *Client) CallContext(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) *errors.Error {
	err := c.rpcClient.Call(ctx, c.cleanServiceMethod(serviceMethod), args, reply)
	if err != nil {
		logrus.WithField("serviceMethod", serviceMethod).WithError(err).Errorf("rpc call fail")
		return errors.Parse(err.Error())
	}
	return nil
}

// GoContext 使用上下文异步执行
func (c *Client) GoContext(ctx context.Context, serviceMethod string, args interface{},
	reply interface{}, done chan *rpcxClient.Call) (*rpcxClient.Call, error) {
	return c.rpcClient.Go(ctx, c.cleanServiceMethod(serviceMethod), args, reply, done)
}
