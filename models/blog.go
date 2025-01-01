package models

import "time"

type Blog struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement;column:id"`
	UserID    uint64    `gorm:"not null;column:user_id"`
	Title     string    `gorm:"not null;column:title"`
	Body      string    `gorm:"not null;column:body"`
	Upvote    int       `gorm:"default:0;column:upvote"`
	Downvote  int       `gorm:"default:0;column:downvote"`
	Comments  []Comment `gorm:"foreignKey:PostID;constraint:onDelete:CASCADE"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`
}
