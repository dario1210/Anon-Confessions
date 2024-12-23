package posts

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandlePosts(c *gin.Context) {
	post := Post{
		ID:      1,
		Title:   "Hello, World!",
		Content: "This is a test post.dsadsa",
	}

	c.JSON(http.StatusOK, post)
}
