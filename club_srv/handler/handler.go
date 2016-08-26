package handler

import (
	proto "club-backend/club_srv/protocol/club"
	"log"
)

type ClubHandler int

func (handler *ClubHandler) Get(req *proto.GetRequest, resp *proto.GetResponse) error {
	log.Print("Get Club info")

	resp.TotalNum = 100
	return nil
}
