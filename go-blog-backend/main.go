package main

import (
	"go-blog-backend/config"
	"go-blog-backend/routes"
	"log"
)

func main() {
	if err := config.LoadEnv(); err != nil {
		log.Fatalf("加载 .env 失败: %v", err)
	}
	config.InitDB() //1.初始化数据库
	config.InitRedis()
	r := routes.SetupRouter() //2.设置路由
	r.Run(":8080")
}
