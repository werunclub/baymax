package main

import (
	club_proto "club-backend/club_srv/protocol/club"
	"club-backend/common/util"
	"fmt"

	logging "third/go-logging"
)

var log = logging.MustGetLogger("client")

func main() {

	var (
		err    error
		client *util.RpcClient
	)

	client, err = util.NewRpcClient("127.0.0.1:8081", "tcp", club_proto.ClubRpcFuncMap, "club_profile", log)
	if nil != err {
		fmt.Println(err)
	}

	var reply club_proto.GetResponse

	client.Call("get", &club_proto.GetRequest{ClubId: 1}, &reply)

	fmt.Println(reply)
}
