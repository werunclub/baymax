package client

import (
	"io"
	"net/rpc"
	"sync"
	"time"

	"baymax/errors"
	"baymax/log"
	"baymax/rpc/registry"
)

// Client represents a RPC client.
type DirectClient struct {
	rpcClient *rpc.Client
	net       string
	Addr      string
	pool      *pool

	timeout time.Duration
	once    sync.Once

	retries int
}

func NewDirectClient(net, addr string, timeout time.Duration) *DirectClient {
	return &DirectClient{
		net:     net,
		Addr:    addr,
		timeout: timeout,
		pool:    newPool(100, time.Minute*30),
		retries: 3,
	}
}

func (c *DirectClient) SetPoolSize(size int) {
	c.pool.size = size
}

// 断开连接
func (c *DirectClient) Close() error {
	return nil
}

// 调用方法
// TODO: 优化连接池
func (c *DirectClient) Call(method string, args interface{}, reply interface{}) *errors.Error {

	call := func(i int) error {

		// 根据执行序号延迟执行
		if t, err := backoff(method, i); err != nil {
			return err
		} else if t.Seconds() > 0 {
			time.Sleep(t)
		}

		// Fixme: 无法连接到服务器时此处有空指针错误
		conn, e := c.pool.GetConn(c.net, c.Addr, c.timeout)
		if e != nil {
			log.SourcedLogrus().WithField("method", method).Errorf("rpc connect error:", e)
			return e
		}

		var err error

		defer func() {
			// 使用后释放
			c.pool.release(c.net, c.Addr, conn, err)
		}()

		err = conn.Call(method, args, reply)
		return err
	}

	ch := make(chan error, c.retries)

	var gerr error

	for i := 0; i < c.retries; i++ {
		go func() {
			ch <- call(i)
		}()

		select {
		case err := <-ch:
			// call 成功即刻返回
			if err == nil {
				return nil
			} else if err != rpc.ErrShutdown &&
				err != registry.ErrNotFound &&
				err != registry.ErrNoneAvailable &&
				err != io.EOF {

				// ErrShutdown ErrNotFound ErrNoneAvailable 需要重试的错误
				// 其它错误直接返回
				log.SourcedLogrus().WithField("method", method).WithError(err).Debugf("rpc call fail")
				return errors.Parse(err.Error())
			}

			gerr = err
		}
	}

	if gerr != nil {
		log.SourcedLogrus().WithField("method", method).WithError(gerr).Debugf("rpc call fail")
	}

	return errors.Parse(gerr.Error())
}
