package model

import (
	"time"
)

// 俱乐部
type Club struct {
	ID uint `gorm:"primary_key"`
	UserID string `gorm:"size:36"`
	Name string `gorm:"size:32;unique"`
	Icon string `gorm:"size:128;default:''"`
	Des string `gorm:"type:text"`
	ShortUrl string `gorm:"type:"`
	PersonCount uint `gorm:"default:0"`
	SortNum int `gorm:"default:0"`
	State bool `gorm:"default:true"`
	Authorized bool `gorm:"default:true"`
	// 保留数据集合
	DataBody string `gorm:"size:1024;default:''"`
	Source int `gorm:"default:0"`
	// 俱乐部城市
	CityCode string `gorm:"size:8;default:0"`
	// 所属行业
	IndustryID int
	// -1位表示 是否参与排序 -2位表示 是否已经验证身份
	CommonByte int `gorm:"default:0"`
	CreateTime time.Time
}

func (c *Club) TableName() string {
	return "club_club"
}