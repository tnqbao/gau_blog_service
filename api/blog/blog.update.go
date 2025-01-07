package blog

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_blog_service/models"
	"github.com/tnqbao/gau_blog_service/providers"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func UpdateBlogById(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var req providers.BlogRequest
	id := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusForbidden, "error": "Please login to access"})
		return
	}

	if id == "" {
		id = "0"
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"status": http.StatusBadRequest, "error": "Invalid request format: " + err.Error()})
		return
	}

	blog := models.Blog{
		Title: req.Title,
		Body:  req.Body,
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&models.Blog{}).Where("id = ? AND user_id = ?", id, userID).Updates(&blog)

		if result.Error != nil {
			return result.Error
		} else if result.RowsAffected == 0 {
			return fmt.Errorf("permission denied or blog not found! ")
		}

		return nil
	})

	if err != nil {
		log.Println("Transaction error:", err)
		if err.Error() == "permission denied or blog not found! " {
			c.JSON(http.StatusForbidden, gin.H{"status": http.StatusForbidden, "error": "permission denied or blog not found!"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot update blog: " + err.Error(), "status": http.StatusInternalServerError})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Update blog successfully!"})
}
