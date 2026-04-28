package config

import (
	"fmt"
	"go-blog-backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := GetDBDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败" + err.Error())
	}
	db.AutoMigrate(&models.Post{}, &models.User{}, &models.Comment{})
	DB = db
	fmt.Println("数据库连接完成并自动迁移完成！")
}
