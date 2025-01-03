package authed

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_blog_service/models"
	"github.com/tnqbao/gau_blog_service/providers"
	"gorm.io/gorm"
)

func CreateBlog(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var req providers.BlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("UserRequest binding error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}
	tokenId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please login to access"})
		return
	}
	tokenIdUint, ok := tokenId.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user_id format in create blog"})
		return
	}
	blog := models.Blog{
		UserID:    tokenIdUint,
		Title:     req.Title,
		Body:      req.Body,
		CreatedAt: time.Now(),
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		if result := tx.Create(&blog); result.Error != nil {
			return result.Error
		} else if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
	if err != nil {
		log.Println("Transaction error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot create blog: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Create blog successfully!")
}
