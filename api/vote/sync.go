package vote

import (
	"context"
	"fmt"
	"github.com/tnqbao/gau_blog_service/config"
	"github.com/tnqbao/gau_blog_service/models"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func SyncVotesToDatabase(db *gorm.DB) {
	ctx := context.Background()
	redisClient := config.GetRedisClient()

	keys, err := redisClient.Keys(ctx, "blog:*:up-down").Result()
	if err != nil {
		fmt.Println("Error fetching keys from Redis:", err)
		return
	}

	for _, key := range keys {
		cache, err := redisClient.HGetAll(ctx, key).Result()
		if err != nil {
			fmt.Println("Error fetching cache data for key:", key, err)
			continue
		}

		upvote := 0
		downvote := 0
		if len(cache) > 0 {
			upvote, _ = strconv.Atoi(cache["upvote"])
			downvote, _ = strconv.Atoi(cache["downvote"])
		}

		var blogID int
		_, err = fmt.Sscanf(key, "blog:%d:up-down", &blogID)
		if err != nil {
			fmt.Println("Error parsing blog ID from key:", key, err)
			continue
		}

		err = db.Transaction(func(tx *gorm.DB) error {
			updateResult := tx.Model(&models.Blog{}).Where("id = ?", blogID).Updates(map[string]interface{}{
				"upvote":   gorm.Expr("upvote + ?", upvote),
				"downvote": gorm.Expr("downvote + ?", downvote),
			})
			if updateResult.Error != nil {
				return updateResult.Error
			}
			return nil
		})
		if err != nil {
			fmt.Println("Error updating database for blog ID:", blogID, err)
			continue
		}

		err = redisClient.Del(ctx, key).Err()
		if err != nil {
			fmt.Println("Error deleting cache for key:", key, err)
			continue
		}
		err = redisClient.Del(ctx, fmt.Sprintf("blog:%d:up", blogID)).Err()
		if err != nil {
			fmt.Println("Error deleting cache for key:", fmt.Sprintf("blog:%d", blogID), err)
			continue
		}

		fmt.Printf("Successfully synced votes for blog ID: %d\n", blogID)
	}
}

func StartSyncJob(db *gorm.DB) {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		SyncVotesToDatabase(db)
	}
}
