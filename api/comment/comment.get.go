package comment

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_blog_service/models"
	"gorm.io/gorm"
)

func GetComentsByBlogId(c *gin.Context) []models.Comment {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	if id == "" {
		id = "0"
	}
	var comments []models.Comment
	err := db.Transaction(func(tx *gorm.DB) error {
		if result := tx.Where("blog_id = ?", id).Find(&comments); result.Error != nil {
			return result.Error
		}
		return nil
	},
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "Cannot get comments: " + err.Error()})
		return []models.Comment{}
	}
	return comments
}
