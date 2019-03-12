package gin_validator

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type Validator struct {
	Func      func(...interface{}) error
	Arguments []string
}

var validators = map[string]Validator{
	"email": {Func: EmailValidator, Arguments: []string{"string"}},
	"min": {Func: MinValidator, Arguments: []string{"int", "int"}},
	"max": {Func: MaxValidator, Arguments: []string{"int", "int"}},
}

func JsonRequiredMiddleware(json interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(&json); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		}
		c.Next()
	}
}

func GetJsonData(c *gin.Context, json interface{}) interface{} {
	if err := c.ShouldBindJSON(&json); err != nil {
		log.Fatal(err)
	}
	return json
}

func ValidData(json interface{}) error {
	v := reflect.ValueOf(json)
	t := reflect.TypeOf(json)

	for i := 0; i < v.NumField(); i++ {
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
