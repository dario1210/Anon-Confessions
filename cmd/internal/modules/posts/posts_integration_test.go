// Package posts_test contains integration tests for the Posts functionality in the application.
// These tests validate the behavior of various endpoints related to posts, including creating,
// retrieving, updating, and deleting posts. The tests ensure that all layers (handler → service → repository)
// work together seamlessly and mimic real-world conditions without relying on mocks for intermediary components.

package posts_test

import (
	"anon-confessions/cmd/internal/helper/testutils"
	"anon-confessions/cmd/internal/models"
	"anon-confessions/cmd/internal/modules/posts"
	"anon-confessions/cmd/internal/websocket"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

// setupPostsTest initializes the test environment for posts-related endpoints, including
// setting up the router, mock database, and required middleware.
func setupPostsTest() *gin.Engine {
	gin.SetMode(gin.TestMode)

	// Mock authentication middleware
	mockAuthMiddleware := func(c *gin.Context) {
		c.Set("userID", 1)
		c.Next()
	}

	// Set up mock database
	db := testutils.SetupMockDB()

	hub := websocket.NewHub()
	go hub.Run()

	// Initialize repository, service, and handler
	repo := posts.NewSQLitePostsRepository(db)
	service := posts.NewPostsService(repo, hub)
	handler := posts.NewPostsHandler(service)

	// Set up router
	router := gin.Default()
	apiGroup := router.Group("/api/v1")
	authenticated := apiGroup.Group("/")
	authenticated.Use(mockAuthMiddleware)
	posts.RegisterPostRoutes(authenticated, handler)

	return router
}

// TestCreatePostHandler tests if a post is created successfully.
func TestCreatePostHandler(t *testing.T) {
	router := setupPostsTest()

	// Prepare the request body
	reqBody := models.PostRequest{
		Content: "This is a test post.",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	w, req := testutils.HTTPTestRequest(http.MethodPost, "/api/v1/posts/", reqBodyBytes)
	router.ServeHTTP(w, req)

	// Assert status code
	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Parse the response
	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Assert the "msg" field in the response
	if resp["msg"] != "Post Created Successfully" {
		t.Errorf("Expected msg 'Post Created Successfully', got '%s'", resp["msg"])
	}
}

// TestGetPostHandler tests if a single post can be retrieved by its ID.
func TestGetPostHandler(t *testing.T) {
	router := setupPostsTest()

	w, req := testutils.HTTPTestRequest(http.MethodGet, "/api/v1/posts/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var post models.GetPost
	if err := json.Unmarshal(w.Body.Bytes(), &post); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if post.ID != 1 {
		t.Errorf("Expected post ID 1, got %d", post.ID)
	}
}

// TestGetPostsCollectionHandler tests if a collection of posts can be retrieved.
func TestGetPostsCollectionHandler(t *testing.T) {
	router := setupPostsTest()

	w, req := testutils.HTTPTestRequest(http.MethodGet, "/api/v1/posts/?page=1&limit=10", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var posts models.GetPostsCollection
	if err := json.Unmarshal(w.Body.Bytes(), &posts); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(posts) == 0 {
		t.Errorf("Expected non-empty post collection, got %d posts", len(posts))
	}
}

// TestUpdatePostsHandler tests if a post can be updated successfully.
func TestUpdatePostsHandler(t *testing.T) {
	router := setupPostsTest()

	// Prepare the request body
	reqBody := models.PostRequest{
		Content: "Updated content for the test post.",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	w, req := testutils.HTTPTestRequest(http.MethodPatch, "/api/v1/posts/1", reqBodyBytes)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp["msg"] != "Updated successfully" {
		t.Errorf("Expected message 'Updated successfully', got '%s'", resp["message"])
	}
}

// TestUpdateLikesHandler tests if likes on a post can be updated successfully.
func TestUpdateLikesHandler(t *testing.T) {
	router := setupPostsTest()

	// Prepare the request body
	reqBody := models.UpdateLikesRequest{
		Action: "Like",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	w, req := testutils.HTTPTestRequest(http.MethodPatch, "/api/v1/posts/1/likes", reqBodyBytes)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp["msg"] != "Post likes updated successfully" {
		t.Errorf("Expected message 'Post likes updated successfully', got '%s'", resp["message"])
	}
}

// TestDeletePostsHandler tests if a post can be deleted successfully.
func TestDeletePostsHandler(t *testing.T) {
	router := setupPostsTest()

	w, req := testutils.HTTPTestRequest(http.MethodDelete, "/api/v1/posts/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp["msg"] != "Post deleted successfully." {
		t.Errorf("Expected message 'Post deleted successfully.', got '%s'", resp["message"])
	}
}
