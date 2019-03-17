package views

import (
	"github.com/gin-gonic/gin"
	"github.com/hwangseonu/gin-restful"
	"github.com/hwangseonu/gin-restful-example/basic/server/models"
	"github.com/hwangseonu/gin-restful-example/basic/server/security"
	"net/http"
	"strconv"
	"time"
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
	if p, err := getPost(c); err != nil {
		return gin.H{"message": err.Message}, err.Status
	} else {
		comment := models.NewComment(req.Content, u)
		p.AddComment(comment)
		return PostResponse(p), http.StatusCreated
	}
}

func (r Comments) Patch(c *gin.Context, cid int, req CommentRequest) (gin.H, int) {
	u := c.MustGet("user").(*models.UserModel)
	if p, err := getPost(c); err != nil {
		return gin.H{"message": err.Message}, err.Status
	} else {
		if comment := p.FindComment(cid); comment == nil {
			return gin.H{"message": err.Error()}, http.StatusBadRequest
		} else {
			if comment.Id != u.Id {
				return gin.H{"message": "cannot access this resource"}, http.StatusForbidden
			}
			comment.Content = req.Content
			comment.UpdateAt = time.Now()
			return PostResponse(p), http.StatusCreated
		}
	}
}

func (r Comments) Delete(c *gin.Context, cid int) (gin.H, int) {
	u := c.MustGet("user").(*models.UserModel)
	if p, err := getPost(c); err != nil {
		return gin.H{"message": err.Message}, err.Status
	} else {
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

func getPost(c *gin.Context) (*models.PostModel, *gin_restful.ApplicationError) {
	str := c.Param("pid")
	if pid, err := strconv.Atoi(str); err != nil {
		return nil, &gin_restful.ApplicationError{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
	} else {
		p := models.FindPostById(pid)
		if p == nil {
			return nil, &gin_restful.ApplicationError{
				Message: "cannot find post by id",
				Status:  http.StatusNotFound,
			}
		} else {
			return p, nil
		}
	}
}
