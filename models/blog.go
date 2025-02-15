package models

import "time"

type Blog struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	UserID    uint64    `gorm:"not null;index;column:user_id" json:"userID"`
	Tag       string    `gorm:"column:tag" json:"tag"`
	Title     string    `gorm:"not null;index:idx_title;column:title" json:"title"`
	Body      string    `gorm:"not null;column:body" json:"body"`
	Upvote    int       `gorm:"default:0;column:upvote" json:"upvote"`
	Downvote  int       `gorm:"default:0;column:downvote" json:"downvote"`
	Comments  []Comment `gorm:"foreignKey:BlogID;constraint:onDelete:CASCADE" json:"comments"`
	CreatedAt time.Time `gorm:"autoCreateTime;index;column:created_at" json:"createdAt"`
}
