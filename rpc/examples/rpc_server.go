package main

import (
	"baymax/errors"
	"baymax/rpc"

	"github.com/prometheus/common/log"
	"github.com/smallnest/rpcx"
)

type Args struct {
	A, B int
}

type Reply struct {
	C int
}

type Arith int

type ArithAddResp struct {
	Id     interface{} `json:"id"`
	Result Reply       `json:"result"`
	Error  interface{} `json:"error"`
}

func (t *Arith) Add(args *Args, reply *Reply) error {
	log.Info("add")
	reply.C = args.A + args.B
	return nil
}

func (t *Arith) Mul(args *Args, reply *Reply) error {
	reply.C = args.A * args.B
	return nil
}

func (t *Arith) Div(args *Args, reply *Reply) error {
	if args.B == 0 {
		return errors.BadRequest("divide by zero")
	}
	reply.C = args.A / args.B
	return nil
}

type Arith2 int

func (t *Arith2) Add(args *Args, reply *Reply) error {
	log.Info("add")
	reply.C = args.A + args.B
	return nil
}

type Arith3 int

func (t *Arith3) Add(args *Args, reply *Reply) error {
	log.Info("add")
	reply.C = args.A + args.B
	return nil
}

func main() {
	server := rpc.NewServer(
		rpc.ConsulAddress("127.0.0.1:8500"))

	server.Handle("Arith", new(Arith))
	server.Handle("Arith2", new(Arith))
	server.Handle("Arith3", new(Arith))
	server.RegisterAndRun()
}
