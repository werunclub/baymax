package db

import (
	"github.com/jinzhu/gorm"
	"club-backend/common/model"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)
var Db *gorm.DB

func Init(database_type string, database_dsn string)  {
	var err error
	Db, err = gorm.Open(database_type, database_dsn)
	Db.SingularTable(true)

	if err != nil {
		panic("fail to connect database")
	}

	Db.AutoMigrate(&model.Club{})
	//defer Db.Close()
}