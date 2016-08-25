package main

import (
    _ "github.com/jinzhu/gorm/dialects/mysql"

    "./model"
    "github.com/jinzhu/gorm"
)

const mysqlDSN = "vagrant:vagrant@(vagrant:3306)" +
    "/goorm?charset=utf8&parseTime=True&loc=Local"

func crateAll() error {
    return nil
}

func dropAll() error {
    return nil
}

func main() {
    db, err := gorm.Open("mysql", mysqlDSN)

    if err != nil {
        panic("failed to connect database")
    } else {
        db.CreateTable(&model.Club{})
        db.CreateTable(&model.Team{})
        db.CreateTable(&model.ClubPerson{})
        defer db.Close()
    }
}