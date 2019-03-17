package views

import (
	"github.com/gin-gonic/gin"
	"github.com/hwangseonu/gin-restful"
	"github.com/hwangseonu/gin-restful-example/basic/server/models"
	"github.com/hwangseonu/gin-restful-example/basic/server/security"
	"net/http"
)

type Users struct {
	*gin_restful.Resource
}

type SignUpRequest struct {
	Username string `json:"username" binding:"required,notblank"`
	Password string `json:"password" binding:"required,notblank"`
	Nickname string `json:"nickname" binding:"required,notblank"`
	Email    string `json:"email" binding:"required,notblank,email"`
}

func InitUsersResource() Users {
	r := Users{Resource: gin_restful.InitResource()}
	r.AddMiddleware(security.AuthRequired(security.ACCESS, "ROLE_USER"), http.MethodGet)
	return r
}

func (r Users) Get(c *gin.Context) (gin.H, int) {
	u := c.MustGet("user").(*models.UserModel)
	return UserResponse(u), http.StatusOK
}

func (r Users) Post(req SignUpRequest) (gin.H, int) {
	if models.ExistsUserByUsernameOrNicknameOrEmail(req.Username, req.Nickname, req.Email) {
		return gin.H{}, http.StatusConflict
	}
	u := models.NewUserModel(req.Username, req.Password, req.Nickname, req.Email, "ROLE_USER")
	u.Save()
	return UserResponse(u), http.StatusCreated
}
