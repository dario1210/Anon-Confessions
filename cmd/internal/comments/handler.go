package comments

import (
	"anon-confessions/cmd/internal/helper"
	"anon-confessions/cmd/internal/posts"
	"log"
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

	var comment CreateCommentRequest
	if err := c.ShouldBindJSON(&comment); err != nil {
		log.Println(http.StatusBadRequest, helper.ErrorMessage{Message: err.Error()})
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
		c.JSON(http.StatusInternalServerError, helper.ErrorMessage{Message: "Failed to create comment on post."})
		return
	}

	c.JSON(http.StatusCreated, helper.SuccessMessage{Message: "Post Created Succesfully"})
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
		c.JSON(http.StatusInternalServerError, helper.ErrorMessage{Message: "Failed to retrieve comments on post."})
		return
	}

	c.JSON(http.StatusOK, comments)
}

func (h *CommentsHandler) UpdateCommentHandler(c *gin.Context) {
	ctx := c.Request.Context()
	userId := helper.RetrieveLoggedInUserId(c)
	postId := helper.ParseIDParam(c, "id")
	commentId := helper.ParseIDParam(c, "commentId")

	var comment CreateCommentRequest
	if err := c.ShouldBindJSON(&comment); err != nil {
		log.Println(http.StatusBadRequest, helper.ErrorMessage{Message: err.Error()})
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
		c.JSON(http.StatusInternalServerError, helper.ErrorMessage{Message: "Failed to update comment."})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, helper.ErrorMessage{Message: "Update was not succesful"})
		return
	}

	c.JSON(http.StatusOK, helper.SuccessMessage{Message: "Comment updated succesfully."})
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
		c.JSON(http.StatusInternalServerError, helper.ErrorMessage{Message: "Failed to delete comment."})
		return
	}
	if comments == 0 {
		c.JSON(http.StatusNotFound, helper.ErrorMessage{Message: "Comment does not exist."})
		return
	}

	c.JSON(http.StatusOK, helper.SuccessMessage{Message: "Comment Deleted Sucesfully"})
}
