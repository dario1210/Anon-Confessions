package posts

import (
	"anon-confessions/cmd/internal/helper"
	"anon-confessions/cmd/internal/models"
	"log"
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

	// Validate Request. (gin uses a validator library)
	var post models.PostRequest
	if err := c.ShouldBindJSON(&post); err != nil {
		log.Println(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, helper.ErrorMessage{Message: "Invalid request body. Please check your input."})
		return
	}
	ctx := c.Request.Context()

	err := h.postsService.CreatePosts(ctx, post, userId)
	if err != nil {
		log.Println(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.JSON(http.StatusInternalServerError, helper.ErrorMessage{Message: "Post Creation Failed"})
		return
	}


	c.JSON(http.StatusCreated, helper.SuccessMessage{Message: "Post Created Succesfully"})
}

func (h *PostsHandler) GetPostHandler(c *gin.Context) {
	ctx := c.Request.Context()
	id := helper.ParseIDParam(c, "id")

	post, err := h.postsService.GetPost(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.ErrorMessage{Message: "Failed to retireve post."})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostsHandler) GetPostsCollectionHandler(c *gin.Context) {
	ctx := c.Request.Context()
	userId := helper.RetrieveLoggedInUserId(c)

	var postQueryParam models.PostQueryParams
	if err := c.ShouldBindQuery(&postQueryParam); err != nil {
		log.Println(http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, helper.ErrorMessage{Message: "Invalid query params. Please check your input."})
		return
	}

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
		c.JSON(http.StatusInternalServerError, helper.ErrorMessage{Message: "Failed to retireve posts."})
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
		c.JSON(http.StatusInternalServerError, helper.ErrorMessage{Message: "Failed to delete post."})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, helper.ErrorMessage{Message: "Post does not exist."})
		return
	}

	c.JSON(http.StatusOK, helper.SuccessMessage{Message: "Post deleted succesfully."})
}

func (h *PostsHandler) UpdatePostsHandler(c *gin.Context) {
	userId := helper.RetrieveLoggedInUserId(c)
	postId := helper.ParseIDParam(c, "id")

	// Validate Request. (gin uses a validator library)
	var post models.PostRequest
	if err := c.ShouldBindJSON(&post); err != nil {
		log.Println(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, helper.ErrorMessage{Message: "Invalid request body. Please check your input."})
		return
	}
	ctx := c.Request.Context()

	rowsAffected, err := h.postsService.UpdatePosts(ctx, postId, userId, post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.ErrorMessage{Message: "Failed to update post."})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, helper.ErrorMessage{Message: "Update was not succesful"})
		return
	}

	c.JSON(http.StatusOK, helper.SuccessMessage{Message: "Updated succesfully"})
}

func (h *PostsHandler) UpdateLikesHandler(c *gin.Context) {
	postId := helper.ParseIDParam(c, "id")
	userId := helper.RetrieveLoggedInUserId(c)
	ctx := c.Request.Context()

	var postsLikes models.UpdateLikesRequest
	if err := c.ShouldBindJSON(&postsLikes); err != nil {
		log.Println("Error parsing request body:", err)
		c.JSON(http.StatusBadRequest, helper.ErrorMessage{Message: "Invalid request body"})
		return
	}
	log.Println("Validation successful.")

	rowsAffected, err := h.postsService.UpdateLikes(ctx, postId, userId, postsLikes)
	if err != nil {
		log.Println("Error updating likes:", err)
		c.JSON(http.StatusInternalServerError, helper.ErrorMessage{Message: "Updating likes failed."})
		return
	}
	if rowsAffected == 0 {
		log.Println("No rows were affected; possibly redundant like/unlike operation")
		c.JSON(http.StatusBadRequest, helper.ErrorMessage{Message: "Action already performed."})
		return
	}

	c.JSON(http.StatusOK, helper.SuccessMessage{Message: "Post likes updated successfully"})
}
