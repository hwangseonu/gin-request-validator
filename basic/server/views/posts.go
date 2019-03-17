package views

import (
	"github.com/gin-gonic/gin"
	"github.com/hwangseonu/gin-restful"
	"github.com/hwangseonu/gin-restful-example/basic/server/models"
	"github.com/hwangseonu/gin-restful-example/basic/server/security"
	"net/http"
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
	r.AddMiddleware(security.AuthRequired(security.ACCESS), http.MethodPost)
	return r
}

func (r Posts) Get(id int) (gin.H, int) {
	p := models.FindPostById(id)
	if p == nil {
		return gin.H{"message": "cannot find post by id"}, http.StatusNotFound
	}
	return createPostResponse(p), http.StatusOK
}

func (r Posts) Post(c *gin.Context, req CreatePostRequest) (gin.H, int) {
	user := c.MustGet("user").(*models.UserModel)
	p := models.NewPost(req.Title, req.Content, user)
	p.Save()
	return createPostResponse(p), http.StatusCreated
}

func createPostResponse(p *models.PostModel) gin.H {
	u := models.FindUserById(p.Writer)
	writer := gin.H{}
	if u != nil {
		writer = gin.H{
			"username": u.Username,
			"nickname": u.Nickname,
			"email":    u.Email,
		}
	}
	return gin.H{
		"id":        p.Id,
		"title":     p.Title,
		"content":   p.Content,
		"comments":  createCommentsResponse(p.Comments),
		"writer":    writer,
		"create_at": p.CreateAt,
		"update_at": p.UpdateAt,
	}
}

func createCommentsResponse(c []*models.CommentModel) []gin.H {
	result := make([]gin.H, 0)
	for _, v := range c {
		u := models.FindUserById(v.Writer)
		writer := gin.H{}
		if u != nil {
			writer = gin.H{
				"username": u.Username,
				"nickname": u.Nickname,
				"email":    u.Email,
			}
		}
		result = append(result, gin.H{
			"id":        v.Id,
			"content":   v.Content,
			"writer":    writer,
			"create_at": v.CreateAt,
			"update_at": v.UpdateAt,
		})
	}
	return result
}
