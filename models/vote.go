package models

import "time"

type Vote struct {
	UserID    uint64    `gorm:"not null;column:user_id" json:"user_id"`
	BlogID    uint64    `gorm:"not null;primaryKey;index;column:blog_id;constraint:onDelete:CASCADE" json:"blog_id"`
	State     bool      `gorm:"column:state" json:"state"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
}
