package rpc

import (
	"net/rpc"
	"sync"
	"time"

	"baymax/errors"
	logger "github.com/Sirupsen/logrus"
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

	// Fixme: 无法连接到服务器时此处有空指针错误
	conn, e := c.pool.GetConn(c.Addr, c.timeout)
	if e != nil {
		logger.Errorf("rpc connect error:", e)
		return errors.InternalServerError(e.Error()).(*errors.Error)
	}

	var err error

	defer func() {
		// 使用后释放
		c.pool.release(c.Addr, conn, err)
	}()

	err = conn.Call(method, args, reply)
	if err != nil {
		logger.Errorf("rpc call error:", err)
		return errors.Parse(err.Error())
	}

	return nil
}
