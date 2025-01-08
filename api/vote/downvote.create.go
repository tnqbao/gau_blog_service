package vote

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_blog_service/config"
	"github.com/tnqbao/gau_blog_service/models"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func AddDownVoteByBlogID(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	ctx := context.Background()
	redisClient := config.GetRedisClient()
	key := "blog:" + id + ":up-down"

	cache, err := redisClient.HGetAll(ctx, key).Result()
	if err != nil {
		c.JSON(500, gin.H{"status": 500, "error": "Redis error: " + err.Error()})
		return
	}

	upvote := 0
	downvote := 0

	if len(cache) == 0 {
		err := redisClient.HSet(ctx, key, "upvote", upvote, "downvote", downvote).Err()
		if err != nil {
			c.JSON(500, gin.H{"status": 500, "error": "Redis error: " + err.Error()})
			return
		}
		err = redisClient.Expire(ctx, key, 5*time.Minute).Err()
		if err != nil {
			c.JSON(500, gin.H{"status": 500, "error": "Redis error: " + err.Error()})
			return
		}
	} else {
		upvote, _ = strconv.Atoi(cache["upvote"])
		downvote, _ = strconv.Atoi(cache["downvote"])
	}

	downvote++

	err = redisClient.HSet(ctx, key, "upvote", upvote, "downvote", downvote).Err()
	if err != nil {
		c.JSON(500, gin.H{"status": 500, "error": "Redis error: " + err.Error()})
		return
	}

	userID, er := c.Get("user_id")
	if !er {
		c.JSON(401, gin.H{"status": 401, "error": "Please login to access"})
		return
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		updateResult := tx.Model(&models.Vote{}).Where("blog_id = ? AND user_id = ?", id, userID).Update("state", 1)
		if updateResult.Error != nil {
			return updateResult.Error
		} else if updateResult.RowsAffected == 0 {
			blogID, err := strconv.ParseUint(id, 10, 64)
			if err != nil {
				return err
			}
			createResult := tx.Create(&models.Vote{BlogID: blogID, UserID: userID.(uint64), State: true})
			if createResult.Error != nil {
				return createResult.Error
			}
		}
		return nil
	})

	if err != nil {
		if err.Error() == "blog not found" {
			c.JSON(404, gin.H{"status": 404, "error": "blog not found"})
		} else {
			c.JSON(500, gin.H{"status": 500, "error": "Failed to update upvote"})
		}
		return
	}

	c.JSON(200, gin.H{"status": 200, "message": "Upvote successfully!"})
}
