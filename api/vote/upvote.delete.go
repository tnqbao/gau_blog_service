package vote

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_blog_service/config"
	"github.com/tnqbao/gau_blog_service/models"
	"gorm.io/gorm"
	"strconv"
)

func DeleteUpvoteByBlogID(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	userID, er := c.Get("user_id")
	if !er {
		c.JSON(401, gin.H{"status": 401, "error": "Please login to access"})
		return
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("blog_id = ? AND user_id = ?", id, userID).Delete(&models.Vote{})
		if result.Error != nil {
			return result.Error
		} else if result.RowsAffected == 0 {
			return nil
		}

		redisClient := config.GetRedisClient()
		key := "blog:" + id + ":up-down"
		cache, _ := redisClient.HGetAll(context.Background(), key).Result()

		upvote := 0
		downvote := 0
		if len(cache) > 0 {
			upvote, _ = strconv.Atoi(cache["upvote"])
			downvote, _ = strconv.Atoi(cache["downvote"])
		}

		upvote--
		err := redisClient.Del(context.Background()).Err()
		err = redisClient.HSet(context.Background(), key, "upvote", upvote, "downvote", downvote).Err()
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(500, gin.H{"status": 500, "error": "Cannot delete upvote: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": 200, "message": "Delete upvote successfully!"})
}
