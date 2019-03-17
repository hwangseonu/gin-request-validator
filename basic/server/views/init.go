package views

import (
	"github.com/gin-gonic/gin"
	"github.com/hwangseonu/gin-restful"
	"github.com/hwangseonu/gin-restful-example/basic/server/models"
)

func RegisterViews(e *gin.Engine) {
	users := InitUsersResource()
	auth := InitAuthResource()
	posts := InitPostsResource()
	comments := InitCommentsResource()
	api := gin_restful.NewApi(e, "/")
	api.AddResource(users, "/users")
	api.AddResource(auth, "/auth")
	api.AddResource(posts, "/posts")
	api.AddResource(comments, "/posts/:pid/comments")
}

func UserResponse(u *models.UserModel) gin.H {
	if u != nil {
		return gin.H{
			"username": u.Username,
			"nickname": u.Nickname,
			"email":    u.Email,
		}
	} else {
		return gin.H{}
	}
}

func AuthResponse(a, r string) gin.H {
	if r != "" {
		return gin.H{"access":  a, "refresh": r}
	} else {
		return gin.H{"access":  a}
	}
}

func PostResponse(p *models.PostModel) gin.H {
	u := models.FindUserById(p.Writer)
	return gin.H{
		"id":        p.Id,
		"title":     p.Title,
		"content":   p.Content,
		"comments":  CommentsResponse(p.Comments),
		"writer":    UserResponse(u),
		"create_at": p.CreateAt,
		"update_at": p.UpdateAt,
	}
}

func CommentsResponse(c []*models.CommentModel) []gin.H {
	result := make([]gin.H, 0)
	for _, v := range c {
		u := models.FindUserById(v.Writer)
		result = append(result, gin.H{
			"id":        v.Id,
			"content":   v.Content,
			"writer":    UserResponse(u),
			"create_at": v.CreateAt,
			"update_at": v.UpdateAt,
		})
	}
	return result
}
