package club

import (
	"baymax/club_srv/model"
	"time"
)

const (
	prefix = "Club."

	// 获取指定的俱乐部信息
	ServiceGetClub = prefix + "Get"
	// 创建新的俱乐部
	ServiceCreateClub = prefix + "Create"
	// 删除俱乐部
	ServiceDeleteClub = prefix + "Delete"
	// 获取多个俱乐部
	ServiceGetManyClub = prefix + "Update"
	// 修改俱乐部信息
	ServiceUpdateClub = prefix + "Update"
)

// TODO 代码复制到这里 注释也不能共享了
// TODO 参数以及返回值的名字不能命名为 req 以及 res　和　rpc　本身的名字不一样
type Club struct {
	ID         uint      `json:"id"`
	UserID     string    `json:"user_id"`
	Name       string    `json:"name"`
	Icon       string    `json:"icon"`
	Des        string    `json:"des"`
	ShortUrl   string    `json:"short_url"`
	SortNum    int       `json:"sort_num"`
	State      bool      `json:"state"`
	Authorized bool      `json:"authorized"`
	DataBody   string    `json:"data_body"`
	Source     int       `json:"source"`
	CityCode   string    `json:"city_code"`
	IndustryID int       `json:"industry_id"`
	CommonByte int       `json:"common_byte"`
	CreateTime time.Time `json:"create_time"`
}

// 根据 model 实例化一个 Club 结构的实例
func InitFromModel(m *model.Club) (Club, error) {
	c := Club{}

	c.ID = m.ID
	c.UserID = m.UserID
	c.Name = m.Name
	c.Icon = m.Icon
	c.Des = m.Des
	c.ShortUrl = m.ShortUrl
	c.SortNum = m.SortNum
	c.State = m.State
	c.Authorized = m.Authorized
	c.DataBody = m.DataBody
	c.Source = m.Source
	c.CityCode = m.CityCode
	c.IndustryID = m.IndustryID
	c.CommonByte = m.CommonByte
	c.CreateTime = m.CreateTime

	return c, nil
}

// 获取单条信息
type GetRequest struct {
	ClubID uint `json:"club_id"`
}

type GetResponse struct {
	Data []Club `json:"data"`
}

// 创建新的俱乐部
type CreateRequest struct {
	Club
}

type CreateResponse struct {
	ClubID uint `json:"club_id"`
}

// 删除俱乐部
type DeleteRequest struct {
	ClubID uint `json:"club_id"`
}

type DeleteResponse struct {
	ClubID uint `json:"club_id"`
}

// 更新俱乐部信息
type UpdateRequest struct {
	Club
}

type UpdateResponse struct {
	Club
}
