package models

import "time"

type Comment struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement;column:id"`
	Body         string    `gorm:"not null" json:"body"`
	UserID       uint64    `gorm:"not null" json:"user_id"`
	UserFullName string    `json:"fullname"`
	BlogID       uint64    `gorm:"not null" json:"blog_id"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}
