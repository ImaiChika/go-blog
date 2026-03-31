package routes

import (
	"go-blog-backend/controllers"
	"go-blog-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 身份验证路由
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// 公开的文章读取路由
	v1 := r.Group("/api/v1")
	{
		v1.GET("/posts", controllers.GetPosts)        // 现在支持分页了
		v1.GET("/posts/:id", controllers.GetPostByID) // 带 Redis 缓存
	}

	// 受保护的文章管理路由
	authorized := r.Group("/api/v1")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.POST("/posts", controllers.CreatePost)       // 创建
		authorized.PUT("/posts/:id", controllers.UpdatePost)    // 更新
		authorized.DELETE("/posts/:id", controllers.DeletePost) // 删除
	}

	return r
}
