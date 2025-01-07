package comments_test

import (
	"anon-confessions/cmd/internal/helper"
	"anon-confessions/cmd/internal/helper/testutils"
	"anon-confessions/cmd/internal/models"
	"anon-confessions/cmd/internal/modules/comments"
	"anon-confessions/cmd/internal/modules/posts"
	"anon-confessions/cmd/internal/websocket"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupCommentsTest() *gin.Engine {
	gin.SetMode(gin.TestMode)

	mockAuthMiddleware := func(c *gin.Context) {
		c.Set("userID", 1)
		c.Next()
	}

	db := testutils.SetupMockDB()

	hub := websocket.NewHub()
	go hub.Run()

	// Initialize repositories, services, and handlers
	postsRepo := posts.NewSQLitePostsRepository(db)
	commentsRepo := comments.NewSQLiteCommentsRepository(db)
	postsService := posts.NewPostsService(postsRepo, hub)
	commentsService := comments.NewCommentsService(commentsRepo, hub)

	handler := comments.NewCommentsHandler(commentsService, postsService)

	// Set up router and register routes
	router := gin.Default()
	apiGroup := router.Group("/api/v1")
	authenticated := apiGroup.Group("/")
	authenticated.Use(mockAuthMiddleware)
	comments.RegisterCommentsRoutes(authenticated, handler)

	return router
}

func TestCreateCommentsHandler(t *testing.T) {
	router := setupCommentsTest()

	// Seed a post to create a comment for
	testutils.SeedPost()

	// Prepare request body
	reqBody := models.CreateCommentRequest{
		Content: "This is a test comment",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	// Make request
	w, req := testutils.HTTPTestRequest(http.MethodPost, "/api/v1/posts/1/comments", reqBodyBytes)
	router.ServeHTTP(w, req)

	t.Logf("Full response: %s", w.Body.String())

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var resp helper.SuccessMessage
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Message != "Comment Created Successfully" {
		t.Errorf("Expected message 'Post Created Successfully', got '%s'", resp.Message)
	}
}

func TestGetCommentsCollection(t *testing.T) {
	router := setupCommentsTest()

	// Seed a post and comments
	testutils.SeedPost()
	testutils.SeedComment(2, "First comment")
	testutils.SeedComment(2, "Second comment")

	// Make request
	w, req := testutils.HTTPTestRequest(http.MethodGet, "/api/v1/posts/2/comments", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var comments []models.Comments
	if err := json.Unmarshal(w.Body.Bytes(), &comments); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	t.Logf("Comments: %v", comments)

	if len(comments) != 2 {
		t.Errorf("Expected 2 comments, got %d", len(comments))
	}
}

func TestUpdateCommentHandler(t *testing.T) {
	router := setupCommentsTest()

	// Seed a post and comment
	testutils.SeedPost()
	testutils.SeedComment(1, "Original comment")

	// Prepare request body
	reqBody := models.CreateCommentRequest{
		Content: "Updated comment",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	// Make request
	w, req := testutils.HTTPTestRequest(http.MethodPatch, "/api/v1/posts/1/comments/1", reqBodyBytes)
	router.ServeHTTP(w, req)

	t.Logf("Full response: %s", w.Body.String())

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var resp helper.SuccessMessage
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Message != "Comment updated successfully." {
		t.Errorf("Expected message 'Comment updated successfully.', got '%s'", resp.Message)
	}
}

func TestDeleteCommentHandler(t *testing.T) {
	router := setupCommentsTest()

	// Seed a post and comment
	testutils.SeedPost()
	testutils.SeedComment(1, "Comment to delete")

	// Make request
	w, req := testutils.HTTPTestRequest(http.MethodDelete, "/api/v1/posts/1/comments/1", nil)
	router.ServeHTTP(w, req)

	t.Logf("Full response: %s", w.Body.String())

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var resp helper.SuccessMessage
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Message != "Comment Deleted Successfully" {
		t.Errorf("Expected message 'Comment Deleted Successfully', got '%s'", resp.Message)
	}
}
