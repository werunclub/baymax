package handler

import (
	"baymax/club_srv/model"
	proto "baymax/club_srv/protocol/club"
	"fmt"
)

type ClubHandler struct {}

func (hdl *ClubHandler) Get(req *proto.GetRequest, resp *proto.GetResponse) error {

	club := model.Club{}
	err := model.DB.Where("id = ?", req.ClubId).First(&club).Error

	if err != nil {
		return err
	} else {
		resp.Data = []proto.Club{}
		return nil
	}
}

func (hdl *ClubHandler) Create(req *proto.CreateRequest, resp *proto.CreateResponse) error {
	return nil
}
