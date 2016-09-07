package model

import (
	"time"
	userProtocol "baymax/user_srv/protocol/user"
	"fmt"
)

type User struct {
	Id           string `gorm:"primary_key"`
	Dob          *time.Time
	AvatarUrl    string

	LastLogin    *time.Time `json:"last_login"`
	Name         string    `json:"name"`
	UpdateTime   *time.Time `json:"update_time"`
	// 是否安装企业咕咚 APP
	AppInstalled bool      `json:"app_installed" valid:"optional"`
	CreateTime   *time.Time `json:"create_time"`
	// 开启免打扰模式
	DNDEnabled   bool `json:"dnd_enabled"`
	// 免打扰开始时间
	DNDEnd       uint `json:"dnd_end"`
	// 免打扰结束时间
	DNDStart     uint `json:"dnd_start"`
	// 生日
	Email        string    `json:"email"`
	Gender       string    `json:"gender"`
	// 目标步数
	Goal         uint    `json:"goal"`
	Height       uint8   `json:"height"`
	Weight       float32 `json:"weight"`
}

// ReadonlyFields 在对用户模型进行更新时, 某些字段是不允许更新的, 或这需要从额外的入口进行更新
// 例如 `last_login` 只能在用户登录时进行修改不能在, 不能在用户修改资料时修改
func (u *User) ReadonlyFields() []string {
	return []string{"id", "last_login", "update_time", "create_time"}
}

func (u *User) Avatar() *userProtocol.Photo {
	return &userProtocol.Photo{
		Size: [3]string{"@small.jpg", "@medium.jpg", "@large.jpg"},
		Url: u.AvatarUrl,
	}
}

//func (User) TableName() string {
//	return "user"
//}

//FromProtocol 将 RPC Protocol 转换为 .
func (user *User) FromProtocol(protocol *userProtocol.UserInfo) (*User, error) {
	if protocol == nil {
		return nil, fmt.Errorf("`*userProtocol.UserInfo` 为空指针")
	}
	user.Id = protocol.Id
	user.Dob = protocol.Dob
	if protocol.Avatar != nil {
		user.AvatarUrl = protocol.Avatar.Url
	}
	user.LastLogin = protocol.LastLogin
	user.Name = protocol.Name
	user.UpdateTime = protocol.UpdateTime
	// 是否安装企业咕咚 APP
	user.AppInstalled = protocol.AppInstalled
	user.CreateTime = protocol.CreateTime
	// 开启免打扰模式
	user.DNDEnabled = protocol.DNDEnabled
	// 免打扰开始时间
	user.DNDEnd = protocol.DNDEnd
	// 免打扰结束时间
	user.DNDStart = protocol.DNDStart
	// 生日
	user.Email = protocol.Email
	user.Gender = protocol.Gender
	// 目标步数
	user.Goal = protocol.Goal
	user.Height = protocol.Height
	user.Weight = protocol.Weight
	return user, nil
}

// ToUserInfo 将模型转换成 protocol.user.UserInfo .
func (user *User) ToUserInfo() *userProtocol.UserInfo {
	info := userProtocol.UserInfo{
		Id: user.Id,
		Dob: user.Dob,
		Avatar: user.Avatar(),
		LastLogin: user.LastLogin,
		Name: user.Name,
		UpdateTime: user.UpdateTime,
		// 是否安装企业咕咚 APP
		AppInstalled: user.AppInstalled,
		CreateTime: user.CreateTime,
		// 开启免打扰模式
		DNDEnabled: user.DNDEnabled,
		// 免打扰开始时间
		DNDEnd: user.DNDEnd,
		// 免打扰结束时间
		DNDStart: user.DNDStart,
		// 生日
		Email: user.Email,
		Gender: user.Gender,
		// 目标步数
		Goal: user.Goal,
		Height: user.Height,
		Weight: user.Weight,
	}
	return &info
}

