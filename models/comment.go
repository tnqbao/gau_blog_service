package models

import "time"

type Comment struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement;column:id"`
	Body      string    `gorm:"not null"`
	UserID    uint      `gorm:"not null"`
	BlogID    uint      `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Blog      Blog      `gorm:"foreignKey:BlogID;constraint:OnDelete:CASCADE"`
}
