package feeds

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_blog_service/models"
	"github.com/tnqbao/gau_blog_service/providers"
	"gorm.io/gorm"
	"strconv"
)

func GetTredingPerPage(c *gin.Context) {
	pageStr := c.Param("page")
	ctx := context.Background()
	key := "/api/forum/treding/" + pageStr

	if pageStr == "" {
		c.JSON(400, gin.H{"status": 400, "error": "Please provide page number"})
		return
	}

	query := func(db *gorm.DB) ([]models.Blog, error) {
		var blogs []models.Blog
		page, err := strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			return nil, fmt.Errorf("invalid page format")
		}

		err = db.Order("upvote desc").
			Offset((page - 1) * 12).
			Limit(12).
			Find(&blogs).Error
		if err != nil {
			return nil, err
		}

		return blogs, nil
	}

	blogsResponse, err := providers.GetListBlog(c, ctx, key, query)
	if err != nil {
		c.JSON(500, gin.H{"status": 500, "error": err.Error()})
		return
	}
	if blogsResponse == nil {
		return
	}
	c.JSON(200, gin.H{"status": 200, "data": blogsResponse})
}
