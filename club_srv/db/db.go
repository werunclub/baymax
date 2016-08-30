package db

import (
	"baymax/common/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Db *gorm.DB

func Init(database_url string) {
	var err error
	Db, err = gorm.Open("mysql", database_url)
	Db.SingularTable(true)

	if err != nil {
		panic("fail to connect database")
	}

	Db.AutoMigrate(&model.Club{})
	//defer Db.Close()
}
