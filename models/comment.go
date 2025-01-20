package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content  string `gorm:"type:text;not null"`
	UserID   uint   `gorm:"not null"` // Foreign key referencing User
	ThreadID uint   `gorm:"not null"` // Foreign key referencing Thread
}
