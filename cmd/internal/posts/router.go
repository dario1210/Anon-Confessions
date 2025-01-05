package posts

import (
	"github.com/gin-gonic/gin"
)

// RegisterPostRoutes registers all routes related to posts.
func RegisterPostRoutes(router *gin.RouterGroup, postsHandler *PostsHandler) {
	postGroup := router.Group("/posts")
	{
		postGroup.POST("/", postsHandler.CreatePostHandler)
		postGroup.GET("/", postsHandler.GetPostsCollectionHandler)
		postGroup.GET("/:id", postsHandler.GetPostHandler)
		postGroup.PATCH("/:id", postsHandler.UpdatePostsHandler)
		postGroup.DELETE("/:id", postsHandler.DeletePostsHandler)
		postGroup.PATCH("/:id/likes", postsHandler.UpdateLikesHandler)

	}
}

// Swagger documentation.

// GetPost handles retrieving a post by its ID.
// @Summary Retrieve a post
// @Description Fetches a post using its unique ID. Requires authentication using X-Account-Number.
// @Tags posts
// @Accept json
// @Produce json
// @Param X-Account-Number header string true "16-digit account number (e.g., 1234567890123456)"
// @Param id path string true "Post ID"
// @Success 200 {object} models.GetPost "Post retrieved successfully"
// @Failure 400 {object} helper.ErrorMessage "Invalid post ID"
// @Failure 401 {object} helper.ErrorMessage "Invalid or missing X-Account-Number"
// @Failure 500 {object} helper.ErrorMessage "Failed to retrieve post"
// @Router /posts/{id} [get]
func getPostsHandler(c *gin.Context) {}

// CreatePostHandler handles the creation of a new post.
// @Summary Create a new post
// @Description Allows authenticated users to create a new post using their X-Account-Number.
// @Tags posts
// @Accept json
// @Produce json
// @Param X-Account-Number header string true "16-digit account number (e.g., 1234567890123456)"
// @Param post body models.PostRequest true "Post content"
// @Success 201 {object} helper.SuccessMessage "Post created successfully"
// @Failure 400 {object} helper.ErrorMessage "Invalid request body"
// @Failure 401 {object} helper.ErrorMessage "Invalid or missing X-Account-Number"
// @Failure 500 {object} helper.ErrorMessage "Internal server error"
// @Router /posts [post]
func createPostHandler(c *gin.Context) {}

// GetPostsCollectionHandler handles retrieving a collection of posts.
// @Summary Retrieve a collection of posts
// @Description Fetches a collection of posts. Requires authentication using X-Account-Number.                           If both sorting options are provided, priority will be given to the SortByCreationDate field.
// @Tags posts
// @Accept json
// @Produce json
// @Param X-Account-Number header string true "16-digit account number (e.g., 1234567890123456)"
// @Param page query int false "Page number (default: 1)" minimum(1) default(1)
// @Param limit query int false "Number of items per page (default: 10)" minimum(1) default(10)
// @Param creation_date query string false "Sort by creation date (asc or desc)" Enums(asc,desc) default()
// @Param sort_by_likes query string false "Sort by likes (asc or desc)" Enums(asc,desc) default()
// @Success 200 {object} models.GetPostsCollection "Posts retrieved successfully"
// @Success 200 {object} map[string]interface{} "{} if no posts are found"
// @Failure 500 {object} helper.ErrorMessage "Failed to retrieve posts"
// @Router /posts [get]
func (h *PostsHandler) getPostsCollectionHandler(c *gin.Context) {}

// DeletePostsHandler handles deleting a post by its ID.
// @Summary Delete a post
// @Description Deletes a post using its unique ID. Requires the user to be logged in and authenticated using X-Account-Number.
// @Tags posts
// @Accept json
// @Produce json
// @Param X-Account-Number header string true "16-digit account number (e.g., 1234567890123456)"
// @Param id path int true "Post ID"
// @Success 200 {object} helper.SuccessMessage "Post deleted successfully"
// @Failure 400 {object} helper.ErrorMessage "Invalid post ID"
// @Failure 401 {object} helper.ErrorMessage "Unauthorized user or missing X-Account-Number"
// @Failure 500 {object} helper.ErrorMessage "Failed to delete post"
// @Router /posts/{id} [delete]
func (h *PostsHandler) deletePostsHandler(c *gin.Context) {}

// UpdatePostsHandler handles updating a post by its ID.
// @Summary Update a post
// @Description Updates a post's content. Requires the user to be authenticated using X-Account-Number.
// @Tags posts
// @Accept json
// @Produce json
// @Param X-Account-Number header string true "16-digit account number (e.g., 1234567890123456)"
// @Param id path int true "Post ID"
// @Param post body models.PostRequest true "Post content"
// @Success 200 {object} helper.SuccessMessage "Updated successfully"
// @Failure 400 {object} helper.ErrorMessage "Invalid request body or parameters"
// @Failure 404 {object} helper.ErrorMessage "Post not found or no updates applied"
// @Failure 500 {object} helper.ErrorMessage "Failed to update post"
// @Router /posts/{id} [patch]
func (h *PostsHandler) updatePostsHandler(c *gin.Context) {}

// UpdateLikesHandler handles liking or unliking a post by a user.
// @Summary Like or Unlike a post
// @Description Updates the like status of a post. Requires the user to be authenticated using X-Account-Number.
// @Tags posts
// @Accept json
// @Produce json
// @Param X-Account-Number header string true "16-digit account number (e.g., 1234567890123456)"
// @Param id path int true "Post ID"
// @Param body body models.UpdateLikesRequest true "Action to like or unlike the post"
// @Success 200 {object} helper.SuccessMessage "Action applied successfully"
// @Failure 400 {object} helper.ErrorMessage "Invalid request body or parameters"
// @Failure 404 {object} helper.ErrorMessage "Post not found or action not applied"
// @Failure 500 {object} helper.ErrorMessage "Failed to apply action on the post"
// @Router /posts/{id}/likes [patch]
func (h *PostsHandler) updateLikesHandler(c *gin.Context) {}
