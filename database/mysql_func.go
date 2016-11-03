package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
)

var MysqlAdvancedFunc mysqlAdvancedFunc

func NewMysqlAdvancedFunc() mysqlAdvancedFunc {
	return mysqlAdvancedFunc{}
}

func init() {
	MysqlAdvancedFunc = NewMysqlAdvancedFunc()
}

type mysqlAdvancedFunc struct{}

// insert into ... on duplicate key update ...
func (mysqlAdvancedFunc) InsertOnDuplicateKeyUpdate(db *gorm.DB, tableName string, insertMap map[string]interface{}, updateMap map[string]interface{}) error {
	insertColumn := make([]string, 0, len(insertMap))
	insertPlaceholders := make([]string, 0, len(insertMap))
	updateExpr := make([]string, 0, len(updateMap))
	args := make([]interface{}, 0, len(insertMap)+len(updateMap))

	for column, value := range insertMap {
		insertColumn = append(insertColumn, column)
		insertPlaceholders = append(insertPlaceholders, "?")
		args = append(args, value)
	}

	for column, value := range updateMap {
		updateExpr = append(updateExpr, fmt.Sprintf("%v=?", column))
		args = append(args, value)
	}

	sqlExpr := fmt.Sprintf(
		"INSERT INTO %v (%v) VALUES (%v) ON DUPLICATE KEY UPDATE %v",
		tableName,
		strings.Join(insertColumn, ","),
		strings.Join(insertPlaceholders, ","),
		strings.Join(updateExpr, ","),
	)
	if err := db.Exec(sqlExpr, args...).Error; err != nil {
		return err
	}
	return nil
}

/*
MySQL uses the following algorithm for REPLACE (and LOAD DATA ... REPLACE):

1.Try to insert the new row into the table

2.While the insertion fails because a duplicate-key error occurs for a primary key or unique index:

	a.Delete from the table the conflicting row that has the duplicate key value

	b.Try again to insert the new row into the table
*/
func (mysqlAdvancedFunc) Replace(db *gorm.DB, tableName string, replaceMap map[string]interface{}) error {
	columns := make([]string, 0, len(replaceMap))
	replacePlaceholders := make([]string, 0, len(replaceMap))
	args := make([]interface{}, 0, len(replaceMap))
	for column, value := range replaceMap {
		columns = append(columns, column)
		replacePlaceholders = append(replacePlaceholders, "?")
		args = append(args, value)
	}

	sqlExpr := fmt.Sprintf(
		"REPLACE INTO %v (%v) VALUES (%v)",
		tableName,
		strings.Join(columns, ","),
		strings.Join(replacePlaceholders, ","),
	)
	if err := db.Exec(sqlExpr, args...).Error; err != nil {
		return err
	}
	return nil
}

/*
insert ignore into ...
the row won't actually be inserted if it results in a duplicate key
*/
func (mysqlAdvancedFunc) InsertIgnore(db *gorm.DB, tableName string, insertMap map[string]interface{}) error {
	columns := make([]string, 0, len(insertMap))
	insertPlaceholders := make([]string, 0, len(insertMap))
	args := make([]interface{}, 0, len(insertMap))
	for column, value := range insertMap {
		columns = append(columns, column)
		insertPlaceholders = append(insertPlaceholders, "?")
		args = append(args, value)
	}

	sqlExpr := fmt.Sprintf(
		"INSERT IGNORE INTO %v (%v) VALUES (%v)",
		tableName,
		strings.Join(columns, ","),
		strings.Join(insertPlaceholders, ","),
	)
	if err := db.Exec(sqlExpr, args...).Error; err != nil {
		return err
	}
	return nil
}

// check rows exists by given conditions
func (mysqlAdvancedFunc) CheckIfRowsExists(db *gorm.DB, tableName string, query string, args ...interface{}) (error, bool) {
	sqlExpr := fmt.Sprintf(
		"SELECT EXISTS(SELECT * FROM %v WHERE %v)",
		tableName,
		query,
	)
	var exist bool
	if err := db.Raw(sqlExpr, args...).Row().Scan(&exist); err != nil {
		return err, false
	}
	return nil, exist
}

func (mysqlAdvancedFunc) CountRows(db *gorm.DB, tableName string, query string, args ...interface{}) (error, int) {
	sqlExpr := fmt.Sprintf(
		"SELECT COUNT(*) FROM %v WHERE %v",
		tableName,
		query,
	)
	var num int
	if err := db.Raw(sqlExpr, args...).Row().Scan(&num); err != nil {
		return err, 0
	}
	return nil, num
}
