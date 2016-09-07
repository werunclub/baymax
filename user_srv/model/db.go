package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/Sirupsen/logrus"
)

var DB *gorm.DB

func Init(database_url string, logMode bool) {
	log.WithField("schema", database_url).Debug("Connecting to database")
	var err error
	DB, err = gorm.Open("mysql", database_url)
	DB.SingularTable(true)
	DB.LogMode(logMode)

	if err != nil {
		panic("fail to connect database")
	}

	DB.AutoMigrate(User{})
	//defer Db.Close()
}

func Close() {
	DB.Close()
}

