package client

import (
	"context"

	"baymax/errors"

	"github.com/Sirupsen/logrus"
	"github.com/smallnest/rpcx"
	"github.com/smallnest/rpcx/clientselector"
	"github.com/smallnest/rpcx/codec"
	"github.com/smallnest/rpcx/core"
)

// Client represents a RPC client.
type Client struct {
	opts Options

	ServiceName string

	rpcClient *rpcx.Client
	Selector  *rpcx.ClientSelector

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

	if options.Registry == "etcd" {
		selector := clientselector.NewEtcdV3ClientSelector(
			options.EtcdAddress,
			"/rpcx/"+serviceName,
			options.SessionTimeout,
			rpcx.RandomSelect,
			options.ConnTimeout,
		)

		client.rpcClient = rpcx.NewClient(selector)
	} else {
		selector := clientselector.NewConsulClientSelector(
			options.ConsulAddress,
			serviceName,
			options.SessionTimeout,
			rpcx.RandomSelect,
			options.ConnTimeout,
		)

		client.rpcClient = rpcx.NewClient(selector)
	}

	// 使用 JSON 编码
	client.rpcClient.ClientCodecFunc = codec.NewJSONRPCClientCodec

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
func (c *Client) Go(serviceMethod string, args interface{}, reply interface{}, done chan *core.Call) *core.Call {
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
	reply interface{}, done chan *core.Call) *core.Call {
	return c.rpcClient.Go(ctx, serviceMethod, args, reply, done)
}
