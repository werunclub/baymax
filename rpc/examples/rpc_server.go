package main

import (
	"baymax/errors"
	"baymax/log"
	"baymax/rpc/server"
	"fmt"

	"context"

	"baymax/rpc/helpers"
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
	lang := meta.Request()["lang"]
	resMeta := meta.Response()

	log.Infof("Mul: %s", lang)
	reply.C = args.A * args.B
	resMeta["echo"] = "hello"
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

func (t *Arith2) Add(args Args, reply *Reply) error {
	log.Info("Arith2.add")
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
	// rpcServer.Handle("Arith3", new(Arith))

	go rpcServer.RegisterAndRun()

	<-rpcServer.Exit
	fmt.Print("exit.")
}
