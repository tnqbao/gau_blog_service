package blog

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_blog_service/models"
	"gorm.io/gorm"
	"net/http"
)

func DeleteBlogById(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	userID, er := c.Get("user_id")
	if !er {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusForbidden, "error": "Please login to access"})
		return
	}

	if id == "" {
		id = "0"
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Blog{})
		if result.Error != nil {
			return result.Error
		} else if result.RowsAffected == 0 {
			return fmt.Errorf("permission denied or blog not found")
		}
		return nil
	})

	if err != nil {
		if err.Error() == "permission denied or blog not found" {
			c.JSON(http.StatusForbidden, gin.H{"error": "permission denied or blog not found!"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": "Cannot delete blog: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Delete blog successfully!"})
}
