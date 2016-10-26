package client

import (
	"testing"

	"baymax/rpc/registry"
	"baymax/rpc/server"
	"time"
)

func TestSelector(t *testing.T) {

	server := server.NewServer(
		server.ConsulAddress("127.0.0.1:8500"),
	)

	server.Handle("Arith2", new(Arith))
	server.Registry.Init()
	server.Start()
	server.Register()

	defer func() {
		server.Deregister()
	}()

	selector := registry.NewSelector(registry.ConsulAddress("127.0.0.1:8500"))

	next, err := selector.Select("Arith2")
	if err != nil {
		t.Errorf("Select: expected no error but got string %q", err.Error())
	}

	node, err := next()
	if err != nil {
		t.Errorf("Select: expected no error but got string %q", err.Error())
	}

	client := NewDirectClient("tcp", node.Address, time.Second*5)

	args := &Args{7, 8}
	reply := new(Reply)

	err1 := client.Call("Arith.Add", args, reply)
	if err1 != nil {
		t.Errorf("Add: expected no error but got string %q", err1.Error())
	}
	if reply.C != args.A+args.B {
		t.Errorf("Add: got %d expected %d", reply.C, args.A+args.B)
	}

	args = &Args{7, 8}
	reply = new(Reply)

	err1 = client.Call("Arith.Mul", args, reply)
	if err1 != nil {
		t.Errorf("Mul: expected no error but got string %q", err1.Error())
	}
	if reply.C != args.A*args.B {
		t.Errorf("Mul: got %d expected %d", reply.C, args.A*args.B)
	}

	args = &Args{7, 0}
	reply = new(Reply)

	err = client.Call("Arith.Div", args, reply)
	if err == nil {
		t.Errorf("Div: expected error but got nil")
	}
}
