package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hwangseonu/gin-request-validator"
	"net/http"
)

//유효성을 검사할 구조체를 정의합니다.
//binding 태그를 이용하여 꼭 필요한 필드를 나타낼 수 있습니다.
type Data struct {
	Email string `json:"email" validate:"email" binding:"required"`
	Age   int    `json:"age" validate:"min=1 max=100"`
}

//요청으로 생성된 데이터를 상태코드 200과 한께 그대로 반환하는 Handler 입니다.
func Handler(c *gin.Context) {
	req := gin_validator.GetJsonData(c).(Data)
	c.JSON(http.StatusOK, req)
}

func main() {
	e := gin.Default()
	g := e.Group("/awesome")

	// "/awesome" url 에 JsonRequiredMiddleware 를 미들웨어로 등록합니다.
	//인자값으로 정의된 구조체의 빈 인스턴스를 넘겨줍니다
	g.Use(gin_validator.JsonRequiredMiddleware(Data{}))

	g.POST("", Handler)
	_ = e.Run(":5000")
}
