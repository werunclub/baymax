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
	var reply clubProto.GetOneReply

	err := c.Call(clubProto.SrvGetOneClub, &clubProto.GetOneArgs{ClubID: 45}, &reply)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(reply)
	}
}

// 根据 ID 获取多条记录
func TestGetManyClubs(t *testing.T) {
	c := getConn()
	var reply clubProto.GetBatchReply

	err := c.Call(clubProto.SrvGetBatchClub, &clubProto.GetBatchArgs{ClubIDS: []int{24, 38, 59}}, &reply)

	if err != nil {
		t.Error(err)
	} else {
		t.Log(reply)
	}
}

// 根据名字查询俱乐部信息
func TestSearchClubs(t *testing.T) {
	c := getConn()
	var reply clubProto.SearchReply

	err := c.Call(clubProto.SrvSearchClub, &clubProto.SearchArgs{Name:"咕咚"}, &reply)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(reply)
	}
}

// 创建新的 club
func TestCreateClub(t *testing.T) {
	c := getConn()
	var reply clubProto.CreateReply

	args := clubProto.CreateArgs{}
	args.Club.Name = "咕咚来了哇"
	args.Club.CreateTime = time.Now()

	err := c.Call(clubProto.SrvCreateClub, &args, &reply)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(reply)
	}
}

// 修改 club 信息
func TestUpdateClub(t *testing.T) {
	c := getConn()
	var reply clubProto.UpdateReply

	args := clubProto.UpdateArgs{}
	args.ClubID = 319
	args.NewClub.Name = "咕咚又来啦哇"
	args.NewClub.CreateTime = time.Now()

	err := c.Call(clubProto.SrvUpdateClub, &args, &reply)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(reply)
	}
}
