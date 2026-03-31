package main

import (
	"go-blog-backend/config"
	"go-blog-backend/routes"
)

func main() {
	config.InitDB() //1.初始化数据库
	config.InitRedis()
	r := routes.SetupRouter() //2.设置路由
	r.Run(":8080")
}
