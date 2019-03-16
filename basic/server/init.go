package server

import (
	"github.com/gin-gonic/gin"
	"github.com/hwangseonu/gin-restful-example/basic/server/views"
)

func CreateServer() *gin.Engine {
	e := gin.Default()
	gin.SetMode(gin.DebugMode)
	views.RegisterViews(e)
	return e
}
