package rpc

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"
	"time"
)

var (
	DefaultPoolSize = 100
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
func (p *pool) DialTimeout(network, addr string, timeout time.Duration) (*rpc.Client, error) {

	var (
		conn net.Conn
		err  error
	)

	conn, err = net.DialTimeout(network, addr, timeout)
	if err != nil {
		return nil, err
	}

	return jsonrpc.NewClient(conn), nil
}

func (p *pool) GetConn(addr string, timeout time.Duration) (*poolConn, error) {
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

	// 新连接
	c, err := p.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return nil, err
	}
	return &poolConn{c, time.Now().Unix()}, nil
}

func (p *pool) release(addr string, conn *poolConn, err error) {

	// 关闭出错的连接
	if err != nil {
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
