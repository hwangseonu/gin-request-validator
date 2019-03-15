package gin_validator

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

func JsonRequiredMiddleware(json interface{}) gin.HandlerFunc {
	t := reflect.TypeOf(json)
	return func(c *gin.Context) {
		var m map[string]interface{}
		if err := c.ShouldBindJSON(&m); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		tmp := reflect.New(t).Interface()
		if err := ValidData(m, tmp); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		c.Set("json", reflect.ValueOf(tmp).Elem().Interface())
		c.Next()
	}
}

func GetJsonData(c *gin.Context) interface{} {
	if json, ok := c.Get("json"); ok {
		return json
	} else {
		return nil
	}
}
