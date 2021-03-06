package views

import (
	"github.com/gin-gonic/gin"
	"github.com/hwangseonu/gin-restful"
	"github.com/hwangseonu/gin-restful-example/basic/server/models"
	"github.com/hwangseonu/gin-restful-example/basic/server/security"
	"net/http"
	"time"
)

type Posts struct {
	*gin_restful.Resource
}

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,notblank"`
	Content string `json:"content" binding:"required,notblank"`
}

func InitPostsResource() Posts {
	r := Posts{gin_restful.InitResource()}
	r.AddMiddleware(security.AuthRequired(security.ACCESS), http.MethodPost, http.MethodPatch, http.MethodDelete)
	return r
}

func (r Posts) Get(id int) (gin.H, int) {
	p := models.FindPostById(id)
	if p == nil {
		return gin.H{"message": "cannot find post by id"}, http.StatusNotFound
	}
	return PostResponse(p), http.StatusOK
}

func (r Posts) Post(c *gin.Context, req CreatePostRequest) (gin.H, int) {
	user := c.MustGet("user").(*models.UserModel)
	p := models.NewPost(req.Title, req.Content, user)
	p.Save()
	return PostResponse(p), http.StatusCreated
}

func (r Posts) Patch(c *gin.Context, id int, req CreatePostRequest) (gin.H, int) {
	user := c.MustGet("user").(*models.UserModel)
	p := models.FindPostById(id)
	if p == nil {
		return gin.H{"message": "cannot find post by id"}, http.StatusNotFound
	}
	if p.Writer != user.Id {
		return gin.H{"message": "cannot access this resource"}, http.StatusForbidden
	}
	p.Title = req.Title
	p.Content = req.Content
	p.UpdateAt = time.Now()
	return PostResponse(p), http.StatusOK
}

func (r Posts) Delete(c *gin.Context, id int) (gin.H, int) {
	user := c.MustGet("user").(*models.UserModel)
	p := models.FindPostById(id)
	if p == nil {
		return gin.H{"message": "cannot find post by id"}, http.StatusNotFound
	}
	if p.Writer != user.Id {
		return gin.H{"message": "cannot access this resource"}, http.StatusForbidden
	}
	models.DeletePostById(id)
	return gin.H{}, http.StatusOK
}