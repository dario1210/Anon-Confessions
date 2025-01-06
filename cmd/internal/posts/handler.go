package posts

import (
	"anon-confessions/cmd/internal/helper"
	"anon-confessions/cmd/internal/models"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PostsHandler struct {
	postsService *PostsService
}

func NewPostsHandler(postsService *PostsService) *PostsHandler {
	return &PostsHandler{postsService: postsService}
}

func (h *PostsHandler) CreatePostHandler(c *gin.Context) {
	userId := helper.RetrieveLoggedInUserId(c)

	// Validate request body.
	var post models.PostRequest
	if err := c.ShouldBindJSON(&post); err != nil {
		slog.Warn("Invalid request body for creating a post", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, helper.ErrorMessage{Message: "Invalid request body. Please check your input."})
		return
	}
	ctx := c.Request.Context()

	err := h.postsService.CreatePosts(ctx, post, userId)
	if err != nil {
		slog.Error("Failed to create post", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, helper.ErrorMessage{Message: "Post Creation Failed"})
		return
	}

	slog.Info("Post created successfully", slog.Int("userId", userId))
	c.JSON(http.StatusCreated, helper.SuccessMessage{Message: "Post Created Successfully"})
}

func (h *PostsHandler) GetPostHandler(c *gin.Context) {
	ctx := c.Request.Context()
	id := helper.ParseIDParam(c, "id")

	post, err := h.postsService.GetPost(ctx, id)
	if err != nil {
		slog.Error("Failed to retrieve post", slog.Int("postId", id), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, helper.ErrorMessage{Message: "Failed to retrieve post."})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostsHandler) GetPostsCollectionHandler(c *gin.Context) {
	ctx := c.Request.Context()
	userId := helper.RetrieveLoggedInUserId(c)

	// Parse query parameters for pagination and sorting.
	var postQueryParam models.PostQueryParams
	if err := c.ShouldBindQuery(&postQueryParam); err != nil {
		slog.Warn("Invalid query parameters for retrieving posts", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, helper.ErrorMessage{Message: "Invalid query params. Please check your input."})
		return
	}

	// Set default values if not provided.
	if postQueryParam.Page == 0 {
		postQueryParam.Page = 1
	}
	if postQueryParam.Limit == 0 {
		postQueryParam.Limit = 10
	}
	if postQueryParam.SortByLikes == "" && postQueryParam.SortByCreationDate == "" {
		postQueryParam.SortByCreationDate = "asc"
	}

	post, err := h.postsService.GetPostsCollection(ctx, userId, postQueryParam)
	if err != nil {
		slog.Error("Failed to retrieve posts", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, helper.ErrorMessage{Message: "Failed to retrieve posts."})
		return
	}

	if post == nil {
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostsHandler) DeletePostsHandler(c *gin.Context) {
	id := helper.ParseIDParam(c, "id")
	userID := helper.RetrieveLoggedInUserId(c)

	rowsAffected, err := h.postsService.DeletePost(id, userID)
	if err != nil {
		slog.Error("Failed to delete post", slog.Int("postId", id), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, helper.ErrorMessage{Message: "Failed to delete post."})
		return
	}
	if rowsAffected == 0 {
		slog.Warn("Post not found during delete operation", slog.Int("postId", id))
		c.JSON(http.StatusNotFound, helper.ErrorMessage{Message: "Post does not exist."})
		return
	}

	c.JSON(http.StatusOK, helper.SuccessMessage{Message: "Post deleted successfully."})
}

func (h *PostsHandler) UpdatePostsHandler(c *gin.Context) {
	userId := helper.RetrieveLoggedInUserId(c)
	postId := helper.ParseIDParam(c, "id")

	// Validate request body.
	var post models.PostRequest
	if err := c.ShouldBindJSON(&post); err != nil {
		slog.Warn("Invalid request body for updating post", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, helper.ErrorMessage{Message: "Invalid request body. Please check your input."})
		return
	}
	ctx := c.Request.Context()

	rowsAffected, err := h.postsService.UpdatePosts(ctx, postId, userId, post)
	if err != nil {
		slog.Error("Failed to update post", slog.Int("postId", postId), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, helper.ErrorMessage{Message: "Failed to update post."})
		return
	}
	if rowsAffected == 0 {
		slog.Warn("No rows updated during update operation", slog.Int("postId", postId))
		c.JSON(http.StatusNotFound, helper.ErrorMessage{Message: "Update was not successful"})
		return
	}

	c.JSON(http.StatusOK, helper.SuccessMessage{Message: "Updated successfully"})
}

func (h *PostsHandler) UpdateLikesHandler(c *gin.Context) {
	postId := helper.ParseIDParam(c, "id")
	userId := helper.RetrieveLoggedInUserId(c)
	ctx := c.Request.Context()

	// Parse request body for updating likes.
	var postsLikes models.UpdateLikesRequest
	if err := c.ShouldBindJSON(&postsLikes); err != nil {
		slog.Warn("Invalid request body for updating likes", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, helper.ErrorMessage{Message: "Invalid request body"})
		return
	}

	rowsAffected, err := h.postsService.UpdateLikes(ctx, postId, userId, postsLikes)
	if err != nil {
		slog.Error("Error updating likes", slog.Int("postId", postId), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, helper.ErrorMessage{Message: "Updating likes failed."})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusBadRequest, helper.ErrorMessage{Message: "Action already performed."})
		return
	}

	c.JSON(http.StatusOK, helper.SuccessMessage{Message: "Post likes updated successfully"})
}
