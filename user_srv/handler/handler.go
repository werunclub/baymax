package handler

import (
	"baymax/rpc"
	"baymax/user_srv/protocol/user"
	"reflect"
	"github.com/jinzhu/gorm"
	"fmt"
	log "github.com/Sirupsen/logrus"
)

func RegisterRPCService(server *rpc.Server) {
	server.RegisterName(user.RegistryName, userHandler{})
}

// 将协议转换成 ORM 可以使用的 WhereCondition
func protocolToConditions(protocol interface{}) *map[string]interface{} {
	var (
		rTypes reflect.Type
		rValues reflect.Value
	)
	conditions := make(map[string]interface{})

	if reflect.TypeOf(protocol).Kind() == reflect.Ptr {
		rTypes = reflect.TypeOf(protocol).Elem()
		rValues = reflect.ValueOf(protocol).Elem()
	} else {
		rTypes = reflect.TypeOf(protocol)
		rValues = reflect.ValueOf(protocol)
	}

	for i := 0; i < rTypes.NumField(); i++ {
		field := rTypes.Field(i)
		value := rValues.FieldByName(field.Name)
		if value.Interface() != reflect.Zero(field.Type).Interface() {
			conditions[field.Tag.Get("condition")] = value.Interface()
		}
	}
	log.WithField("Conditions", conditions).Debug("解析后的查询条件")
	return &conditions
}


// 根据 protocol 参数执行 db.where
func applyWhereFilter(db *gorm.DB, protocol interface{}) *gorm.DB {
	conditions := protocolToConditions(protocol)
	query := ""
	args := []interface{}{}
	for key, value := range *conditions {
		//db = db.Where(key, value)
		if query == "" {
			query = key
		} else {
			query += fmt.Sprintf(" AND %s ", key)
		}
		args = append(args, value)
	}
	return db.Where(query, args...)
}