package client

import (
	"net"
	"net/rpc"

	"github.com/Sirupsen/logrus"
)

type LogRpc struct {
}

func (l LogRpc) Name() string {
	return "log connect plugin"
}

func (l LogRpc) HandleConnected(conn net.Conn) (net.Conn, bool) {
	logrus.WithField("address", conn.RemoteAddr()).Debugf("connect to server")
	return conn, true
}

func (l LogRpc) PreReadResponseHeader(resp *rpc.Response) error {
	logrus.WithField("resp", resp).Debugf("pre read resp header")
	return nil
}
func (l LogRpc) PostReadResponseHeader(resp *rpc.Response) error {
	logrus.WithField("resp", resp).Debugf("read resp header")
	return nil
}
func (l LogRpc) PreReadResponseBody(args interface{}) error {
	logrus.WithField("args", args).Debugf("pre read Response body")
	return nil
}
func (l LogRpc) PostReadResponseBody(args interface{}) error {
	logrus.WithField("args", args).Debugf("read Response body")
	return nil
}

func (l LogRpc) PreWriteRequest(req *rpc.Request, args interface{}) error {
	logrus.WithField("req", req).Debugf("pre write Request")
	return nil
}
func (l LogRpc) PostWriteRequest(req *rpc.Request, args interface{}) error {
	logrus.WithField("req", req).Debugf("write Request")
	return nil
}
