package gin_validator

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

func JsonRequiredMiddleware(json interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := reflect.TypeOf(json)
		must := reflect.New(t).Interface()
		if err := c.ShouldBindJSON(&json); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if err := ValidData(json, must); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		must = reflect.ValueOf(must).Elem().Interface()
		c.Set("json", must)
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