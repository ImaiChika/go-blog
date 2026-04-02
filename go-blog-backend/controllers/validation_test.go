package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterRejectsShortPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/register", Register)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(`{"username":"imai","password":"123"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusBadRequest, rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "密码长度不能少于6") {
		t.Fatalf("expected password validation error, got %s", rec.Body.String())
	}
}

func TestLoginRejectsBlankUsername(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/login", Login)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(`{"username":"   ","password":"123456"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusBadRequest, rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "用户名不能为空") {
		t.Fatalf("expected username validation error, got %s", rec.Body.String())
	}
}

func TestCreatePostRejectsBlankTitle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("username", "imai")
		c.Next()
	})
	router.POST("/api/v1/posts", CreatePost)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", bytes.NewBufferString(`{"title":"   ","content":"正文内容","cover_image":"http://localhost:8080/uploads/test.png"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusBadRequest, rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "标题不能为空") {
		t.Fatalf("expected title validation error, got %s", rec.Body.String())
	}
}

func TestGetPostsRejectsInvalidPageSize(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/api/v1/posts", GetPosts)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts?page=1&page_size=0", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusBadRequest, rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "每页数量不能小于1") {
		t.Fatalf("expected page_size validation error, got %s", rec.Body.String())
	}
}
