package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-blog-backend/config"
	"go-blog-backend/models"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

var initPostControllerDeps sync.Once

func setupPostControllerTestDeps(t *testing.T) {
	t.Helper()

	initPostControllerDeps.Do(func() {
		config.InitDB()
		config.InitRedis()
	})
}

func createTestPost(t *testing.T, author string) models.Post {
	t.Helper()

	post := models.Post{
		Title:      fmt.Sprintf("title-%d", time.Now().UnixNano()),
		Content:    "original content",
		Author:     author,
		CoverImage: "http://localhost:8080/uploads/original.png",
	}
	if err := config.DB.Create(&post).Error; err != nil {
		t.Fatalf("failed to create test post: %v", err)
	}
	t.Cleanup(func() {
		config.DB.Unscoped().Delete(&models.Post{}, post.ID)
		config.RDB.Del(ctx, fmt.Sprintf("post:%d", post.ID))
	})

	return post
}

func newAuthRouter(username string) *gin.Engine {
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("username", username)
		c.Next()
	})
	router.PUT("/api/v1/posts/:id", UpdatePost)
	router.DELETE("/api/v1/posts/:id", DeletePost)
	return router
}

func TestUpdatePostRejectsNonAuthor(t *testing.T) {
	setupPostControllerTestDeps(t)
	gin.SetMode(gin.TestMode)

	post := createTestPost(t, "author_owner")
	router := newAuthRouter("another_user")

	body := bytes.NewBufferString(`{"title":"hacked title","content":"bad","cover_image":"http://localhost:8080/uploads/hacked.png","author":"another_user"}`)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/posts/%d", post.ID), body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusForbidden, rec.Code, rec.Body.String())
	}

	var updated models.Post
	if err := config.DB.First(&updated, post.ID).Error; err != nil {
		t.Fatalf("failed to fetch post after forbidden update: %v", err)
	}
	if updated.Title != post.Title || updated.Author != post.Author {
		t.Fatalf("post should not be modified by non-author, got %+v", updated)
	}
}

func TestUpdatePostKeepsAuthorFromExistingRecord(t *testing.T) {
	setupPostControllerTestDeps(t)
	gin.SetMode(gin.TestMode)

	post := createTestPost(t, "author_owner")
	router := newAuthRouter("author_owner")

	body := bytes.NewBufferString(`{"title":"updated title","content":"updated content","cover_image":"http://localhost:8080/uploads/updated.png","author":"fake_author"}`)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/posts/%d", post.ID), body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, rec.Code, rec.Body.String())
	}

	var resp models.Post
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if resp.Author != "author_owner" {
		t.Fatalf("expected author to remain author_owner, got %q", resp.Author)
	}

	var updated models.Post
	if err := config.DB.First(&updated, post.ID).Error; err != nil {
		t.Fatalf("failed to fetch updated post: %v", err)
	}
	if updated.Author != "author_owner" {
		t.Fatalf("expected stored author to remain author_owner, got %q", updated.Author)
	}
}

func TestDeletePostRejectsNonAuthor(t *testing.T) {
	setupPostControllerTestDeps(t)
	gin.SetMode(gin.TestMode)

	post := createTestPost(t, "delete_owner")
	router := newAuthRouter("another_user")

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/posts/%d", post.ID), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusForbidden, rec.Code, rec.Body.String())
	}

	var count int64
	if err := config.DB.Model(&models.Post{}).Where("id = ?", post.ID).Count(&count).Error; err != nil {
		t.Fatalf("failed to count post after forbidden delete: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected post to remain after forbidden delete, count=%d", count)
	}
}

func TestDeletePostAllowsAuthor(t *testing.T) {
	setupPostControllerTestDeps(t)
	gin.SetMode(gin.TestMode)

	post := createTestPost(t, "delete_owner")
	router := newAuthRouter("delete_owner")

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/posts/%d", post.ID), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, rec.Code, rec.Body.String())
	}

	var count int64
	if err := config.DB.Model(&models.Post{}).Where("id = ?", post.ID).Count(&count).Error; err != nil {
		t.Fatalf("failed to count post after delete: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected post to be deleted, count=%d", count)
	}
}

func TestGetPostsReturnsNewestFirst(t *testing.T) {
	setupPostControllerTestDeps(t)
	gin.SetMode(gin.TestMode)

	older := createTestPost(t, "sort_owner")
	newer := createTestPost(t, "sort_owner")

	olderTime := time.Date(2098, time.January, 1, 0, 0, 0, 0, time.UTC)
	newerTime := time.Date(2099, time.January, 1, 0, 0, 0, 0, time.UTC)

	if err := config.DB.Model(&models.Post{}).Where("id = ?", older.ID).Update("created_at", olderTime).Error; err != nil {
		t.Fatalf("failed to update older created_at: %v", err)
	}
	if err := config.DB.Model(&models.Post{}).Where("id = ?", newer.ID).Update("created_at", newerTime).Error; err != nil {
		t.Fatalf("failed to update newer created_at: %v", err)
	}

	router := gin.New()
	router.GET("/api/v1/posts", GetPosts)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts?page=1&page_size=2", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, rec.Code, rec.Body.String())
	}

	var resp struct {
		Data []models.Post `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(resp.Data) < 2 {
		t.Fatalf("expected at least two posts in response, got %d", len(resp.Data))
	}
	if resp.Data[0].ID != newer.ID || resp.Data[1].ID != older.ID {
		t.Fatalf("expected newest posts first, got ids %d then %d", resp.Data[0].ID, resp.Data[1].ID)
	}
}
