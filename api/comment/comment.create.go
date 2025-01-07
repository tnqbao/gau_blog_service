package comment

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_blog_service/models"
	"github.com/tnqbao/gau_blog_service/providers"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func CreateComment(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	fullname, exists := c.Get("fullname")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "error": "Please login to access"})
		return
	}
	if fullname == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "error": "Please login to access"})
		return
	}
	var req providers.CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": "Invalid request format: " + err.Error()})
		return
	}
	tokenId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "error": "Please login to access"})
		return
	}

	tokenIdUint, ok := tokenId.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": "Invalid user_id format in create comment"})
		return
	}

	var comment = models.Comment{
		UserID:       tokenIdUint,
		UserFullName: fullname.(string),
		BlogID:       req.BlogID,
		Body:         req.Body,
		CreatedAt:    time.Now(),
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if result := tx.Create(&comment); result.Error != nil {
			return result.Error
		} else if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": "Cannot create comment: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Create comment successfully!"})
}
