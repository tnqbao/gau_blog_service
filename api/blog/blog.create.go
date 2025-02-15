package blog

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_blog_service/lib/ai"
	"github.com/tnqbao/gau_blog_service/models"
	"github.com/tnqbao/gau_blog_service/providers"
	"gorm.io/gorm"
)

func CreateBlog(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var req providers.BlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("UserRequest binding error:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "Invalid request format: " + err.Error(),
		})
		return
	}

	tokenID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  "Please login to access",
		})
		return
	}

	tokenIDUint, ok := tokenID.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Invalid user_id format in create blog",
		})
		return
	}

	blog := models.Blog{
		UserID:    tokenIDUint,
		Title:     req.Title,
		Body:      req.Body,
		CreatedAt: time.Now(),
	}

	blog.Tag = ai.DetectTag(blog.Body)

	err := db.Transaction(func(tx *gorm.DB) error {
		if result := tx.Create(&blog); result.Error != nil {
			return result.Error
		}
		return nil
	})

	if err != nil {
		log.Println("Transaction error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Cannot create blog: " + err.Error(),
		})
		return
	}

	response := providers.BriefBlogResponse{
		ID:        blog.ID,
		Title:     blog.Title,
		Body:      blog.Body,
		Tag:       blog.Tag,
		CreatedAt: blog.CreatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Create blog successfully!",
		"blog":    response,
	})
}
