package client

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
)

// Client represents a RPC client.
type Client struct {
	rpcClient *rpc.Client
	net       string
	addr      string
	timeout   time.Duration
}

func NewClient(net, addr string) *Client {
	return &Client{
		net:     net,
		addr:    addr,
		timeout: time.Second * 10,
	}
}

// 建立连接
func (c *Client) Connect() error {

	if c.rpcClient != nil {
		return nil
	}

	var conn net.Conn
	var err error

	conn, err = net.DialTimeout(c.net, c.addr, c.timeout)

	if err != nil {
		return err
	}

	c.rpcClient = jsonrpc.NewClient(conn)

	return nil
}

// 断开连接
func (c *Client) Close() error {
	if c.rpcClient != nil {
		return c.rpcClient.Close()
	}
	return nil
}

// 调用方法
func (c *Client) Call(method string, args interface{}, reply interface{}) error {

	var err error

	if err = c.Connect(); err != nil {
		return err
	}

	err = c.rpcClient.Call(method, args, reply)

	return err
}
