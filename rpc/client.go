package rpc

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"
	"time"
	//"baymax/errors"
	"baymax/errors"
)

// Client represents a RPC client.
type Client struct {
	rpcClient *rpc.Client
	net       string
	addr      string
	timeout   time.Duration
	once      sync.Once
}

func NewClient(net, addr string, timeout time.Duration) *Client {
	return &Client{
		net:     net,
		addr:    addr,
		timeout: timeout,
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
	return c.rpcClient.Close()
}

// 调用方法
func (c *Client) Call(method string, args interface{}, reply interface{}) *errors.Error {

	c.once.Do(func() {
		c.Connect()
	})

	err := c.rpcClient.Call(method, args, reply)

	if err != nil {
		return errors.Parse(err.Error())
	}

	return nil
}
