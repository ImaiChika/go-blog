package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model        // GORM 自动添加 ID, CreatedAt, UpdatedAt, DeletedAt
	Title      string `json:"title" gorm:"type:varchar(100);not null"`
	Content    string `json:"content" gorm:"type:text;not null"`
	Author     string `json:"author" gorm:"type:varchar(20);default:'Anonymous'"`
	CoverImage string `json:"cover_image" gorm:"type:varchar(255)"` // 新增封面图字段，保存图片的 URL
}

//gorm.Model：这是一个嵌套结构体，帮我们自动管理主键和时间戳。

//结构体标签（如 json:"title"）：告诉 Go 在转成 JSON 给前端时叫什么名字。

//gorm:"..."：告诉数据库这一列的约束条件。
