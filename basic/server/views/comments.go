package views

import (
	"github.com/gin-gonic/gin"
	"github.com/hwangseonu/gin-restful"
	"github.com/hwangseonu/gin-restful-example/basic/server/models"
	"github.com/hwangseonu/gin-restful-example/basic/server/security"
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
	r.AddMiddleware(security.AuthRequired(security.ACCESS, "ROLE_USER"), http.MethodPost, http.MethodPatch)
	return r
}

func (r Comments) Post(c *gin.Context, req CommentRequest) (gin.H, int) {
	u := c.MustGet("user").(*models.UserModel)
	if pid, err := getPostId(c); err != nil {
		return gin.H{"message": err.Error()}, http.StatusBadRequest
	} else {
		p := models.FindPostById(pid)
		if p == nil {
			return gin.H{"message": "cannot find post by id"}, http.StatusNotFound
		}
		c := models.NewComment(req.Content, u)
		p.AddComment(c)
		return PostResponse(p), http.StatusCreated
	}
}

func (r Comments) Patch(c *gin.Context, cid int, req CommentRequest) (gin.H, int) {
	u := c.MustGet("user").(*models.UserModel)
	if pid, err := getPostId(c); err != nil {
		return gin.H{"message": err.Error()}, http.StatusBadRequest
	} else {
		p := models.FindPostById(pid)
		if p == nil {
			return gin.H{"message": "cannot find post by id"}, http.StatusNotFound
		}
		if comment := p.FindComment(cid); comment == nil {
			return gin.H{"message": err.Error()}, http.StatusBadRequest
		} else {
			if comment.Id != u.Id {
				return gin.H{"message": "cannot access this resource"}, http.StatusForbidden
			}
			comment.Content = req.Content
			return PostResponse(p), http.StatusCreated
		}
	}
}

func (r Comments) Delete(c *gin.Context, cid int) (gin.H, int) {
	u := c.MustGet("user").(*models.UserModel)
	if pid, err := getPostId(c); err != nil {
		return gin.H{"message": err.Error()}, http.StatusBadRequest
	} else {
		p := models.FindPostById(pid)
		if p == nil {
			return gin.H{"message": "cannot find post by id"}, http.StatusNotFound
		}
		if comment := p.FindComment(cid); comment == nil {
			return gin.H{"message": err.Error()}, http.StatusNotFound
		} else {
			if comment.Id != u.Id {
				return gin.H{"message": "cannot access this resource"}, http.StatusForbidden
			}
			if ok := p.RemoveComment(comment.Id); !ok {
				return gin.H{"message": "cannot find comment by id"}, http.StatusNotFound
			}
			return PostResponse(p), http.StatusOK
		}
	}
}

func getPostId(c *gin.Context) (int, error) {
	pid := c.Param("pid")
	return strconv.Atoi(pid)
}
