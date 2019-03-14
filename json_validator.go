package gin_validator

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"strconv"
)

/*
	JsonRequiredMiddleware 는 구조체의 인스턴스를 인자로 받아 요청에서 json 데이터를 추출하여 매칭하고 유효성을 검사하여
	문제가 있으면 400 을 반환하고 문제가 없으면 json 데이터를 context 에 저장합니다.
 */
func JsonRequiredMiddleware(json interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := reflect.TypeOf(json)
		json = mapToStruct(json, t)
		if err := c.ShouldBindJSON(&json); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		json = mapToStruct(json, t)
		if err := ValidData(json, t); err != nil {
			msg := make([]string, 0)
			for _, e := range err {
				msg = append(msg, e.Error())
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		json = mapToStruct(json, t)
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