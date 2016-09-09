package rpc

import (
	"net/rpc"
	"sync"
	"time"

	"baymax/errors"
)

// Client represents a RPC client.
type Client struct {
	rpcClient *rpc.Client
	net       string
	Addr      string

	pool    *pool
	timeout time.Duration
	once    sync.Once
}

func NewClient(net, addr string, timeout time.Duration) *Client {
	return &Client{
		net:     net,
		Addr:    addr,
		timeout: timeout,
		pool:    newPool(100, time.Minute*30),
	}
}

func (c *Client) SetPoolSize(size int) {
	c.pool.size = size
}

// 断开连接
func (c *Client) Close() error {
	return nil
}

// 调用方法
// TODO: 优化连接池
func (c *Client) Call(method string, args interface{}, reply interface{}) *errors.Error {

	conn, e := c.pool.GetConn(c.Addr, c.timeout)
	if e != nil {
		return errors.Parse(e.Error())
	}

	var grr error

	defer func() {
		// 使用后释放
		c.pool.release(c.Addr, conn, grr)
	}()

	err := conn.Call(method, args, reply)
	if err != nil {
		grr = err
		return errors.Parse(err.Error())
	}

	return nil
}
