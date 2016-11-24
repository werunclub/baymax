package client

import (
	"fmt"
	"io"
	"math"
	"net/rpc"
	"sync"
	"time"

	"baymax/errors"
	"baymax/log"
	"baymax/rpc/registry"
	"strings"
)

// Client represents a RPC client.
type Client struct {
	opts Options
	pool *pool

	once sync.Once

	ServiceName string
	Selector    *registry.Selector

	//重试次数
	Retries int
}

func NewClient(serviceName string, opts ...Option) *Client {

	options := newOptions(opts...)

	client := Client{
		opts: options,
		pool: newPool(options.PoolSize, options.PoolTTL),

		ServiceName: serviceName,
		Selector:    registry.NewSelector(registry.ConsulAddress(options.ConsulAddress)),
		Retries:     3,
	}

	// 初始化选择器
	client.Selector.AddServices(client.getServiceName())

	return &client
}

//　完整名称:　名称空间+名称
func (c *Client) getServiceName() string {
	return c.opts.Namespace + c.ServiceName
}

func (c *Client) SetPoolSize(size int) {
	c.pool.size = size
}

// 断开连接
func (c *Client) Close() error {
	return nil
}

// 调用方法
func (c *Client) call(network, address, method string, args interface{}, reply interface{}) error {

	// Fixme: 无法连接到服务器时此处有空指针错误
	conn, e := c.pool.Dial(network, address, c.opts.ConnTimeout)
	if e != nil {
		log.SourcedLogrus().WithError(e).Errorf("rpc connect error")
		return registry.ErrConnectIsLost
	}

	err := conn.Call(method, args, reply)
	if err != nil {
		log.SourcedLogrus().WithError(err).
			WithField("server", address).
			WithField("method", method).
			WithField("args", args).
			Errorf("Call rpc method error")
	}

	conn.Close()
	//c.pool.release(network, address, conn, err)

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
	next, err := c.Selector.Select(c.getServiceName())

	if err != nil && err == registry.ErrNotFound {
		log.SourcedLogrus().WithField("method", method).
			WithField("service", c.getServiceName()).
			WithError(err).Errorf("rpc service not found")

		return errors.Parse(errors.NotFound(err.Error()).Error())
	} else if err != nil {
		log.SourcedLogrus().WithField("method", method).
			WithField("service", c.getServiceName()).
			WithError(err).Errorf("get service nodes fail")

		return errors.Parse(errors.InternalServerError(err.Error()).Error())
	}

	call := func(i int) error {
		// 根据执行序号延迟执行
		if t, err := backoff(method, i); err != nil {
			return err
		} else if t.Seconds() > 0 {
			time.Sleep(t)
		}

		// 获取服务地址
		node, err := next()
		if err != nil {
			return err
		}

		address := node.Address
		if !strings.Contains(address, ":") && node.Port > 0 {
			address = fmt.Sprintf("%s:%d", address, node.Port)
		}

		var network string
		if strings.Contains(address, "@") {
			parts := strings.Split(address, "@")
			network = parts[0]
			address = parts[1]
		} else {
			network = "tcp"
		}

		// 调用rpc
		if err := c.call(network, address, method, args, reply); err != nil {

			if err == registry.ErrConnectIsLost {
				nodes, _ := c.Selector.GetNodes(c.getServiceName())

				log.SourcedLogrus().WithField("retry", i).
					WithField("network", network).
					WithField("address", address).
					WithField("method", method).
					WithField("args", args).
					WithField("servers", nodes).
					WithError(err).Errorf("call err %d times", i)

				c.Selector.Mark(c.ServiceName, node.Id, err)
			}
			return err
		}
		return nil
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
				err != registry.ErrNotFound &&
				err != registry.ErrNoneAvailable &&
				err != io.EOF &&
				err != io.ErrUnexpectedEOF &&
				err != registry.ErrConnectIsLost {

				// ErrShutdown ErrNotFound ErrNoneAvailable 需要重试的错误
				// 其它错误直接返回
				return errors.Parse(err.Error())
			}

			gerr = err
		}
	}

	if gerr != nil && gerr.Error() != "" {
		log.SourcedLogrus().WithField("method", method).
			WithError(gerr).Errorf("rpc call got system error")
		return errors.Parse(gerr.Error())
	}
	return nil
}
