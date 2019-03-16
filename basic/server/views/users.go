package views

import (
	"github.com/gin-gonic/gin"
	"github.com/hwangseonu/gin-restful"
	"github.com/hwangseonu/gin-restful-example/basic/server/models"
	"net/http"
)

type Users struct {
	*gin_restful.Resource
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,notblank"`
	Password string `json:"password" binding:"required,notblank"`
}

func InitUsersResource() Users {
	r := Users{Resource: gin_restful.InitResource()}
	return r
}

func (r Users) Post(req CreateUserRequest) (gin.H, int) {
	if models.ExistsUserByUsername(req.Username) {
		return gin.H{}, http.StatusConflict
	}
	u := models.NewUserModel(req.Username, req.Password)
	u.Save()
	return gin.H{
		"username": u.Username,
		"password": u.Password,
	}, http.StatusCreated
}


