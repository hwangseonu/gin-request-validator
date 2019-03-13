package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hwangseonu/gin-request-validator"
	"log"
	"net/http"
)

type Data struct {
	Email string `json:"email" validate:"pattern=^[A-Z]$"`
	Age   int    `json:"age" validate:"min=1 max=7"`
}


func main() {
	r := gin.Default()

	r.Use(gin_validator.JsonRequiredMiddleware(Data{}))

	r.POST("/", func(c *gin.Context) {
		req := gin_validator.GetJsonData(c).(Data)
		c.JSON(http.StatusOK, req)
	})

	log.Fatal(r.Run(":5000"))
}
