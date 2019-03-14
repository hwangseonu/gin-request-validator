package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/hwangseonu/gin-request-validator"
	"log"
	"net/http"
)

type Data struct {
	Email  string `json:"email"`
	Age    int    `json:"age" validate:"min=1 max=7"`
	Status string `json:"status" validate:"custom"`
}

func CustomValidator(name string, data interface{}, interfaces ...interface{}) error {
	str, ok := data.(string)
	if !ok {
		return errors.New(name + "must string")
	}
	if str != "happy" {
		return errors.New("you must be happy")
	}
	return nil
}

func main() {
	gin_validator.RegisterValidator("custom", CustomValidator, "string")
	r := gin.Default()
	r.Use(gin_validator.JsonRequiredMiddleware(Data{}))
	r.POST("/", func(c *gin.Context) {
		req := gin_validator.GetJsonData(c).(Data)
		c.JSON(http.StatusOK, req)
	})

	log.Fatal(r.Run(":5000"))
}
