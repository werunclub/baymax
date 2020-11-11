package main

import (
	"context"
	"fmt"

	"baymax/errors"
	"baymax/log"
	"baymax/rpc/helpers"
	"baymax/rpc/server"
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

func (t *Arith) Add(ctx context.Context, args *Args, reply *Reply) error {
	log.Info("add")
	reply.C = args.A + args.B
	return nil
}

func (t *Arith) Mul(ctx context.Context, args *Args, reply *Reply) error {
	meta := helpers.NewMetaDataFormContext(ctx)
	lang := meta.Get("lang")
	resMeta := meta.Response()

	log.Infof("Mul: %s", lang)
	reply.C = args.A * args.B
	resMeta["echo"] = "hello"
	meta.Set("hello", "world")
	return nil
}

func (t *Arith) Div(ctx context.Context, args *Args, reply *Reply) error {
	if args.B == 0 {
		return errors.BadRequest("divide by zero")
	}
	reply.C = args.A / args.B
	return nil
}

type Arith2 int

func (t *Arith2) Add2(args Args, reply *Reply) error {
	log.Info("Arith2.add")
	reply.C = args.A + args.B
	return nil
}

type Arith3 int

func (t *Arith3) Add3(args Args, reply *Reply) error {
	log.Info("Arith3.add")
	reply.C = args.A + args.B
	return nil
}

type Arith4 int

func (t *Arith4) Add4(args Args, reply *Reply) error {
	log.Info("Arith4.add")
	reply.C = args.A + args.B
	return nil
}

type Arith5 int

func (t *Arith5) Add5(args Args, reply *Reply) error {
	log.Info("Arith5.add")
	reply.C = args.A + args.B
	return nil
}

type Arith6 int

func (t *Arith6) Add6(args Args, reply *Reply) error {
	log.Info("Arith6.add")
	reply.C = args.A + args.B
	return nil
}

type Arith7 int

func (t *Arith7) Add7(args Args, reply *Reply) error {
	log.Info("Arith7.add")
	reply.C = args.A + args.B
	return nil
}

type Arith8 int

func (t *Arith8) Add8(args Args, reply *Reply) error {
	log.Info("Arith8.add")
	reply.C = args.A + args.B
	return nil
}

func main() {
	rpcServer := server.NewServer(
		server.EtcdAddress([]string{"127.0.0.1:2379"}),
		server.StopWait(1),
	)

	rpcServer.Handle("Arith", new(Arith))
	rpcServer.Handle("Arith2", new(Arith2))
	rpcServer.Handle("Arith3", new(Arith3))
	rpcServer.Handle("Arith4", new(Arith4))
	rpcServer.Handle("Arith5", new(Arith5))
	rpcServer.Handle("Arith6", new(Arith6))
	rpcServer.Handle("Arith7", new(Arith7))
	rpcServer.Handle("Arith8", new(Arith8))

	go rpcServer.RegisterAndRun()

	<-rpcServer.Exit
	fmt.Print("exit.")
}
