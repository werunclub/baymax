package handler

import (
	"baymax/club_srv/model"
	protoClub "baymax/club_srv/protocol/club"
	"log"
)


type ClubHandler struct{}

func (*ClubHandler) GetOne(args *protoClub.GetOneArgs, reply *protoClub.GetOneReply) error {
	log.Println("call GetOne")

	var club model.Club
	err := model.DB.Where("id = ?", args.ClubID).First(&club).Error

	if err != nil {
		return err
	} else {
		c, err := model.ToProtoStruct(club)
		if err != nil {
			return err
		} else {
			reply.Data = c
			return nil
		}
	}
}

func (*ClubHandler) GetBatch(args *protoClub.GetBatchArgs, reply *protoClub.GetBatchReply) error {
	log.Println("call GetBatch")

	var clubs []model.Club
	err := model.DB.Where("id in (?)", args.ClubIDS).Find(&clubs).Error

	if err != nil {
		return err
	} else {
		protoClubs, err := model.ToBatchProtoStruct(clubs)
		if err != nil {
			return err
		} else {
			reply.Total = len(protoClubs)
			reply.Data = protoClubs
			return nil
		}
	}

	return nil
}

func (*ClubHandler) Search(args *protoClub.SearchArgs, reply *protoClub.SearchReply) error {
	log.Println("call Search")

	var clubs []model.Club

	likeArg := "%" + args.Name + "%"
	err := model.DB.Where("name LIKE ?", likeArg).Find(&clubs).Offset(args.Offset).Limit(args.Limit).Error

	if err != nil {
		return nil
	} else {
		protoClubs, err := model.ToBatchProtoStruct(clubs)
		if err != nil {
			return err
		} else {
			reply.Total = len(protoClubs)
			reply.Data = protoClubs
			return nil
		}
	}
}

func (*ClubHandler) Create(args *protoClub.CreateArgs, reply *protoClub.CreateReply) error {
	log.Println("call Create")
	club, err := model.FromProtoStruct(args.Club)

	if err != nil {
		return err
	} else {
		err := model.DB.Create(&club).Error
		if err != nil {
			return nil
		} else {
			protoClub, err := model.ToProtoStruct(club)
			if err != nil {
				return err
			} else {
				reply.Club = protoClub
				return nil
			}
		}
	}
}

func (*ClubHandler) Update(args *protoClub.UpdateArgs, reply *protoClub.UpdateReply) error {
	log.Println("call Update")

	var club model.Club
	err := model.DB.Where("id = ?", args.ClubID).First(&club).Error

	if err != nil {
		return err
	} else {
		newClub, err := model.FromProtoStruct(args.NewClub)
		if err != nil {
			return err
		} else {
			err := model.DB.Model(&club).Update(newClub).Error
			if err != nil {
				return err
			} else {
				c, err := model.ToProtoStruct(club)
				if err != nil {
					return nil
				} else {
					reply.Club = c
					return nil
				}
			}
		}
	}

	return nil
}
