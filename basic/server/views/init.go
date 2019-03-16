package views

import (
	"github.com/gin-gonic/gin"
	"github.com/hwangseonu/gin-restful"
)

func RegisterViews(e *gin.Engine) {
	users := InitUsersResource()
	api := gin_restful.NewApi(e, "/")
	api.AddResource(users, "/users")
}
