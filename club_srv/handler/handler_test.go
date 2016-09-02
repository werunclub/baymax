package handler_test


import (
	"testing"

	clubProto "baymax/club_srv/protocol/club"
	"baymax/rpc"
	"time"
	"log"
)

const SERVER_ADDR = "127.0.0.1:8091"


// 测试获取单条记录
func TestGetOneClub(t *testing.T) {

	c := rpc.NewClient("tcp", SERVER_ADDR, 10 * time.Minute)

	var reply clubProto.GetResponse

	err := c.Call(clubProto.ServiceGetClub, &clubProto.GetRequest{ClubID: 45}, &reply)
	if err != nil {
		t.Error(err)
	} else {
		log.Println(reply)
	}
}