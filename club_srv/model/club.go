package model

import (
	proto "baymax/club_srv/protocol/club"
	"fmt"
	"time"
)

// 俱乐部 model
type Club struct {
	ID          int   `gorm:"primary_key"`
	UserID      string `gorm:"size:36"`
	Name        string `gorm:"size:32;unique"`
	Icon        string `gorm:"size:128;default:''"`
	Des         string `gorm:"type:text"`
	ShortUrl    string
	PersonCount int   `gorm:"default:0"`
	SortNum     int    `gorm:"default:0"`
	State       bool   `gorm:"default:true"`
	Authorized  bool   `gorm:"default:true"`
	// 保留数据集合
	DataBody string `gorm:"size:1024;default:''"`
	Source   int    `gorm:"default:0"`
	// 俱乐部城市
	CityCode string `gorm:"size:8;default:0"`
	// 所属行业
	IndustryID int
	// -1位表示 是否参与排序 -2位表示 是否已经验证身份
	CommonByte int `gorm:"default:0"`
	CreateTime time.Time
}

// orm 用来定义数据库表名
func (c *Club) TableName() string {
	return "club_club"
}

func (c *Club) String() string {
	return fmt.Sprintf("club: %s", c.ID)
}

func FromProtoStruct(protoClub proto.Club) (Club, error) {
	m := Club{}

	m.ID = protoClub.ID
	m.UserID = protoClub.UserID
	m.Name = protoClub.Name
	m.Icon = protoClub.Icon
	m.Des = protoClub.Des
	m.ShortUrl = protoClub.ShortUrl
	m.SortNum = protoClub.SortNum
	m.State = protoClub.State
	m.Authorized = protoClub.Authorized
	m.DataBody = protoClub.DataBody
	m.Source = protoClub.Source
	m.CityCode = protoClub.CityCode
	m.IndustryID = protoClub.IndustryID
	m.CommonByte = protoClub.CommonByte
	m.CreateTime = protoClub.CreateTime

	return m, nil
}

// 根据 model 实例化一个 proto.Club 结构的实例
func ToProtoStruct(m Club) (proto.Club, error) {
	c := proto.Club{}

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

// 根据 model 数组实例化多个 proto.Club 结构的实例
func ToBatchProtoStruct(models []Club) ([]proto.Club, error) {
	var protoClubs []proto.Club

	for _, model := range models {
		protoClub, err := ToProtoStruct(model)
		if err != nil {
			return []proto.Club{}, nil
		} else {
			protoClubs = append(protoClubs, protoClub)
		}
	}

	return protoClubs, nil
}
