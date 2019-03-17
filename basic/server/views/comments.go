package views

import (
	"github.com/gin-gonic/gin"
	"github.com/hwangseonu/gin-restful"
	"github.com/hwangseonu/gin-restful-example/basic/server/models"
	"net/http"
	"strconv"
)

type Comments struct {
	*gin_restful.Resource
}

type CommentRequest struct {
	Content string `json:"content" binding:"required,notblank"`
}

func InitCommentsResource() Comments {
	r := Comments{gin_restful.InitResource()}
	return r
}

func (r Comments) Post(c *gin.Context, req CommentRequest) (gin.H, int) {
	u := c.MustGet("user").(*models.UserModel)
	if pid, err := getPostId(c); err != nil {
		return gin.H{"message": err.Error()}, http.StatusBadRequest
	} else {
		p := models.FindPostById(pid)
		c := models.NewComment(req.Content, u)
		p.AddComment(c)
		p.Save()
		return PostResponse(p), http.StatusCreated
	}
}

func getPostId(c *gin.Context) (int, error) {
	pid := c.Param("pid")
	return strconv.Atoi(pid)
}
