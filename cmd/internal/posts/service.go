package posts

import (
	"anon-confessions/cmd/internal/models"
	"anon-confessions/cmd/internal/websocket"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

type PostsService struct {
	PostsRepo PostsRepository
	hub       *websocket.Hub
}

func NewPostsService(PostsRepo PostsRepository, hub *websocket.Hub) *PostsService {
	return &PostsService{PostsRepo: PostsRepo, hub: hub}
}

func (s *PostsService) CreatePosts(ctx context.Context, post models.PostRequest, userID int) error {
	slog.Info("Creating a new post", slog.Int("userId", userID))

	postDBModel := models.PostDBModel{
		Content:   post.Content,
		CreatedAt: time.Now(),
		UserId:    userID,
	}

	err := s.PostsRepo.CreatePosts(ctx, postDBModel)
	if err != nil {
		slog.Error("Failed to create post", slog.String("error", err.Error()), slog.Int("userId", userID))
		return err
	}

	slog.Info("Post created successfully", slog.Int("userId", userID))

	wsMsg := models.WebSocketMessage{
		Type:    "newPost",
		Message: "New Post was created",
	}

	marshalledWSMsg, err := json.Marshal(wsMsg)
	if err != nil {
		slog.Warn("Failed to marshal websocket message", slog.String("error", err.Error()))
		return nil
	}
	s.hub.Broadcast <- marshalledWSMsg

	slog.Debug("Broadcasted new post message via WebSocket", slog.Int("userId", userID))
	return nil
}

func (s *PostsService) GetPost(ctx context.Context, postID int) (*models.GetPostWithComments, error) {
	slog.Info("Retrieving post", slog.Int("postId", postID))

	post, err := s.PostsRepo.GetPost(ctx, postID)
	if err != nil {
		slog.Error("Failed to retrieve post", slog.Int("postId", postID), slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to retrieve post: %w", err)
	}

	slog.Info("Post retrieved successfully", slog.Int("postId", postID))
	return post, nil
}

func (s *PostsService) GetPostsCollection(ctx context.Context, userId int, postQueryParam models.PostQueryParams) (*models.GetPostsCollection, error) {
	slog.Info("Fetching posts collection", slog.Int("userId", userId))

	postCollection, err := s.PostsRepo.GetPostsCollection(ctx, userId, postQueryParam)
	if err != nil {
		slog.Error("Failed to retrieve posts collection", slog.String("error", err.Error()), slog.Int("userId", userId))
		return nil, fmt.Errorf("failed to retrieve posts: %w", err)
	}

	slog.Info("Posts collection retrieved successfully", slog.Int("userId", userId))
	return postCollection, nil
}

func (s *PostsService) DeletePost(id int, userId int) (int64, error) {
	slog.Info("Attempting to delete post", slog.Int("postId", id), slog.Int("userId", userId))

	rowsAffected, err := s.PostsRepo.DeletePost(id, userId)
	if err != nil {
		slog.Error("Failed to delete post", slog.Int("postId", id), slog.String("error", err.Error()))
		return -1, fmt.Errorf("failed to delete post: %w", err)
	}

	if rowsAffected > 0 {
		slog.Info("Post deleted successfully", slog.Int("postId", id), slog.Int64("rowsAffected", rowsAffected))
	} else {
		slog.Warn("No post found to delete", slog.Int("postId", id))
	}

	return rowsAffected, nil
}

func (s *PostsService) UpdatePosts(ctx context.Context, postId, userId int, post models.PostRequest) (int64, error) {
	slog.Info("Attempting to update post", slog.Int("postId", postId), slog.Int("userId", userId))

	rowsAffected, err := s.PostsRepo.UpdatePosts(ctx, postId, userId, post)
	if err != nil {
		slog.Error("Failed to update post", slog.Int("postId", postId), slog.String("error", err.Error()))
		return -1, fmt.Errorf("failed to update post: %w", err)
	}

	if rowsAffected > 0 {
		slog.Info("Post updated successfully", slog.Int("postId", postId), slog.Int64("rowsAffected", rowsAffected))
	} else {
		slog.Warn("No rows updated", slog.Int("postId", postId))
	}

	return rowsAffected, nil
}

func (s *PostsService) UpdateLikes(ctx context.Context, postId, userId int, postsLikes models.UpdateLikesRequest) (int64, error) {
	slog.Info("Updating likes for post", slog.Int("postId", postId), slog.Int("userId", userId), slog.String("action", postsLikes.Action))

	var value int
	switch postsLikes.Action {
	case "Like":
		value = +1
	case "Unlike":
		value = -1
	default:
		slog.Warn("Invalid action for updating likes", slog.String("action", postsLikes.Action))
		return -1, fmt.Errorf("invalid action: %s", postsLikes.Action)
	}

	rowsAffected, err := s.PostsRepo.UpdateLikes(ctx, postId, userId, value)
	if err != nil {
		slog.Error("Failed to update likes", slog.Int("postId", postId), slog.String("error", err.Error()))
		return -1, fmt.Errorf("failed to update likes: %w", err)
	}

	// Broadcast updated likes message
	wsMsg := models.WebSocketMessage{
		Type:    "updatedLikes",
		Message: "Likes Updated",
		Content: map[string]interface{}{
			"postId": postId,
		},
	}

	marshalledWSMsg, err := json.Marshal(wsMsg)
	if err != nil {
		slog.Warn("Failed to marshal websocket message", slog.String("error", err.Error()))
		return rowsAffected, nil 
	}
	s.hub.Broadcast <- marshalledWSMsg

	return rowsAffected, nil
}
