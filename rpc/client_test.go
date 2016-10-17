package rpc

import (
	"testing"
	"time"

	"sync"
)

var (
	server *Server
	once   sync.Once

	serviceName   = "Arith"
	consulAddress = "127.0.0.1:8500"
)

func startAndRegistServer() {
	server = NewServer(
		ConsulAddress(consulAddress),
	)
	server.RegisterName(serviceName, new(Arith))
	go server.RegisterAndRun()
}

func TestClient(t *testing.T) {
	once.Do(startAndRegistServer)

	client := NewClient(serviceName, consulAddress, time.Second*5)

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

func TestNoneServer(t *testing.T) {

	client := NewDirectClient("tcp", "127.0.0.1:11212", time.Second*5)

	args := &Args{7, 8}
	reply := new(Reply)

	err := client.Call("Arith.Add", args, reply)
	if err == nil {
		t.Error("Add: expected an error but got nil")
	}
}
