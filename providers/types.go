package providers

import (
	"time"
)

type User struct {
	Id       uint64 `json:"user_id"`
	Fullname string `json:"fullname"`
}

type BlogRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}
type BlogResponse struct {
	ID        uint64    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Upvote    int       `json:"upvote"`
	Downvote  int       `json:"downvote"`
	CreatedAt time.Time `json:"createdAt"`
	User      User      `json:"user"`
}

type BriefBlogResponse struct {
	ID        uint64    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
}

type CommentRequest struct {
	BlogID uint64 `json:"blog_id"`
	Body   string `json:"body"`
}

type ListBlogResponse struct {
	Blogs       []BlogResponse `json:"blogs"`
	ItemPerPage int            `json:"itemPerPage"`
	TotalItem   int            `json:"totalItem"`
	Page        int            `json:"page"`
}
