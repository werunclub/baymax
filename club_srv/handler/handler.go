package handler

import (
	"baymax/club_srv/model"
	proto "baymax/club_srv/protocol/club"
)

type ClubHandler struct{}

func (*ClubHandler) Get(req *proto.GetRequest, resp *proto.GetResponse) error {

	club := model.Club{}
	err := model.DB.Where("id = ?", req.ClubID).First(&club).Error

	if err != nil {
		return err
	} else {
		c, err := proto.InitFromModel(&club)
		if err != nil {
			return err
		} else {
			resp.Data = []proto.Club{c}
			return nil
		}
	}
}

// 创建新的 club 返回 id
func (*ClubHandler) Create(req *proto.CreateRequest, resp *proto.CreateResponse) error {

	club := model.Club{}

	// TODO 两种类型的 struct 之间的转换
	club.Name = req.Name

	var err error
	err = model.DB.Create(&club).Error

	if err != nil {
		return err
	} else {
		resp.ClubID = club.ID
		return nil
	}
}

// 删除指定的 club
func (*ClubHandler) Delete(req *proto.DeleteRequest, resp *proto.DeleteResponse) error {

	club := model.Club{ID: req.ClubID}

	var err error
	// 原数据库没有 delete_at 字段, 指定记录会被直接删除
	err = model.DB.Delete(&club).Error
	if err != nil {
		return err
	} else {
		resp.ClubID = club.ID
	}

	return nil
}

// 更新指定的 club
func (*ClubHandler) Update(req *proto.UpdateRequest, res *proto.UpdateResponse) error {
	club := model.Club{}

	var err error
	err = model.DB.Where("id = ?", req.Club.ID).First(&club).Error

	if err != nil {
		return err
	} else {
		club.Name = req.Name

		err = model.DB.Model(&club).Update().Error
		if err != nil {
			return err
		} else {
			(*res).Club, err = proto.InitFromModel(&club)
			if err != nil {
				return err
			} else {
				return nil
			}
		}
	}

}
