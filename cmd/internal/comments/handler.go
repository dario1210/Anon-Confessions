package comments

import (
	"anon-confessions/cmd/internal/helper"
	"anon-confessions/cmd/internal/models"
	"anon-confessions/cmd/internal/posts"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommentsHandler struct {
	commentsService *CommentsService
	postsService    *posts.PostsService
}

func NewCommentsHandler(commentsService *CommentsService, postsService *posts.PostsService) *CommentsHandler {
	return &CommentsHandler{commentsService: commentsService, postsService: postsService}
}

func (h *CommentsHandler) CreateCommentsHandler(c *gin.Context) {
	ctx := c.Request.Context()
	userId := helper.RetrieveLoggedInUserId(c)
	postId := helper.ParseIDParam(c, "id")

	var comment models.CreateCommentRequest
	if err := c.ShouldBindJSON(&comment); err != nil {
		slog.Warn("Invalid request body", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, helper.ErrorMessage{Message: "Invalid request body. Please check your input."})
		return
	}

	_, err := h.postsService.GetPost(ctx, postId)
	if err != nil {
		c.JSON(http.StatusNotFound, helper.ErrorMessage{Message: "Post not found"})
		return
	}

	err = h.commentsService.CreateComments(ctx, postId, userId, comment)
	if err != nil {
		slog.Error("Failed to create comment", slog.String("error", err.Error()), slog.Int("postId", postId), slog.Int("userId", userId))
		c.JSON(http.StatusInternalServerError, helper.ErrorMessage{Message: "Failed to create comment on post."})
		return
	}

	c.JSON(http.StatusCreated, helper.SuccessMessage{Message: "Post Created Successfully"})
}

func (h *CommentsHandler) GetCommentsCollection(c *gin.Context) {
	ctx := c.Request.Context()
	postId := helper.ParseIDParam(c, "id")

	_, err := h.postsService.GetPost(ctx, postId)
	if err != nil {
		c.JSON(http.StatusNotFound, helper.ErrorMessage{Message: "Post not found"})
		return
	}

	comments, err := h.commentsService.GetCommentsCollection(ctx, postId)
	if err != nil {
		slog.Error("Failed to retrieve comments", slog.String("error", err.Error()), slog.Int("postId", postId))
		c.JSON(http.StatusInternalServerError, helper.ErrorMessage{Message: "Failed to retrieve comments on post."})
		return
	}

	slog.Info("Comments retrieved successfully", slog.Int("postId", postId))
	c.JSON(http.StatusOK, comments)
}

func (h *CommentsHandler) UpdateCommentHandler(c *gin.Context) {
	ctx := c.Request.Context()
	userId := helper.RetrieveLoggedInUserId(c)
	postId := helper.ParseIDParam(c, "id")
	commentId := helper.ParseIDParam(c, "commentId")

	var comment models.CreateCommentRequest
	if err := c.ShouldBindJSON(&comment); err != nil {
		slog.Warn("Invalid request body", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, helper.ErrorMessage{Message: "Invalid request body. Please check your input."})
		return
	}

	_, err := h.postsService.GetPost(ctx, postId)
	if err != nil {
		c.JSON(http.StatusNotFound, helper.ErrorMessage{Message: "Post not found"})
		return
	}

	rowsAffected, err := h.commentsService.UpdateComments(ctx, commentId, postId, userId, comment)
	if err != nil {
		slog.Error("Failed to update comment", slog.String("error", err.Error()), slog.Int("commentId", commentId))
		c.JSON(http.StatusInternalServerError, helper.ErrorMessage{Message: "Failed to update comment."})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, helper.ErrorMessage{Message: "Update was not successful"})
		return
	}

	c.JSON(http.StatusOK, helper.SuccessMessage{Message: "Comment updated successfully."})
}

func (h *CommentsHandler) DeleteCommentHandler(c *gin.Context) {
	ctx := c.Request.Context()
	userId := helper.RetrieveLoggedInUserId(c)
	postId := helper.ParseIDParam(c, "id")
	commentId := helper.ParseIDParam(c, "commentId")

	_, err := h.postsService.GetPost(ctx, postId)
	if err != nil {
		c.JSON(http.StatusNotFound, helper.ErrorMessage{Message: "Post not found"})
		return
	}

	comments, err := h.commentsService.DeleteComments(ctx, postId, userId, commentId)
	if err != nil {
		slog.Error("Failed to delete comment", slog.String("error", err.Error()), slog.Int("commentId", commentId))
		c.JSON(http.StatusInternalServerError, helper.ErrorMessage{Message: "Failed to delete comment."})
		return
	}
	if comments == 0 {
		c.JSON(http.StatusNotFound, helper.ErrorMessage{Message: "Comment does not exist."})
		return
	}

	c.JSON(http.StatusOK, helper.SuccessMessage{Message: "Comment Deleted Successfully"})
}
