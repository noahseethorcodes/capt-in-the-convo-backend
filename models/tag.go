package models

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Name     string   `gorm:"unique;not null"`
	IsActive bool     `gorm:"default:true"`
	Threads  []Thread `gorm:"many2many:thread_tags"` // Many-to-many relationship
}
