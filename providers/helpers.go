package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/tnqbao/gau_blog_service/config"
	"github.com/tnqbao/gau_blog_service/models"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

func GetUserByID(userID uint64) (*User, error) {
	url := fmt.Sprintf("http://localhost/api/user/public/%d", userID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch user: %s", resp.Status)
	}
	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func GetListUser(userIDs []uint64) ([]User, error) {
	url := "http://localhost/api/user/public/list"
	body, err := json.Marshal(map[string][]uint64{"ids": userIDs})
	if err != nil {
		return nil, fmt.Errorf("error marshalling userIDs: %w", err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error while sending request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch users: received status code %d", resp.StatusCode)
	}
	var users []User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	return users, nil
}

func ParseCacheValue(cacheValue string) int {
	if cacheValue == "" {
		return 0
	}
	value, err := strconv.ParseUint(cacheValue, 10, 32)
	if err != nil {
		return 0
	}
	return int(value)
}

func getCachedData(ctx context.Context, redisClient *redis.Client, key string, query func(db *gorm.DB) ([]interface{}, error)) ([]interface{}, error) {
	cache, err := redisClient.HGetAll(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("redis error: %v", err)
	}
	if len(cache) > 0 {
		var data []interface{}
		for _, itemStr := range cache {
			var item interface{}
			err := json.Unmarshal([]byte(itemStr), &item)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal blog response: %v", err)
			}
			data = append(data, item)
		}
		return data, nil
	}
	db := ctx.Value("db").(*gorm.DB) // Assuming db is passed in context
	blogs, err := query(db)
	if err != nil {
		return nil, fmt.Errorf("failed to query blogs: %v", err)
	}
	blogData := make(map[string]interface{})
	for _, blog := range blogs {
		blogJSON, err := json.Marshal(blog)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal blog response: %v", err)
		}
		blogData[strconv.Itoa(int(blog.(BlogResponse).ID))] = string(blogJSON)
	}

	err = redisClient.HSet(ctx, key, blogData).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to store data in Redis: %v", err)
	}
	err = redisClient.Expire(ctx, key, 10*time.Minute).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to set TTL for Redis cache: %v", err)
	}

	return blogs, nil
}

func GetListBlog(c *gin.Context, ctx context.Context, key string, query func(db *gorm.DB) ([]models.Blog, error)) ([]BlogResponse, error) {
	redisClient := config.GetRedisClient()
	db, exists := c.Get("db")
	if !exists || db == nil {
		return nil, fmt.Errorf("failed to retrieve DB from context")
	}

	gdb, ok := db.(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("failed to assert DB type")
	}

	cache, err := redisClient.HGetAll(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("redis error: %v", err)
	}

	var blogsResponse []BlogResponse
	if len(cache) > 0 {
		// Cache hit, process cache
		for _, blogStr := range cache {
			var blog BlogResponse
			err := json.Unmarshal([]byte(blogStr), &blog)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal blog response: %v", err)
			}

			upvote, downvote, err := getUpDownVotes(ctx, redisClient, blog.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch vote data: %v", err)
			}
			blog.Upvote = upvote
			blog.Downvote = downvote
			blogsResponse = append(blogsResponse, blog)
		}
		return blogsResponse, nil
	}

	// Cache miss, fetch from database
	blogs, err := query(gdb)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch blogs: %v", err)
	}

	// Process users and blogs
	var userIDs []uint64
	for _, blog := range blogs {
		userIDs = append(userIDs, blog.UserID)
	}

	users, err := GetListUser(userIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %v", err)
	}

	userMap := make(map[uint64]User)
	for _, user := range users {
		userMap[user.Id] = user
	}

	for i := range blogs {
		user, exists := userMap[blogs[i].UserID]
		if exists {
			upvote, downvote, err := getUpDownVotes(ctx, redisClient, blogs[i].ID)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch vote data: %v", err)
			}
			blogsResponse = append(blogsResponse, BlogResponse{
				ID:        blogs[i].ID,
				Title:     blogs[i].Title,
				Body:      blogs[i].Body,
				Upvote:    upvote,
				Downvote:  downvote,
				Comments:  blogs[i].Comments,
				CreatedAt: blogs[i].CreatedAt,
				User: User{
					Id:       user.Id,
					Fullname: user.Fullname,
				},
			})
		}
	}

	// Store the fetched blogs in Redis for future use
	blogData := make(map[string]interface{})
	for _, blogResp := range blogsResponse {
		blogJSON, err := json.Marshal(blogResp)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal blog response: %v", err)
		}
		blogData[strconv.Itoa(int(blogResp.ID))] = string(blogJSON)
	}

	err = redisClient.HSet(ctx, key, blogData).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to store data in Redis: %v", err)
	}

	err = redisClient.Expire(ctx, key, 10*time.Minute).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to set TTL for Redis cache: %v", err)
	}

	return blogsResponse, nil
}

func getUpDownVotes(ctx context.Context, redisClient *redis.Client, blogID uint64) (upvote, downvote int, err error) {
	updownKey := "blog:" + strconv.Itoa(int(blogID)) + ":up-down"
	updownCache, err := redisClient.HGetAll(ctx, updownKey).Result()
	if err != nil && err != redis.Nil {
		return 0, 0, fmt.Errorf("redis error: %v", err)
	}
	if err == redis.Nil {
		return 0, 0, nil
	}

	upvote = ParseCacheValue(updownCache["upvote"])
	downvote = ParseCacheValue(updownCache["downvote"])
	return
}
