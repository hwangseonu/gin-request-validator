package gin_validator

import (
	"errors"
	"reflect"
	"strings"
)

var validators = map[string]Validator{
	"email":    EmailValidator,
	"notblank": NotBlankValidator,
	"min": MinValidator,
	"max": MaxValidator,
}

/*
	Validator 함수는 구조체 필드의 유효성을 검사하는 함수입니다.
	name 은 구조체 필드의 이름입니다.
	data 는 구조체 인스턴스 필드에 들어있는 데이터입니다.
	args 는 콤마(,)로 구분되는 validate 태그의 인자들입니다.
*/
type Validator func(name string, data interface{}, args ...string) error

func ValidData(json map[string]interface{}, must interface{}) error {
	if err := mapToStruct(json, must); err != nil {
		return err
	}
	data := reflect.ValueOf(must).Elem().Interface()
	structType := reflect.TypeOf(data)

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		validateTag := field.Tag.Get("validate")
		if validateTag == "" {
			continue
		}
		isRequired, data := checkRequiredField(json, must, field.Name)
		if isRequired && data == nil {
			return errors.New(field.Name + " is required field")
		} else if !isRequired && data == nil {
			continue
		}

		for _, c := range strings.Split(validateTag, " ") {
			args := strings.Split(c, "=")
			validator := validators[args[0]]
			if validator == nil {
				continue
			}
			if err := validator(field.Name, data, args[1:]...); err != nil {
				return err
			}
		}
	}
	return nil
}