package comment

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_blog_service/models"
	"github.com/tnqbao/gau_blog_service/providers"
	"gorm.io/gorm"
)

func UpdateCommentById(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	var req providers.CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"status": 400, "error": "Invalid request format: " + err.Error()})
		return
	}

	if id == "" {
		id = "0"
	}
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(403, gin.H{"status": 403, "error": "Please login to access"})
		return
	}

	comment := models.Comment{
		Body: req.Body,
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		if result := tx.Model(&models.Comment{}).Where("id = ? AND user_id = ?", id, userID).Updates(&comment); result.Error != nil {
			return result.Error
		} else if result.RowsAffected == 0 {
			return fmt.Errorf("comment not found or permission denied: ")
		}
		return nil
	},
	)

	if err != nil {
		if err.Error() == "comment not found or permission denied: " {
			c.JSON(403, gin.H{"status": 403, "error": "comment not found or permission denied"})
		} else {
			c.JSON(500, gin.H{"error": "Cannot update comment: " + err.Error()})
		}
		return
	}

	c.JSON(200, "Update comment successfully!")
}
