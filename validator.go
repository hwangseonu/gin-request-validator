//gin_validator 는 gin 의 request body 의 유효성을 구조체의 정의를 통해 검사하는 라이브러리입니다.
//Handler 에 불필요한 구조체의 유효성을 검사하는 코드가 포함되는 것을 막기 위해 만들었습니다.
package gin_validator

import (
	"errors"
	"reflect"
	"strings"
)

//기본 제공되는 유효성 검사 함수들입니다.
var validators = map[string]Validator{
	"email":    EmailValidator,
	"notblank": NotBlankValidator,
	"min":      MinValidator,
	"max":      MaxValidator,
	"pattern":  PatternValidator,
}

//Validator 함수는 구조체 필드의 유효성을 검사하는 함수입니다.
//name 은 구조체 필드의 이름입니다.
//data 는 구조체 인스턴스 필드에 들어있는 데이터입니다.
//args 는 콤마(,)로 구분되는 validate 태그의 인자들입니다.
type Validator func(name string, data interface{}, args ...string) error


//인자로 받은 map 을 구조체로 변환하여 유효성을 검사합니다.
//Required 필드가 없으면 error 를 반환합니다.
//Required 하지 않은 필드에 데이터가 없으면 유효성을 검사하지 않습니다.
//Required 하지 않은 필드에 데이터가 있으면 유효성을 검사합니다.
//유효성검사를 통과하지 못하면 error 를 반환합니다.
//유효성검사를 통과하면 nil 을 반환합니다.
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
