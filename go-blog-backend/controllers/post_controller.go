package controllers

import (
	"context"
	"fmt"
	"go-blog-backend/config"
	"go-blog-backend/models"
	"net/http"
	"strconv"
	"time"

	"encoding/json"

	"github.com/gin-gonic/gin"
)

var ctx = context.Background()

func getCurrentUsername(c *gin.Context) (string, bool) {
	usernameValue, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未找到用户信息，请重新登录"})
		return "", false
	}

	username, ok := usernameValue.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户信息格式错误"})
		return "", false
	}

	return username, true
}

func CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	username, ok := getCurrentUsername(c)
	if !ok {
		return
	}
	post.Author = username

	result := config.DB.Create(&post)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, post)
}
func GetPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var posts []models.Post
	var total int64
	config.DB.Model(&models.Post{}).Count(&total)
	config.DB.Limit(pageSize).Offset(offset).Find(&posts)

	c.JSON(http.StatusOK, gin.H{
		"data":      posts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetPostByID(c *gin.Context) {
	id := c.Param("id")
	redisKey := fmt.Sprintf("post:%s", id)
	var post models.Post
	val, err := config.RDB.Get(ctx, redisKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &post)
		c.JSON(http.StatusOK, gin.H{"data": post, "source": "redis"})
		return
	}
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	data, _ := json.Marshal(post)
	config.RDB.Set(ctx, redisKey, data, time.Hour)
	c.JSON(http.StatusOK, gin.H{"data": post, "source": "db"})
}
func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	username, ok := getCurrentUsername(c)
	if !ok {
		return
	}
	if post.Author != username {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权修改他人的文章"})
		return
	}

	var input struct {
		Title      string `json:"title"`
		Content    string `json:"content"`
		CoverImage string `json:"cover_image"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	post.Title = input.Title
	post.Content = input.Content
	post.CoverImage = input.CoverImage
	config.DB.Save(&post)
	redisKey := fmt.Sprintf("post:%s", id)
	config.RDB.Del(ctx, redisKey)
	c.JSON(http.StatusOK, post)
}
func DeletePost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	username, ok := getCurrentUsername(c)
	if !ok {
		return
	}
	if post.Author != username {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权删除他人的文章"})
		return
	}

	config.DB.Delete(&post)
	redisKey := fmt.Sprintf("post:%s", id)
	config.RDB.Del(ctx, redisKey)
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
