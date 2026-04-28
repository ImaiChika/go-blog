package controllers

import (
	"go-blog-backend/config"
	"go-blog-backend/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type commentRequest struct {
	Content string `json:"content" binding:"required,max=500"`
}
type commentListQuery struct {
	Page     int `form:"page" binding:"min=1"`
	PageSize int `form:"page_size" binding:"min=1,max=50"`
}

func GetComments(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文章 ID 不正确"})
		return
	}
	query := commentListQuery{
		Page:     1,
		PageSize: 10,
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindErrorMessage(err)})
		return
	}
	var post models.Post
	if err := config.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文章不存在"})
		return
	}
	offset := (query.Page - 1) * query.PageSize
	var comments []models.Comment
	var total int64
	config.DB.Model(&models.Comment{}).Where("post_id=?", postID).Count(&total)
	config.DB.Where("post_id=?", postID).
		Order("created_at ASC").
		Limit(query.PageSize).
		Offset(offset).
		Find(&comments)
	c.JSON(http.StatusOK, gin.H{
		"data":      comments,
		"total":     total,
		"page":      query.Page,
		"page_size": query.PageSize,
	})
}

func CreateComment(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文章 ID 不正确"})
		return
	}
	var post models.Post
	if err := config.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文章不存在"})
		return
	}
	var input commentRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindErrorMessage(err)})
		return
	}
	input.Content = strings.TrimSpace(input.Content)
	if input.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "评论内容不能为空"})
		return
	}
	username, ok := getCurrentUsername(c)
	if !ok {
		return
	}
	comment := models.Comment{
		PostID:  uint(postID),
		Content: input.Content,
		Author:  username,
	}
	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建评论失败"})
		return
	}
	c.JSON(http.StatusOK, comment)
}
func DeleteComment(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "评论 ID 不正确"})
		return
	}
	var comment models.Comment
	if err := config.DB.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "评论不存在"})
		return
	}
	var post models.Post
	if err := config.DB.First(&post, comment.PostID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文章不存在"})
		return
	}
	username, ok := getCurrentUsername(c)
	if !ok {
		return
	}
	if comment.Author != username && post.Author != username {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无权删除评论"})
		return
	}
	config.DB.Delete(&comment)
	c.JSON(http.StatusOK, gin.H{"message": "评论删除成功"})

}
