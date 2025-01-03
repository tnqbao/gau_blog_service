package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_blog_service/api/authed"
	"github.com/tnqbao/gau_blog_service/api/public"
	"github.com/tnqbao/gau_blog_service/middlewares"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.CORSMiddleware())
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})
	apiRoutes := r.Group("/api")
	{
		forumRoutes := apiRoutes.Group("/forum")
		{
			authedRoutes := forumRoutes.Group("/authed")
			{
				authedRoutes.Use(middlewares.AuthMiddleware())
				authedRoutes.PUT("/blog", authed.CreateBlog)
			}
			publicRoutes := forumRoutes.Group("/public")
			{
				publicRoutes.GET("/blog/:id", public.GetBlogById)
			}
		}
	}
	return r
}
