package gin_validator

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

//JsonRequiredMiddleware 를 url 의 미들웨어로 등록하여 유효성검사를 할 수 있습니다.
//인자로 검사할 구조체의 비어있는 인스턴스를 받습니다.
//요청을 인자로 받은 구조체로 변황하여 유효성을 검사합니다.
//유효한 요청일 경우 context 에 정보를 저장합니다.
//만약 요청이 유효하지 않으면 상태코드 400을 반환하며 요청을 종료합니다.
func JsonRequiredMiddleware(json interface{}) gin.HandlerFunc {
	mustType := reflect.TypeOf(json)
	return func(c *gin.Context) {
		m := make(map[string]interface{})
		if err := c.ShouldBindJSON(&m); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		must := reflect.New(mustType).Interface()
		if err := ValidData(m, must); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.Set("json", reflect.ValueOf(must).Elem().Interface())
		c.Next()
	}
}

//반드시 JsonRequiredMiddleware 가 적용된 Handler 에서 사용해야합니다.
//유효성 검사로 context 에 저장된 데이터를 구조체로 반환합니다.
//함수 호출로 얻은 결과를 원하는 구조체로 변환하여 사용합니다.
func GetJsonData(c *gin.Context) interface{} {
	if json, ok := c.Get("json"); ok {
		return json
	} else {
		return nil
	}
}
