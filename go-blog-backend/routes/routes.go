package routes

import (
	"go-blog-backend/controllers"
	"go-blog-backend/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// 【新增】：配置静态文件路由映射
	// 第一个参数是浏览器访问的路径，第二个参数是服务器本地的文件夹路径
	r.Static("/uploads", "./uploads")
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
		// 【新增】：图片上传接口（防止陌生人往你的服务器乱传文件，必须登录）
		authorized.POST("/upload", controllers.UploadImage)
	}

	return r
}
