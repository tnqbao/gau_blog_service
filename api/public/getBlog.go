package public

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_blog_service/models"
	"github.com/tnqbao/gau_blog_service/providers"
	"gorm.io/gorm"
)

func GetBlogById(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}
	var blog models.Blog
	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&blog, "id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	user, err := providers.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lost author data of blog"})
	}
	response := providers.BlogResponse{
		ID:        blog.ID,
		Title:     blog.Title,
		Body:      blog.Body,
		Upvote:    blog.Upvote,
		Downvote:  blog.Downvote,
		Comments:  blog.Comments,
		CreatedAt: blog.CreatedAt,
		User:      *user,
	}
	c.JSON(http.StatusOK, response)
}
