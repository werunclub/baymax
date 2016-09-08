package v1

import (
	"baymax/errors"
	"baymax/rpc"
	userProtocol "baymax/user_srv/protocol/user"
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
)

type UserHandler struct {
}

func NewUserHandler() UserHandler {
	return UserHandler{}
}

func handleRPCError(ctx *gin.Context, err *errors.Error, desc string) {
	// TODO: RPC client 使用 baymax/errors.Error 过后下面代码可以优化
	if err != nil {
		logrus.WithField(
			"error", logrus.Fields{
				"Id":     err.Id,
				"Status": err.Status,
				"Code":   err.Code,
				"Detail": err.Detail,
			}).Error(desc)
		ctx.JSON(http.StatusInternalServerError, err)
	}
}

func (handler *UserHandler) UserProfile(c *gin.Context) {
	// Fixme: RPC Client
	client := rpc.NewClient("tcp", "127.0.0.1:5000", 5*time.Second)

	userId := c.Param("userId")
	args := userProtocol.UserDetailRequest{Id: userId}
	reply := userProtocol.UserDetailReply{}
	err := client.Call(userProtocol.UserDetail, &args, &reply)
	if err != nil {
		handleRPCError(c, err, "RPC 获取用户详情失败")
		return
	}
	c.JSON(200, reply)
}

func (handler *UserHandler) PatchUser(ctx *gin.Context) {
	// Fixme: RPC Client
	client := rpc.NewClient("tcp", "127.0.0.1:5000", 5*time.Second)

	// 获取当前登录用户 userID
	user_id, exist := ctx.Get("userID")
	if !exist {
		logrus.Debug("从 gin.Context 中获取 userID 失败")
		ctx.JSON(http.StatusUnauthorized, errors.Unauthorized(""))
		return
	}
	// fucking type assert
	userId, ok := user_id.(string)
	if !ok {
		logrus.Error("gin.Context 获取到的 userId 不是 string 类型")
		ctx.JSON(http.StatusUnauthorized, errors.Unauthorized(""))
		return
	}

	// 从 RPC Server 获取当前登录用户的详情
	args := userProtocol.UserDetailRequest{Id: userId}
	user := userProtocol.UserDetailReply{}
	if err := client.Call(userProtocol.UserDetail, args, &user); err != nil {
		ctx.JSON(int(err.Status), err)
		return
	}

	// 解析 HTTP Request Body 解析为 Struct
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	form := userProtocol.UserInfo{}
	if err := json.Unmarshal(body, &form); err != nil {
		logrus.WithField("error", err).Debug("解析 JSON 失败")
	}

	// 使用 govalidator 校验参数是否正确
	valid, err := govalidator.ValidateStruct(form)
	logrus.WithFields(logrus.Fields{"valid": valid, "error": err}).Debug("校验结果")

	updated := userProtocol.UserDetailReply{}
	// fixme: userId 必须是当前登录用户的 ID, 因为 RPC Server 无法判断 targetUser 和 currentUser 是否同一个
	form.Id = userId
	arg := userProtocol.PatchUserRequest{Id: userId, Payload: &form}
	if err := client.Call(userProtocol.UserPatch, &arg, &updated); err != nil {
		handleRPCError(ctx, err, "RPC 修改用户失败")
		return
	}
	ctx.JSON(http.StatusOK, updated)
}
