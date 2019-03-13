package gin_validator

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

var validators = map[string]Validator{
	"email": {Func: EmailValidator, Arguments: []string{"string"}},
	"min": {Func: MinValidator, Arguments: []string{"int", "int"}},
	"max": {Func: MaxValidator, Arguments: []string{"int", "int"}},
	"pattern": {Func: PatternValidator, Arguments: []string{"string", "string"}},
	"notblank": {Func: NotBlackValidator, Arguments: []string{"string"}},
}

/*
	Validator 구조체는 validator 함수와 인자의 종류를 함께 묶어놓은 구조체입니다.
	Custom Validator 를 등록하면 Validator 구조체의 인스턴스를 생성하여 등록합니다.
 */
type Validator struct {
	Func      func(name string, interfaces ...interface{}) error
	Arguments []string
}

/*
	JsonRequiredMiddleware 는 구조체의 인스턴스를 인자로 받아 요청에서 json 데이터를 추출하여 매칭하고 유효성을 검사하여
	문제가 있으면 400 을 반환하고 문제가 없으면 json 데이터를 context 에 저장합니다.
 */
func JsonRequiredMiddleware(json interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := reflect.TypeOf(json)
		if err := c.ShouldBindJSON(&json); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		if err := ValidData(json, t); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		m, ok := json.(map[string]interface{})
		if ok {
			json = mapToStruct(m, t)
		}
		c.Set("json", json)
		c.Next()
	}
}

/*
	GetJsonData 는 JsonRequiredMiddleware 가 context 에 저장한 json 데이터를 추출하여 반환합니다.
	데이터가 없다면 nil 을 반환합니다.
 */
func GetJsonData(c *gin.Context) interface{} {
	json, ok := c.Get("json")
	if !ok {
		return nil
	}
	return json
}

/*
	ValidData 는 직접적으로 유효성을 검사하는 함수입니다.
	json 으로 받은 데이터를 must 와 비교하여 유효성을 검사합니다.
	문제가 없으면 nil 을, 문제가 있으면 error 를 반환합니다.
 */
func ValidData(json interface{}, must reflect.Type) error {
	m, ok := json.(map[string]interface{})
	if ok {
		json = mapToStruct(m, must)
	}
	v := reflect.ValueOf(json)
	t := reflect.TypeOf(json)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := strings.Split(f.Tag.Get("validate"), " ")
		for _, raw := range tag {
			t := strings.Split(raw, "=")
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
				as := []interface{}{data.Interface()}
				for i, k := range args {
					k := getTrueType(k, validator.Arguments[i+1])
					as = append(as, k)
				}
				if err := validator.Func(f.Name, as...); err != nil {
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
func RegistValidator(name string, f func(name string, interfaces ...interface{}) error, args ...string) {
	v := Validator{
		Func: EmailValidator,
		Arguments: args,
	}
	validators[name] = v
}

func getTrueType(i interface{}, must string) interface{} {
	str := i.(string)
	switch must {
	case "int":
		if i, err := strconv.Atoi(str); err != nil {
			return 0
		} else {
			return i
		}
	case "float":
		if f, err := strconv.ParseFloat(str, 64); err != nil {
			return 0.0
		} else {
			return f
		}
	case "bool":
		if str == "" || str == "false" || str == "0" || str == "null" || str == "nil" || str == "off" {
			return true
		} else {
			return false
		}
	case "string":
		return str
	default:
		return i
	}
}

func mapToStruct(m map[string]interface{}, t reflect.Type) interface{} {
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