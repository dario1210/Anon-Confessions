package comments

import (
	"github.com/gin-gonic/gin"
)

// RegisterCommentsRoutes registers all routes related to posts.
func RegisterCommentsRoutes(router *gin.RouterGroup, h *CommentsHandler) {
	commentGroup := router.Group("/posts/:id/comments")
	{
		commentGroup.POST("", h.CreateCommentsHandler)
		commentGroup.GET("", h.GetCommentsCollection)
		commentGroup.PATCH("/:commentId", h.UpdateCommentHandler)
		commentGroup.DELETE("/:commentId", h.DeleteCommentHandler)
	}
}

// Swagger documentation.

// CreateCommentsHandler handles the creation of a comment for a specific post.
// @Summary Create a comment
// @Description Allows authenticated users to add a comment to a specific post. Requires authentication using X-Account-Number.
// @Tags comments
// @Accept json
// @Produce json
// @Param X-Account-Number header string true "16-digit account number (e.g., 1234567890123456)"
// @Param id path int true "Post ID"
// @Param comment body models.CreateCommentRequest true "Comment content"
// @Success 201 {object} helper.SuccessMessage "Comment created successfully"
// @Failure 400 {object} helper.ErrorMessage "Invalid request body"
// @Failure 401 {object} helper.ErrorMessage "Invalid or missing X-Account-Number"
// @Failure 404 {object} helper.ErrorMessage "Post not found"
// @Failure 500 {object} helper.ErrorMessage "Internal server error"
// @Router /posts/{id}/comments [post]
func (h *CommentsHandler) createCommentsHandler(c *gin.Context) {}

// GetCommentsCollection retrieves a collection of comments for a specific post.
// @Summary Retrieve comments for a post
// @Description Fetches all comments associated with a specific post ID. Requires authentication using X-Account-Number.
// @Tags comments
// @Accept json
// @Produce json
// @Param X-Account-Number header string true "16-digit account number (e.g., 1234567890123456)"
// @Param id path int true "Post ID"
// @Success 200 {object} models.GetCommentsCollection "Comments retrieved successfully"
// @Failure 500 {object} helper.ErrorMessage "Failed to retrieve comments"
// @Router /posts/{id}/comments [get]
func (h *CommentsHandler) getCommentsCollection(c *gin.Context) {}

// @Summary Update a comment
// @Description Updates the content of a specific comment in a post. Requires authentication using X-Account-Number.
// @Tags comments
// @Accept json
// @Produce json
// @Param X-Account-Number header string true "16-digit account number (e.g., 1234567890123456)"
// @Param id path int true "Post ID"
// @Param commentId path int true "Comment ID"
// @Param body body models.CreateCommentRequest true "Updated comment content"
// @Success 200 {object} helper.SuccessMessage "Comment updated successfully"
// @Failure 400 {object} helper.ErrorMessage "Invalid request body or input"
// @Failure 401 {object} helper.ErrorMessage "Invalid or missing X-Account-Number"
// @Failure 404 {object} helper.ErrorMessage "Post or comment not found"
// @Failure 500 {object} helper.ErrorMessage "Failed to update comment"
// @Router /posts/{id}/comments/{commentId} [patch]
func (h *CommentsHandler) updateCommentsHandler(c *gin.Context) {}

// @Summary      Delete a comment
// @Description  Deletes a specific comment from a post. The user must be authenticated and authorized to delete the comment.
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param X-Account-Number header string true "16-digit account number (e.g., 1234567890123456)"
// @Param        id          path      int  true  "Post ID"
// @Param        commentId   path      int  true  "Comment ID"
// @Success      200 {object} helper.SuccessMessage "Comment deleted successfully"
// @Failure      400 {object} helper.ErrorMessage   "Invalid post ID or comment ID"
// @Failure      404 {object} helper.ErrorMessage   "Post not found or comment not found"
// @Failure      500 {object} helper.ErrorMessage   "Failed to delete comment"
// @Router       /posts/{id}/comments/{commentId} [delete]
func (h *CommentsHandler) deleteComment(c *gin.Context) {}
