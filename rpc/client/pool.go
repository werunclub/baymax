package client

import (
	"io"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"
	"time"
)

type pool struct {
	size int
	ttl  int64

	sync.Mutex
	conns map[string][]*poolConn
}

type poolConn struct {
	*rpc.Client
	created int64
}

func newPool(size int, ttl time.Duration) *pool {
	return &pool{
		size:  size,
		ttl:   int64(ttl.Seconds()),
		conns: make(map[string][]*poolConn),
	}
}

// 建立连接
func (p *pool) dialTcp(addr string, timeout time.Duration) (*rpc.Client, error) {
	var (
		conn net.Conn
		err  error
	)

	conn, err = net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return nil, err
	}

	return jsonrpc.NewClient(conn), nil
}

//　建立　http　连接
func (p *pool) dialHTTP(addr string) (*rpc.Client, error) {
	return rpc.DialHTTP("tcp", addr)
}

//
func (p *pool) Dial(network, addr string, connTimeout time.Duration) (*rpc.Client, error) {
	if network == "http" {
		return p.dialHTTP(addr)
	} else {
		return p.dialTcp(addr, connTimeout)
	}
}

//　获取一个连接
func (p *pool) GetConn(network, addr string, connTimeout time.Duration) (*poolConn, error) {
	p.Lock()
	conns := p.conns[addr]
	now := time.Now().Unix()

	// 优化从连接池获取连接
	for len(conns) > 0 {
		conn := conns[len(conns)-1]
		conns = conns[:len(conns)-1]
		p.conns[addr] = conns

		// 关闭过期连接
		if d := now - conn.created; d > p.ttl {
			conn.Close()
			continue
		}

		p.Unlock()

		return conn, nil
	}
	p.Unlock()

	c, err := p.Dial(network, addr, connTimeout)
	if err != nil {
		return nil, err
	}

	return &poolConn{c, time.Now().Unix()}, nil
}

func (p *pool) release(net, addr string, conn *poolConn, err error) {

	// 关闭出错的连接
	if err == rpc.ErrShutdown ||
		err == io.ErrUnexpectedEOF ||
		err == io.EOF {

		conn.Close()
		return
	}

	// 放回池子
	p.Lock()
	conns := p.conns[addr]
	if len(conns) >= p.size {
		p.Unlock()
		conn.Close()
		return
	}
	p.conns[addr] = append(conns, conn)
	p.Unlock()
}
