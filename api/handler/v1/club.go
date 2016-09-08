package v1

import (
	"github.com/gin-gonic/gin"

	clubProto "baymax/club_srv/protocol/club"
	"baymax/rpc"
	"strconv"
	"time"
	"log"
)

func getClubRpcConn() *rpc.Client {
	return rpc.NewClient("tcp", "127.0.0.1:8091", 10*time.Minute)
}

// 根据 id 获取 club 信息
func GetClub(c *gin.Context) {
	var err error

	clubIDStr := c.Params.ByName("clubId")
	clubID, err := strconv.Atoi(clubIDStr)

	if err != nil {
		c.JSON(400, err)
	} else {
		clubRpcConn := getClubRpcConn()

		var reply clubProto.GetOneReply
		err = clubRpcConn.Call(clubProto.SrvGetOneClub, &clubProto.GetOneArgs{ClubID: clubID}, &reply)

		if err != nil {
			c.JSON(400, err)
		} else {
			c.JSON(200, reply.Data)
		}
	}
}

// 创建 club
type CreateClubQuery struct {
	CityCode    string `json:"city_code" binding:"required"`
	Description string `json:"description"`
	IndustryID  int    `json:"industry_id"`
	Name        string `json:"name"`
}
func CreateClub(c *gin.Context) {
	var (
		q CreateClubQuery
		err error
	)

	err = c.BindJSON(&q)
	if err != nil {
		log.Println(err)
		c.JSON(400, err)
	} else {
		clubRpcConn := getClubRpcConn()

		var reply clubProto.CreateReply
		args := clubProto.CreateArgs{}

		args.Club.Name = q.Name
		args.Club.CityCode = q.CityCode
		args.Club.IndustryID = q.IndustryID
		args.Club.Des = q.Description

		err = clubRpcConn.Call(clubProto.SrvCreateClub, &args, &reply)
		if err != nil {
			c.JSON(400, err)
		} else {
			c.JSON(200, reply.Club)
		}
	}
}


func UpdateClub(c *gin.Context) {
	var (
		err error
		q UpdateClubQuery
	)

	clubIDStr := c.Params.ByName("clubId")
	clubID, err := strconv.Atoi(clubIDStr)

	if err != nil {
		log.Println(err)
		c.JSON(400, err)
	} else {
		err := c.BindJSON(&q)

		if err != nil {
			c.JSON(400, err)
		} else {
			clubRpcConn := getClubRpcConn()

			reply := clubProto.UpdateReply{}
			args := clubProto.UpdateArgs{}

			args.ClubID = clubID

			args.NewClub.CityCode = q.CityCode
			args.NewClub.Des = q.Description
			args.NewClub.IndustryID = q.IndustryID
			args.NewClub.Name = q.Name
			args.NewClub.VerifyCode = q.VerifyCode

			err := clubRpcConn.Call(clubProto.SrvUpdateClub, &args, &reply)
			if err != nil {
				log.Println(err)
				c.JSON(400, err)
			} else {
				c.JSON(200, reply.Club)
			}
		}
	}

}