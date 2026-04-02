package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

const maxUploadSize = 5 << 20 // 5MB
//	func UploadImage(c *gin.Context) {
//		file, err := c.FormFile("file")
//		if err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{
//				"error":  "文件上传失败，请使用 multipart/form-data 并确保字段名为 file",
//				"detail": err.Error(),
//			})
//			return
//		}
//		uploadDir := "./uploads"
//		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建上传目录失败"})
//			return
//		}
//		ext := filepath.Ext(file.Filename)
//		newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
//		// 拼接完整的保存路径 (例如: ./uploads/168123456789.jpg)
//		dst := filepath.Join(uploadDir, newFileName)
//		if err := c.SaveUploadedFile(file, dst); err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
//			return
//		}
//		imageUrl := fmt.Sprintf("http://localhost:8080/uploads/%s", newFileName)
//		c.JSON(http.StatusOK, gin.H{"message": "文件上传成功", "url": imageUrl})
//	}
func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "文件上传失败，请使用 multipart/form-data 并确保字段名为 file",
			"detail": err.Error(),
		})
		return
	}
	if file.Size > maxUploadSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "文件大小不能超过 5MB",
			"detail": fmt.Sprintf("当前文件大小: %d", file.Size),
		})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "打开文件失败"})
		return
	}
	defer src.Close()
	buffer := make([]byte, 512)
	n, err := src.Read(buffer)
	if err != nil && err != io.EOF {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取文件失败"})
		return
	}
	contentType := http.DetectContentType(buffer[:n])
	allowedTypes := map[string]string{
		"image/jpeg": ".jpg",
		"image/png":  ".png",
	}
	ext, ok := allowedTypes[contentType]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只允许上传 jpg 和 png 文件类型"})
		return
	}
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建上传目录失败"})
		return
	}
	newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join(uploadDir, newFileName)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}
	imageURL := fmt.Sprintf("http://localhost:8080/uploads/%s", newFileName)
	c.JSON(http.StatusOK, gin.H{"message": "文件上传成功", "url": imageURL})

}
