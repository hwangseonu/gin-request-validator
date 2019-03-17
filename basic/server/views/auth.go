package views

import (
	"github.com/gin-gonic/gin"
	"github.com/hwangseonu/gin-restful"
	"github.com/hwangseonu/gin-restful-example/basic/server/models"
	"github.com/hwangseonu/gin-restful-example/basic/server/security"
	"net/http"
	"time"
)

type Auth struct {
	*gin_restful.Resource
}

type SignInRequest struct {
	Username string `json:"username" binding:"required,notblank"`
	Password string `json:"password" binding:"required,notblank"`
}

func InitAuthResource() Auth {
	r := Auth{gin_restful.InitResource()}
	return r
}

func (r Auth) Post(req SignInRequest) (gin.H, int) {
	u := models.FindUserByUsername(req.Username)
	if u == nil {
		return gin.H{"message": "cannot find user by username"}, http.StatusNotFound
	}
	if u.Password != req.Password {
		return gin.H{"message": "invalid password"}, http.StatusUnauthorized
	}
	access := security.GenerateToken(security.ACCESS, u.Username)
	refresh := security.GenerateToken(security.REFRESH, u.Username)
	return AuthResponse(access, refresh), http.StatusOK
}

const DAY = 24 * time.Hour

type Refresh struct {
	*gin_restful.Resource
}

func InitRefreshResource() Refresh {
	r := Refresh{gin_restful.InitResource()}
	r.AddMiddleware(security.AuthRequired(security.REFRESH), http.MethodGet)
	return r
}

func (r Refresh) Get(c *gin.Context) (gin.H, int) {
	u := c.MustGet("user").(*models.UserModel)
	exp := c.MustGet("exp").(int64)

	access := security.GenerateToken(security.ACCESS, u.Username)
	refresh := security.GenerateToken(security.REFRESH, u.Username)

	if time.Unix(exp, 0).Before(time.Now().Add(7 * DAY)) {
		return AuthResponse(access, refresh), http.StatusOK
	} else {
		return AuthResponse(access, ""), http.StatusOK
	}
}
