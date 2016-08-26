package handler

import (
	proto "club-backend/club_srv/protocol/club"
	"log"
	db "club-backend/club_srv/db"
	"club-backend/common/model"
)

type ClubHandler int

func (handler *ClubHandler) Get(req *proto.GetRequest, resp *proto.GetResponse) error {
	log.Print("Get Club info")

	var club model.Club

	db.Db.First(&club)

	log.Print(club)

	resp.TotalNum = 100
	return nil
}
