package rpc

import (
	"testing"
	"time"

	"baymax/errors"
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

func (t *Arith) Error(args *Args, reply *Reply) error {
	panic("ERROR")
}

func startServer() {
	server = NewServer()
	server.RegisterName("Arith", new(Arith))
	server.Start()
}

func TestDirectClient(t *testing.T) {
	once.Do(startServer)

	client := NewDirectClient("tcp", server.Address(), time.Second*5)

	args := &Args{7, 8}
	reply := new(Reply)

	err := client.Call("Arith.Add", args, reply)
	if err != nil {
		t.Errorf("Add: expected no error but got string %q", err.Error())
	}
	if reply.C != args.A+args.B {
		t.Errorf("Add: got %d expected %d", reply.C, args.A+args.B)
	}

	args = &Args{7, 8}
	reply = new(Reply)
	err = client.Call("Arith.Mul", args, reply)
	if err != nil {
		t.Errorf("Mul: expected no error but got string %q", err.Error())
	}
	if reply.C != args.A*args.B {
		t.Errorf("Mul: got %d expected %d", reply.C, args.A*args.B)
	}

	args = &Args{7, 0}
	reply = new(Reply)
	err = client.Call("Arith.Div", args, reply)
	if err == nil {
		t.Error("Div: expected error but got nil")
	}
}

func TestDirectClientNoneServer(t *testing.T) {

	client := NewDirectClient("tcp", "127.0.0.1:11212", time.Second*5)

	args := &Args{7, 8}
	reply := new(Reply)

	err := client.Call("Arith.Add", args, reply)
	if err == nil {
		t.Error("Add: expected an error but got nil")
	}
}
