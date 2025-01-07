package comment

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_blog_service/models"
	"gorm.io/gorm"
)

func DeleteCommentById(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	if id == "" {
		id = "0"
	}
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(403, gin.H{"status": 403, "error": "Please login to access"})
		return
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		if result := tx.Delete(&models.Comment{}).Where("id = ? AND user_id = ?", id, userID); result.Error != nil {
			return result.Error
		} else if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})

	if err != nil {
		c.JSON(500, gin.H{"error": "Cannot delete comment: " + err.Error()})
		return
	}

	c.JSON(200, "Delete comment successfully!")
}
