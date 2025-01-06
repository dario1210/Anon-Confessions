package posts

import (
	"anon-confessions/cmd/internal/models"
	"anon-confessions/cmd/internal/websocket"
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	postDBModel := models.PostDBModel{
		Content:   post.Content,
		CreatedAt: time.Now(),
		UserId:    userID,
	}

	err := s.PostsRepo.CreatePosts(ctx, postDBModel)
	if err != nil {
		log.Println("Create Posts Service:", err)
		return err
	}

	wsMsg := models.WebSocketMessage{
		Type:    "newPost",
		Message: "New Post was created",
	}

	marshalledWSMsg, err := json.Marshal(wsMsg)
	if err != nil {
		slog.Warn("Failed to marshal websocket message" + string(marshalledWSMsg))
	}
	s.hub.Broadcast <- marshalledWSMsg

	return nil
}

func (s *PostsService) GetPost(ctx context.Context, postID int) (*models.GetPostWithComments, error) {
	post, err := s.PostsRepo.GetPost(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve post: %w", err)
	}

	return post, nil
}

func (s *PostsService) GetPostsCollection(ctx context.Context, userId int, postQueryParam models.PostQueryParams) (*models.GetPostsCollection, error) {
	postCollection, err := s.PostsRepo.GetPostsCollection(ctx, userId, postQueryParam)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve posts: %w", err)
	}

	return postCollection, nil
}

func (s *PostsService) DeletePost(id int, userId int) (int64, error) {
	rowsAffected, err := s.PostsRepo.DeletePost(id, userId)
	if err != nil {
		return -1, fmt.Errorf("failed to delete post: %w", err)
	}

	return rowsAffected, nil
}

func (s *PostsService) UpdatePosts(ctx context.Context, postId, userId int, post models.PostRequest) (int64, error) {
	rowsAffected, err := s.PostsRepo.UpdatePosts(ctx, postId, userId, post)
	if err != nil {
		return -1, fmt.Errorf("failed to update post: %w", err)
	}

	return rowsAffected, nil
}

func (s *PostsService) UpdateLikes(ctx context.Context, postId, userId int, postsLikes models.UpdateLikesRequest) (int64, error) {
	var value int

	switch postsLikes.Action {
	case "Like":
		value = +1
	case "Unlike":
		value = -1
	default:
		return -1, fmt.Errorf("invalid action: %s", postsLikes.Action)
	}

	rowsAffected, err := s.PostsRepo.UpdateLikes(ctx, postId, userId, value)
	if err != nil {
		return -1, fmt.Errorf("failed to update likes: %w", err)
	}

	// Broadcast to all connected clients
	wsMsg := models.WebSocketMessage{
		Type:    "updatedLikes",
		Message: "Likes Updated",
		Content: map[string]interface{}{
			"postId": postId,
		},
	}

	marshalledWSMsg, err := json.Marshal(wsMsg)
	if err != nil {
		slog.Warn("Failed to marshal websocket message" + string(marshalledWSMsg))
	}
	s.hub.Broadcast <- marshalledWSMsg

	return rowsAffected, nil
}
