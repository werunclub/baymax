package handler

import (
	"baymax/errors"
	"baymax/user_srv/model"
	userProtocol "baymax/user_srv/protocol/user"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/pborman/uuid"
)

type userHandler struct{}

func (userHandler) Create(args *userProtocol.UserInfo, reply *userProtocol.UserDetailReply) error {
	user, err := (&model.User{}).FromProtocol(args)
	if err != nil {
		log.WithField("error", err).Error("协议转模型出错")
		return errors.InternalServerError("系统内部错误")
	}
	user.Id = uuid.New()
	err = model.DB.Create(&user).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err.Error(), "user": user}).Error("创建用户失败")
		return errors.InternalServerError(err.Error())
	}
	*reply = userProtocol.UserDetailReply(*user.ToUserInfo())
	return nil
}

// Get  根据用户 ID 获取用户详情
func (userHandler) Get(args *userProtocol.UserDetailRequest, reply *userProtocol.UserDetailReply) error {
	user := model.User{Id: args.Id}
	//db := applyWhereFilter(model.DB, args)
	if err := model.DB.Find(&user).Error; err != nil {
		log.WithField("error", err).Debug("获取用户出错")
		return errors.NotFound("用户不存在")
	}

	reply.Id = user.Id
	reply.Avatar = user.Avatar()
	reply.Name = user.Name
	reply.Dob = user.Dob

	return nil
}

// 批量获取用户详情
func (userHandler) List(userIds *[]string, reply *userProtocol.UserListReply) error {
	users := []model.User{}
	err := model.DB.Where("id in (?)", *userIds).Find(&users).Error
	if err != nil {
		fmt.Println(err)
		return err
	}

	results := []userProtocol.UserInfo{}
	for _, user := range users {
		log.WithField("user", user).Debug("获取到用户")
		info := userProtocol.UserInfo{
			Id:   user.Id,
			Name: user.Name,
		}
		results = append(results, info)
	}
	reply.Users = &results

	return nil
}

func (userHandler) Patch(args *userProtocol.PatchUserRequest, reply *userProtocol.UserDetailReply) error {
	update, _ := (&model.User{}).FromProtocol(args.Payload)
	userId := args.Id
	if userId == "" {
		log.WithField("userId", userId).Error("修改用户 userId 为空!")
		return errors.BadRequest("必须提供要修改用户的 userId")
	}
	log.WithField("userId", userId).Debugf("修改用户 User[%s]\n", userId)

	user := model.User{}
	model.DB.Find(&user, "id = ?", userId)

	model.DB.Model(&user).Omit(user.ReadonlyFields()...).Update(update)

	*reply = userProtocol.UserDetailReply(*user.ToUserInfo())

	return nil
}
