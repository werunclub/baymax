package main

import (
	"fmt"
	"log"
	"time"

	"github.com/werunclub/baymax/errors"
	rpcClient "github.com/werunclub/baymax/rpc/client"
	"github.com/werunclub/baymax/rpc/helpers"

	"github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/share"
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

	logrus.SetLevel(logrus.DebugLevel)

	start := time.Now()

	client := rpcClient.NewClient("Arith",
		rpcClient.EtcdAddress([]string{"127.0.0.1:2379"}),
		rpcClient.ConnTimeout(time.Second*5),
	)

	for i := 0; i <= 10; i++ {

		args := &Args{7, 8}
		reply := new(Reply)

		err1 := client.Call("Add", args, reply)
		if err1 != nil {
			log.Printf("error: %v", err1.Error())
		}

		log.Printf("res: %v", reply)

		{
			args = &Args{7, 8}
			reply = new(Reply)

			ctx := helpers.NewMetaDataContext(map[string]string{
				"lang": "cn",
			})

			err1 := client.CallWithContext(ctx, "Arith.Mul", args, reply)
			if err1 != nil {
				log.Printf("error: %v", err1.Error())
			}
			log.Printf("res: %v", reply)
			log.Printf("received meta: %+v", ctx.Value(share.ResMetaDataKey))
		}
	}

	{
		client := rpcClient.NewClient("Arith2",
			rpcClient.EtcdAddress([]string{"127.0.0.1:2379"}),
			rpcClient.ConnTimeout(time.Second*5),
		)

		args := &Args{7, 8}
		reply := new(Reply)

		err1 := client.Call("Add2", args, reply)
		if err1 != nil {
			log.Printf("error: %v", err1.Error())
		}

		log.Printf("Arith2.Add2: %v", reply)
	}

	{
		client := rpcClient.NewClient("Arith3",
			rpcClient.EtcdAddress([]string{"127.0.0.1:2379"}),
			rpcClient.ConnTimeout(time.Second*5),
		)

		args := &Args{7, 8}
		reply := new(Reply)

		err1 := client.Call("Add3", args, reply)
		if err1 != nil {
			log.Printf("error: %v", err1.Error())
		}

		log.Printf("Arith3.Add3: %v", reply)
	}

	{
		client := rpcClient.NewClient("Arith4",
			rpcClient.EtcdAddress([]string{"127.0.0.1:2379"}),
			rpcClient.ConnTimeout(time.Second*5),
		)

		args := &Args{7, 8}
		reply := new(Reply)

		err1 := client.Call("Add4", args, reply)
		if err1 != nil {
			log.Printf("error: %v", err1.Error())
		}

		log.Printf("Arith4.Add4: %v", reply)
	}

	{
		client := rpcClient.NewClient("Arith5",
			rpcClient.EtcdAddress([]string{"127.0.0.1:2379"}),
			rpcClient.ConnTimeout(time.Second*5),
		)

		args := &Args{7, 8}
		reply := new(Reply)

		err1 := client.Call("Add5", args, reply)
		if err1 != nil {
			log.Printf("error: %v", err1.Error())
		}

		log.Printf("Arith5.Add5: %v", reply)
	}

	{
		client := rpcClient.NewClient("Arith6",
			rpcClient.EtcdAddress([]string{"127.0.0.1:2379"}),
			rpcClient.ConnTimeout(time.Second*5),
		)

		args := &Args{7, 8}
		reply := new(Reply)

		err1 := client.Call("Add6", args, reply)
		if err1 != nil {
			log.Printf("error: %v", err1.Error())
		}

		log.Printf("Arith6.Add6: %v", reply)
	}

	{
		client := rpcClient.NewClient("Arith7",
			rpcClient.EtcdAddress([]string{"127.0.0.1:2379"}),
			rpcClient.ConnTimeout(time.Second*5),
		)

		args := &Args{7, 8}
		reply := new(Reply)

		err1 := client.Call("Add7", args, reply)
		if err1 != nil {
			log.Printf("error: %v", err1.Error())
		}

		log.Printf("Arith7.Add7: %v", reply)
	}

	{
		client := rpcClient.NewClient("Arith8",
			rpcClient.EtcdAddress([]string{"127.0.0.1:2379"}),
			rpcClient.ConnTimeout(time.Second*5),
		)

		args := &Args{7, 8}
		reply := new(Reply)

		err1 := client.Call("Add8", args, reply)
		if err1 != nil {
			log.Printf("error: %v", err1.Error())
		}

		log.Printf("Arith8.Add8: %v", reply)
	}

	fmt.Printf("du: %d", time.Now().Nanosecond()-start.Nanosecond())
}
