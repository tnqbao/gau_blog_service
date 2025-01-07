package providers

import (
	"time"

	"github.com/tnqbao/gau_blog_service/models"
)

type User struct {
	Id       uint16 `json:"id"`
	Fullname string `json:"fullname"`
}

type BlogRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}
type BlogResponse struct {
	ID        uint64             `json:"id"`
	Title     string             `json:"title"`
	Body      string             `json:"body"`
	Upvote    int                `json:"upvote"`
	Downvote  int                `json:"downvote"`
	Comments  [](models.Comment) `json:"comments"`
	CreatedAt time.Time          `json:"createdAt"`
	User      User               `json:"user"`
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
