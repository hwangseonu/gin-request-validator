package gin_validator

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var validators = map[string]Validator{
	"email": {Func: EmailValidator, Arguments: []string{"string"}},
	"min": {Func: MinValidator, Arguments: []string{"int", "int"}},
	"max": {Func: MaxValidator, Arguments: []string{"int", "int"}},
	"pattern": {Func: PatternValidator, Arguments: []string{"string", "string"}},
	"notblank": {Func: NotBlankValidator, Arguments: []string{"string"}},
}
var emailRegex = regexp.MustCompile(`^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`)

/*
	Validator 구조체는 validator 함수와 인자의 종류를 함께 묶어놓은 구조체입니다.
	Custom Validator 를 등록하면 Validator 구조체의 인스턴스를 생성하여 등록합니다.
 */
type Validator struct {
	Func      func(name string, data interface{}, interfaces ...interface{}) error
	Arguments []string
}

/*
	ValidData 는 직접적으로 유효성을 검사하는 함수입니다.
	json 으로 받은 데이터를 must 와 비교하여 유효성을 검사합니다.
	문제가 없으면 nil 을, 문제가 있으면 error 를 반환합니다.
 */
func ValidData(json interface{}, must reflect.Type) error {
	json = mapToStruct(json, must)
	v := reflect.ValueOf(json)
	t := reflect.TypeOf(json)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fTag := f.Tag.Get("validate")
		if fTag == "" {
			continue
		}
		tags := strings.Split(fTag, " ")
		for _, tag := range tags {
			t := strings.Split(tag, "=")
			data := v.Field(i)
			validator := validators[t[0]]
			if data.Kind().String() != validator.Arguments[0] {
				return errors.New(f.Name + " must " + validator.Arguments[0])
			}
			if len(validator.Arguments) == 1 {
				if err := validator.Func(f.Name, data.Interface()); err != nil {
					return err
				}
			} else {
				args := strings.Split(t[1], ",")
				as := make([]interface{}, 0)
				for i, k := range args {
					k := getTrueType(k, validator.Arguments[i+1])
					as = append(as, k)
				}
				if err := validator.Func(f.Name, data.Interface(), as...); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

/*
	CustomValidator 를 등록합니다.
	name 은 validator 의 이름을, f 는 구현된 validator 함수를, args 는 validator 함수에 필요한 인자들의 타입을 순서대로 적습니다.
	args 의 제일 처음은 구조체 Field 의 타입입니다.
	나머지 args 는 `validate: "in=A,B,C"` 일때 A, B, C 같은 인자를 의미합니다.
 */
func RegisterValidator(name string, f func(name string, data interface{}, interfaces ...interface{}) error, dataType string, args ...string) {
	v := Validator{
		Func: f,
		Arguments: append([]string{dataType}, args...),
	}
	validators[name] = v
}

/*
	구조체의 Field 가 email 임을 확인합니다.
	아니면 error 를 반환합니다.
	구조체 Field Tag 에 `validate: "email"`처럼 사용합니다.
 */
func EmailValidator(name string, data interface{}, i ...interface{}) error {
	str, err := mustString(data)
	if err != nil {
		return err
	}
	if !emailRegex.MatchString(str) {
		return errors.New(str + " is not email")
	}
	return nil
}

/*
	구조체의 Field 가 최소 min 임을 확인합니다.
	아니면 error 를 반환합니다.
	구조체 Field Tag 에 `validate: "min=5"`처럼 사용합니다.
 */
func MinValidator(name string, data interface{}, i ...interface{}) error {
	d, err := mustInt(data)
	if err != nil {
		return err
	}
	min, err := mustInt(i[0])
	if err != nil {
		return err
	}
	if min > d {
		return errors.New(name + " must greater than " + strconv.Itoa(min) + ", but actual value is " + strconv.Itoa(d))
	}
	return nil
}

/*
	구조체의 Field 가 최대 max 임을 확인합니다.
	아니면 error 를 반환합니다.
	구조체 Field Tag 에 `validate: "max=5"`처럼 사용합니다.
 */
func MaxValidator(name string, data interface{}, i ...interface{}) error {
	d, err := mustInt(data)
	if err != nil {
		return err
	}
	max, err := mustInt(i[0])
	if err != nil {
		return err
	}
	if max < d {
		return errors.New(name + " must less than " + strconv.Itoa(max) + ", but actual value is " + strconv.Itoa(d))
	}
	return nil
}

/*
	구조체의 Field 가 정규식 pattern 과 일치함을 확인합니다.
	아니면 error 를 반환합니다.
	구조체 Field Tag 에 `validate: "pattern=^[A-Z]$"`처럼 사용합니다.
 */
func PatternValidator(name string, data interface{}, i ...interface{}) error {
	var str, regexStr string
	var err error
	var r *regexp.Regexp

	if str, err = mustString(data); err != nil {
		return err
	}
	if regexStr, err = mustString(i[0]); err != nil {
		return err
	}
	if r, err = regexp.Compile(regexStr); err != nil {
		return nil
	}
	if r.MatchString(str) {
		return nil
	} else {
		return errors.New(name + " is not matched pattern")
	}
}

/*
	구조체의 Field 가 문자열이고 비어있지 않음을 확인합니다.
	아니면 error 를 반환합니다.
	구조체 Field Tag 에 `validate: "notblank"`처럼 사용합니다.
 */
func NotBlankValidator(name string, data interface{}, i ...interface{}) error {
	if str, err := mustString(data); err != nil {
		return err
	} else if str == "" {
		return errors.New(name + " should not blank")
	} else {
		return nil
	}
}

func mustString(i interface{}) (string, error) {
	if s, ok := i.(string); ok {
		return s, nil
	} else {
		return "", errors.New("interface is not string")
	}
}

func mustFloat(i interface{}) (float64, error) {
	if s, ok := i.(float64); ok {
		return s, nil
	} else {
		return 0, errors.New("interface is not float")
	}
}

func mustInt(i interface{}) (int, error) {
	if s, ok := i.(int); ok {
		return s, nil
	} else {
		return 0, errors.New("interface is not integer")
	}
}

func mapToStruct(v interface{}, t reflect.Type) interface{} {
	m, ok := v.(map[string]interface{})
	if !ok {
		return v
	}
	s := reflect.New(t).Elem()

	for k, v := range m {
		if s.Kind() == reflect.Struct {
			k = string(unicode.ToUpper(rune(k[0]))) + k[1:]
			f := s.FieldByName(k)
			if f.IsValid() {
				if f.CanSet() {
					switch f.Type().Kind() {
					case reflect.Int:
						x := int64(v.(float64))
						if !f.OverflowInt(x) {
							f.SetInt(x)
						}
						break
					case reflect.Float64:
						x := v.(float64)
						if !f.OverflowFloat(x) {
							f.SetFloat(x)
						}
						break
					case reflect.String:
						x := v.(string)
						f.SetString(x)
						break
					}
				}
			}
		}
	}
	return s.Interface()
}