package models

import "time"

type PostModel struct {
	Id       int             `json:"id"`
	Title    string          `json:"title"`
	Content  string          `json:"content"`
	Comments []*CommentModel `json:"comments"`
	Writer   int             `json:"writer"`
	CreateAt time.Time       `json:"create_at"`
	UpdateAt time.Time       `json:"update_at"`
}

type CommentModel struct {
	Id       int       `json:"id"`
	Content  string    `json:"content"`
	Writer   int       `json:"writer"`
	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
}

func (p *PostModel) Save() {
	posts[p.Id] = p
}

func (p *PostModel) AddComment(c *CommentModel) {
	p.Comments = append(p.Comments, c)
}

func (p *PostModel) RemoveComment(id int) bool {
	for i, c := range p.Comments {
		if c.Id == id {
			p.Comments = append(p.Comments[:i], p.Comments[i+1:]...)
			return true
		}
	}
	return false
}

func NewPost(title, content string, writer *UserModel) *PostModel {
	return &PostModel{
		Id:       GetNextId("posts"),
		Title:    title,
		Content:  content,
		Writer:   writer.Id,
		Comments: make([]*CommentModel, 0),
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}
}

func NewComment(content string, writer *UserModel) *CommentModel {
	return &CommentModel{
		Id:       GetNextId("comments"),
		Content:  content,
		Writer:   writer.Id,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}
}

func FindPostById(id int) *PostModel {
	p, ok := posts[id]
	if !ok {
		return nil
	}
	if !ExistsUserById(p.Writer) {
		DeletePostById(p.Id)
		return nil
	}
	for _, c := range p.Comments {
		if !ExistsUserById(c.Writer) {
			p.RemoveComment(c.Id)
		}
	}
	return p
}

func DeletePostById(id int) {
	delete(posts, id)
}

func ExistsPostById(id int) bool {
	return FindPostById(id) != nil
}
