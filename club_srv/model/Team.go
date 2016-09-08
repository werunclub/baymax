package model

import (
	"fmt"
	"time"
)

// 部门 model
type Team struct {
	ID     uint `gorm:"primary_key"`
	ClubID int  `gorm:"index"`
	Parent int  `gorm:"index"`

	Name     string `gorm:"size:32"`
	DataBody string `gorm:"size:512;default:''"`
	Count    int    `gorm:"default:0"`

	CreateTime time.Time
}

func (*Team) TableName() {
	return "club_team"
}

func (t *Team) String() string {
	return fmt.Sprintf("team: %s", t.ID)
}
