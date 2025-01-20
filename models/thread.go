package models

import "gorm.io/gorm"

type Thread struct {
	gorm.Model
	Title   string `gorm:"not null"`
	Content string `gorm:"type:text;not null"`
	UserID  uint   `gorm:"not null"`              // Foreign key
	Tags    []Tag  `gorm:"many2many:thread_tags"` // Many-to-many relationship
}
