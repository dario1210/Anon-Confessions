package posts

import (
	"anon-confessions/cmd/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterPostRoutes registers all routes related to posts.
func RegisterPostRoutes(router *gin.RouterGroup, db *gorm.DB) {
	postGroup := router.Group("/posts")
	{
		postGroup.GET("/", middleware.Authentication(db), HandlePosts)
		postGroup.POST("/", createPostHandler)
	}
}

// @Summary Get all posts
// @Description Retrieve a list of all posts
// @Tags posts
// @Accept json
// @Produce json
// @Success 200 {array} posts.Post
// @Router /posts [get]
func getPostsHandler(c *gin.Context) {}

func createPostHandler(c *gin.Context) {}
