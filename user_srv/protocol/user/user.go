package user

import (
	"time"
)

const (
	// used for rpc.Server.RegisterName
	RegistryName = "User"
	// 获取用户详情
	UserDetail = RegistryName + ".Get"
	// 获取用户列表
	UserList = RegistryName + ".List"
	// 创建用户
	UserCreate = RegistryName + ".Create"
	// 修改用户
	UserPatch = RegistryName + ".Patch"
)

type UserId uint

type UserIds []uint

type Photo struct {
	Size [3]string `json:"sizes"`
	Url  string    `json:"url"`
}

type UserInfo struct {
	Id         string     `json:"id"`
	LastLogin  *time.Time `json:"last_login"`
	Name       string     `json:"name"`
	UpdateTime *time.Time `json:"update_time"`
	// 是否安装企业咕咚 APP
	AppInstalled bool       `json:"app_installed"`
	Avatar       *Photo     `json:"avatar"`
	CreateTime   *time.Time `json:"create_time"`
	// 开启免打扰模式
	DNDEnabled bool `json:"dnd_enabled"`
	// 免打扰开始时间
	DNDEnd uint `json:"dnd_end"`
	// 免打扰结束时间
	DNDStart uint `json:"dnd_start"`
	// 生日
	Dob    *time.Time `json:"dob"`
	Email  string     `json:"email"`
	Gender string     `json:"gender"`
	// 目标步数
	Goal   uint    `json:"goal"`
	Height uint8   `json:"height"`
	Weight float32 `json:"weight"`
}

type PatchUser struct {
	// 开启免打扰模式
	DNDEnabled bool `json:"dnd_enabled"`
	// 免打扰开始时间
	DNDEnd uint `json:"dnd_end"`
	// 免打扰结束时间
	DNDStart uint `json:"dnd_start"`
	// 生日
	Dob time.Time `json:"dob"`
	//Email  string    `json:"email"`
	Gender string `json:"gender"`
	// 目标步数
	Goal   uint    `json:"goal"`
	Height uint8   `json:"height"`
	Weight float32 `json:"weight"`
}

type PatchUserRequest struct {
	Id      string
	Payload *UserInfo
}

type UserDetailRequest struct {
	Id   string    `condition:"id = ?"`
	Name string    `condition:"name = ?"`
	Dob  time.Time `condition:"dob >= ?"`
}

type UserDetailReply UserInfo

type ResponseReply struct {
	Success string
}

type UserListRequest struct {
	Ids []uint
}
type UserListReply struct{ Users *[]UserInfo }
