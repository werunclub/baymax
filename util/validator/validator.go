package validator

import (
	"fmt"
	"reflect"
	"strings"
)

/* ========== JSON Form Validator ============ */
// FieldError 参数校验时, 特定字段的错误信息
type FieldError struct {
	// 对应 struct 的属性
	Field string

	// struct 的 json 输出名
	JSONField string

	// 错误信息
	Errors []string
}

func (e *FieldError) AppendError(err string) {
	e.Errors = append(e.Errors, err)
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("invalid agrument `%s`, with errors `%s`", e.JSONField, strings.Join(e.Errors, ";"))
}

// Validate 校验表单参数; 不支持嵌套的 struct, array, slice
// form 要校验的表单 struct, 必须传指针
// in 客户端提交的 fields
// out 输出校验通过的字段和值
func ValidateJSONStruct(form interface{}, in map[string]interface{}, out *map[string]interface{}) error {
	formTyp := reflect.TypeOf(form)
	formValue := reflect.ValueOf(form)
	if formTyp.Kind() == reflect.Ptr {
		formTyp = formTyp.Elem()
	}

	for i := 0; i < formTyp.NumField(); i++ {
		field := formTyp.Field(i)

		// 获取对应的 json 字段名
		tagString, ok := field.Tag.Lookup("json")
		if !ok {
			tagString = field.Name
		}

		tagValues := strings.Split(tagString, ",")
		jsonName := tagValues[0]

		// 如果有提交则进行参数校验
		value, exist := in[jsonName]
		if exist {
			e := &FieldError{Field: field.Name, JSONField: jsonName}

			//// 1. 先校验数据的原始类型是否相同; 如果类型错误, 立即终止当前字段的校验
			//lTyp := field.Type
			//if lTyp.Kind() == reflect.Ptr {
			//	lTyp = field.Type.Elem()
			//}
			//rTyp := reflect.TypeOf(value)
			//if rTyp.Kind() == reflect.Ptr {
			//	rTyp = reflect.TypeOf(value).Elem()
			//}
			//
			//if lTyp.Kind() != rTyp.Kind() {
			//	return fmt.Errorf("类型错误 %s != %s", lTyp.Kind(), rTyp.Kind())
			//	//return fmt.Errorf("类型错误")
			//}

			// 2. 获取 Field 对应的校验方法
			validators := findValidators(formValue, field)

			// 3. 执行校验方法
			for _, validator := range *validators {
				returnValues := validator.Call([]reflect.Value{reflect.ValueOf(value)})
				if len(returnValues) == 1 {
					// 1 个返回值时只能返回 error;
					if rtn := returnValues[0].Interface(); rtn != nil {
						e.AppendError(rtn.(error).Error())
					}
				} else {
					if rtn := returnValues[1].Interface(); rtn != nil {
						e.AppendError(rtn.(error).Error())
					} else {
						value = returnValues[0].Interface()
					}
				}
			}

			// 4. 校验后如果有错误信息, 返回 error, 校验失败
			if len(e.Errors) != 0 {
				return e
			}
			//formValue.Elem().Field(i).Set(reflect.ValueOf(value))
			(*out)[field.Name] = value
		} else if isRequired(field) {
			e := &FieldError{Field: field.Name, JSONField: jsonName}
			e.AppendError("该字段为必填项")
			return e
		}
	}

	if method := formValue.MethodByName("Validate"); method != (reflect.Value{}) {
		rtnValues := method.Call([]reflect.Value{})
		if rtn := rtnValues[0].Interface(); rtn != nil {
			return rtn.(error)
		}
	}
	return nil
}

func findValidators(formTyp reflect.Value, field reflect.StructField) *[]reflect.Value {
	validators := []reflect.Value{}

	// 查看是否有自定义的校验方法, 如果有就执行; 校验方法必须是绑定该 form 的实例方法; 并且校验方法必须是可导出的
	if validatorTag, ok := field.Tag.Lookup("validator"); ok {
		// 同一个字段可以有多个自定义校验方法
		_validators := strings.Split(validatorTag, ",")
		for _, validator := range _validators {
			if method := formTyp.MethodByName(validator); method != (reflect.Value{}) {
				validators = append(validators, method)
			}
		}
	}

	// 同时会调用 Validate_{FieldName}
	if method := formTyp.MethodByName("Validate_" + field.Name); method != (reflect.Value{}) {
		validators = append(validators, method)
	}

	return &validators
}

func isRequired(field reflect.StructField) bool {
	if tagValidators, ok := field.Tag.Lookup("validator"); ok {
		_validators := strings.Split(tagValidators, ",")
		for _, _validator := range _validators {
			if _validator == "required" {
				return true
			}
		}
	}
	return false
}

/* ========== JSON Form Validator End ============ */
