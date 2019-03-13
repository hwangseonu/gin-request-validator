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

type Validator struct {
	Func      func(...interface{}) error
	Arguments []string
}

var validators = map[string]Validator{
	"email": {Func: EmailValidator, Arguments: []string{"string"}},
	"min": {Func: MinValidator, Arguments: []string{"int", "int"}},
	"max": {Func: MaxValidator, Arguments: []string{"int", "int"}},
	"pattern": {Func: PatternValidator, Arguments: []string{"string", "string"}},
	"notblank": {Func: NotBlackValidator, Arguments: []string{"string"}},
}

func JsonRequiredMiddleware(json interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := reflect.TypeOf(json)
		if err := c.ShouldBindJSON(&json); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		str := mapToStruct(json.(map[string]interface{}), t)
		if err := ValidData(str); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.Set("json", str)
		c.Next()
	}
}

func GetJsonData(c *gin.Context) interface{} {
	json, ok := c.Get("json")
	if !ok {
		return nil
	}
	return json
}

func ValidData(json interface{}) error {
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
				return errors.New("data1 is must " + validator.Arguments[0])
			}
			if len(validator.Arguments) == 1 {
				if err := validator.Func(data.Interface()); err != nil {
					return err
				}
			} else {
				args := strings.Split(t[1], ",")
				as := []interface{}{data.Interface()}
				for i, k := range args {
					k := getTrueType(k, validator.Arguments[i+1])
					as = append(as, k)
				}
				if err := validator.Func(as...); err != nil {
					return err
				}
			}
		}
	}
	return nil
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