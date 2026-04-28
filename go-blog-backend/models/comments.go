package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	PostID  uint   `json:"post_id" gorm:"not null;index"`
	Post    Post   `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Content string `json:"content" gorm:"type:varchar(500);not null"`
	Author  string `json:"author" gorm:"type:varchar(20);not null"`
}
