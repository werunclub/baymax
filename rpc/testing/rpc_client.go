package main

import (
	"baymax/errors"
	"baymax/rpc"
	//"time"
	//"time"
	"log"
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

func main() {

	registry := rpc.NewConsulRegistry()
	registry.Init()

	services, _ := registry.GetService("RpcTest")

	for _, s := range services {
		log.Print(s.Address)
	}

	selector := rpc.NewSelector(rpc.PoolSize(10))

	client, err := selector.Select("RpcTest")

	if err != nil {
		log.Printf("error: %v", err.Error())
		return
	}

	args := &Args{7, 8}
	reply := new(Reply)

	err1 := client.Call("Arith.Add", args, reply)
	if err1 != nil {
		log.Printf("error: %v", err1.Error())
	}

	log.Printf("res: %v", reply)

	args = &Args{7, 8}
	reply = new(Reply)

	err1 = client.Call("Arith.Mul", args, reply)
	if err1 != nil {
		log.Printf("error: %v", err1.Error())
	}
	log.Printf("res: %v", reply)
}
