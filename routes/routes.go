package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_blog_service/api/blog"
	"github.com/tnqbao/gau_blog_service/api/comment"
	"github.com/tnqbao/gau_blog_service/api/feeds"
	"github.com/tnqbao/gau_blog_service/api/public"
	"github.com/tnqbao/gau_blog_service/api/vote"
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
		forumRoutes := apiRoutes.Group("/storisy")
		{
			blogRoutes := forumRoutes.Group("/blog")
			{
				authedBlogRoutes := blogRoutes.Group("/")
				{
					authedBlogRoutes.Use(middlewares.AuthMiddleware())
					authedBlogRoutes.POST("/:id", blog.UpdateBlogById)
					authedBlogRoutes.DELETE("/:id", blog.DeleteBlogById)
					authedBlogRoutes.PUT("/", blog.CreateBlog)

					voteRoutes := authedBlogRoutes.Group("/vote")
					{
						voteRoutes.PUT("upvote/:id", vote.AddUpVoteByBlogID)
						voteRoutes.DELETE("upvote/:id", vote.DeleteUpvoteByBlogID)

						voteRoutes.PUT("downvote/:id", vote.AddDownVoteByBlogID)
						voteRoutes.DELETE("downvote/:id", vote.DeleteDownvoteByBlogID)
					}
				}
				blogRoutes.GET("/:id", blog.GetBlogByID)
			}

			commentRoutes := forumRoutes.Group("/comment")
			{
				commentRoutes.Use(middlewares.AuthMiddleware())
				commentRoutes.PUT("/", comment.CreateComment)
				commentRoutes.DELETE("/:id", comment.DeleteCommentById)
				commentRoutes.POST("/:id", comment.UpdateCommentById)
			}

			forumRoutes.GET("/new-feed/:page", feeds.GetNewFeedPerPage)
			forumRoutes.GET("/trending/:page", feeds.GetTredingPerPage)
			forumRoutes.GET("/check", public.CheckHealth)
		}
	}
	return r
}
