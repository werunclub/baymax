package handler_test


import (
	"testing"

	clubProto "baymax/club_srv/protocol/club"
	"baymax/rpc"
	"time"
)

func getConn() *rpc.Client {
	return rpc.NewClient("tcp", "127.0.0.1:8091", 10 * time.Minute)
}

// TODO 启动一个 server

// 获取单条记录
func TestGetOneClub(t *testing.T) {
	c := getConn()
	var reply clubProto.GetResponse

	err := c.Call(clubProto.ServiceGetClub, &clubProto.GetRequest{ClubID: 45}, &reply)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(reply)
	}
}

// 根据 ID 获取多条记录
func TestGetManyClubs(t *testing.T) {
	c := getConn()
	var reply clubProto.GetManyResponse

	err := c.Call(clubProto.ServiceGetManyClub, &clubProto.GetManyRequest{}, &reply)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(reply)
	}
}