package rpcx

import (
	"github.com/smallnest/rpcx"
	"net/rpc/jsonrpc"
	"time"
)

type Client struct {
	rpcxClient *rpcx.Client
}

func NewClient(net, addr string, timeout time.Duration) *Client {

	s := &rpcx.DirectClientSelector{
		Network: net,
		Address: addr,
		Timeout: timeout,
	}
	client := rpcx.NewClient(s)
	client.ClientCodecFunc = jsonrpc.NewClientCodec

	return &Client{
		rpcxClient: client,
	}
}

// 断开连接
func (c *Client) Close() error {
	return c.rpcxClient.Close()
}

// 调用方法
func (c *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	return c.rpcxClient.Call(serviceMethod, args, reply)
}
