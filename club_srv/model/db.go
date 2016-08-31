package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func Init(database_url string) {
	var err error
	DB, err = gorm.Open("mysql", database_url)
	if err != nil {
		panic("fail to connect database")
	}

	// 设置默认表名不再是复数
	DB.SingularTable(true)

	// DB.AutoMigrate(&model.Club{})
	// defer DB.Close()
}

func Close() {
	Db.Close()
}
