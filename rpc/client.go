package rpc

import (
	"fmt"
	"io"
	"math"
	"net/rpc"
	"sync"
	"time"

	"baymax/errors"
	"baymax/log"
)

// Client represents a RPC client.
type Client struct {
	pool    *pool
	timeout time.Duration
	once    sync.Once

	ServiceName string
	Selector    *Selector

	//重试次数
	Retries int
}

func NewClient(serviceName, consulAddress string, timeout time.Duration) *Client {
	return &Client{
		timeout: timeout,
		pool:    newPool(100, time.Minute*30),

		ServiceName: serviceName,
		Selector:    NewSelector(ConsulAddress(consulAddress)),
		Retries:     3,
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
func (c *Client) call(address, method string, args interface{}, reply interface{}) error {

	// Fixme: 无法连接到服务器时此处有空指针错误
	conn, e := c.pool.GetConn(address, c.timeout)
	if e != nil {
		log.SourcedLogrus().WithError(e).Errorf("rpc connect error")
		return e
	}

	var err error

	defer func() {
		// 使用后释放
		c.pool.release(address, conn, err)
	}()

	err = conn.Call(method, args, reply)
	return err
}

// exponential backoff
func backoff(method string, attempts int) (time.Duration, error) {
	if attempts == 0 {
		return time.Duration(0), nil
	}
	return time.Duration(math.Pow(10, float64(attempts))) * time.Millisecond, nil
}

// 调用RPC方法
func (c *Client) Call(method string, args interface{}, reply interface{}) *errors.Error {

	// 获取一个服务地址选择器
	next, err := c.Selector.SelectNodes(c.ServiceName)
	if err != nil && err == ErrNotFound {
		log.SourcedLogrus().WithField("method", method).WithError(err).Debugf("rpc service not found")
		return errors.Parse(errors.NotFound(err.Error()).Error())
	} else if err != nil {
		log.SourcedLogrus().WithField("method", method).WithError(err).Debugf("rpc call fail")
		return errors.Parse(errors.InternalServerError(err.Error()).Error())
	}

	call := func(i int) error {
		// 根据执行序号延迟执行
		if t, err := backoff(method, i); err != nil {
			return errors.InternalServerError(err.Error())
		} else if t.Seconds() > 0 {
			time.Sleep(t)
		}

		// 获取服务地址
		node, err := next()
		if err != nil && err == ErrNotFound {
			return errors.NotFound(err.Error())
		} else if err != nil {
			return errors.InternalServerError(err.Error())
		}

		address := node.Address
		if node.Port > 0 {
			address = fmt.Sprintf("%s:%d", address, node.Port)
		}

		// 调用rpc
		err = c.call(address, method, args, reply)
		//c.Selector.Mark(c.ServiceName, address, err)
		return err
	}

	ch := make(chan error, c.Retries)
	var gerr error

	for i := 0; i < c.Retries; i++ {
		go func() {
			ch <- call(i)
		}()

		select {
		case err := <-ch:
			// 调用成功
			if err == nil {
				return nil

			} else if err != rpc.ErrShutdown &&
				err != ErrNotFound &&
				err != ErrNoneAvailable &&
				err != io.EOF &&
				err != io.ErrUnexpectedEOF {

				// ErrShutdown  ErrNotFound ErrNoneAvailable 需要重试的错误
				// 其它错误直接返回
				log.SourcedLogrus().WithField("method", method).WithError(err).Debugf("rpc call got error")
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
