package util

import (
	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(false)
}

type Validator struct {
}

func (v *Validator) ValidateStruct(post interface{}) error {
	_, err := govalidator.ValidateStruct(post)
	return err
}
