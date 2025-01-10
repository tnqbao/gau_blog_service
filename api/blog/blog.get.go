package blog

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/tnqbao/gau_blog_service/config"
	"github.com/tnqbao/gau_blog_service/models"
	"github.com/tnqbao/gau_blog_service/providers"
	"gorm.io/gorm"
)

func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"status": code, "error": message})
}

func GetBlogByID(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	idStr := c.Param("id")
	ctx := context.Background()

	redisClient := config.GetRedisClient()
	key := "blog:" + idStr
	updownKey := "blog:" + idStr + ":up-down"

	var response providers.BlogResponse

	updownCache, err := redisClient.HGetAll(ctx, updownKey).Result()
	if err != nil && err != redis.Nil {
		respondWithError(c, http.StatusInternalServerError, "Redis error: "+err.Error())
		return
	}

	blogStr, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid blog ID format")
			return
		}

		var blog models.Blog
		if err := db.Preload("Comments").First(&blog, "id = ?", id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				respondWithError(c, http.StatusNotFound, "Blog not found")
			} else {
				respondWithError(c, http.StatusInternalServerError, err.Error())
			}
			return
		}

		user, err := providers.GetUserByID(blog.UserID)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "Failed to fetch author data")
			return
		}

		upvote := blog.Upvote + providers.ParseCacheValue(updownCache["upvote"])
		downvote := blog.Downvote + providers.ParseCacheValue(updownCache["downvote"])

		response = providers.BlogResponse{
			ID:        blog.ID,
			Title:     blog.Title,
			Body:      blog.Body,
			Upvote:    upvote,
			Downvote:  downvote,
			Comments:  blog.Comments,
			CreatedAt: blog.CreatedAt,
			User:      *user,
		}

		cacheData, err := json.Marshal(response)
		if err == nil {
			_ = redisClient.Set(ctx, key, cacheData, 1*time.Minute).Err()
		}
	} else if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to fetch data from cache")
		return
	} else {
		if err := json.Unmarshal([]byte(blogStr), &response); err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to parse cache data")
			return
		}

		response.Upvote += providers.ParseCacheValue(updownCache["upvote"])
		response.Downvote += providers.ParseCacheValue(updownCache["downvote"])
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": response})
}
