package posts

import (
	"github.com/gin-gonic/gin"
)

// RegisterPostRoutes registers all routes related to posts.
func RegisterPostRoutes(router *gin.RouterGroup) {
	postGroup := router.Group("/posts")
	{
		postGroup.GET("/", HandlePosts)
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
