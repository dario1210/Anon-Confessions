package comments

import (
	"anon-confessions/cmd/internal/models"
	"anon-confessions/cmd/internal/websocket"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

type CommentsService struct {
	CommentsRepo CommentsRepository
	hub          *websocket.Hub
}

func NewCommentsService(CommentsRepo CommentsRepository, hub *websocket.Hub) *CommentsService {
	return &CommentsService{CommentsRepo: CommentsRepo, hub: hub}
}

func (s *CommentsService) CreateComments(ctx context.Context, postId, userId int, comment models.CreateCommentRequest) error {
	slog.Debug("Creating a new comment", slog.Int("postId", postId), slog.Int("userId", userId))

	commentsDbModel := models.CommentsDbModel{
		Content:   comment.Content,
		CreatedAt: time.Now(),
		UserId:    userId,
		PostId:    postId,
	}

	err := s.CommentsRepo.CreateComments(ctx, commentsDbModel)
	if err != nil {
		slog.Error("Failed to create comment in repository", slog.String("error", err.Error()), slog.Int("postId", postId), slog.Int("userId", userId))
		return err
	}

	// Broadcast to all connected clients
	wsMsg := models.WebSocketMessage{
		Type:    "newComment",
		Message: "New comment created.",
		Content: map[string]interface{}{
			"postId": postId,
		},
	}

	marshalledWSMsg, err := json.Marshal(wsMsg)
	if err != nil {
		slog.Warn("Failed to marshal websocket message", slog.String("error", err.Error()), slog.Any("message", wsMsg))
		return nil
	}
	s.hub.Broadcast <- marshalledWSMsg

	return nil
}

func (s *CommentsService) GetCommentsCollection(ctx context.Context, postId int) (*models.GetCommentsCollection, error) {

	commentsCollection, err := s.CommentsRepo.GetCommentsCollection(ctx, postId)
	if err != nil {
		slog.Error("Failed to retrieve comments collection", slog.String("error", err.Error()), slog.Int("postId", postId))
		return nil, fmt.Errorf("failed to retrieve comments: %w", err)
	}

	return commentsCollection, nil
}

func (s *CommentsService) UpdateComments(ctx context.Context, commentId, postId, userId int, comment models.CreateCommentRequest) (int64, error) {
	slog.Debug("Updating comment", slog.Int("commentId", commentId), slog.Int("postId", postId), slog.Int("userId", userId))

	rowsAffected, err := s.CommentsRepo.UpdateComments(ctx, commentId, postId, userId, comment)
	if err != nil {
		slog.Error("Failed to update comment", slog.String("error", err.Error()), slog.Int("commentId", commentId), slog.Int("postId", postId), slog.Int("userId", userId))
		return -1, fmt.Errorf("failed to update comment: %w", err)
	}

	return rowsAffected, nil
}

func (s *CommentsService) DeleteComments(ctx context.Context, postId, userId, commentId int) (int64, error) {
	slog.Debug("Deleting comment", slog.Int("commentId", commentId), slog.Int("postId", postId), slog.Int("userId", userId))

	rowsAffected, err := s.CommentsRepo.DeleteComments(ctx, postId, userId, commentId)
	if err != nil {
		slog.Error("Failed to delete comment", slog.String("error", err.Error()), slog.Int("commentId", commentId), slog.Int("postId", postId), slog.Int("userId", userId))
		return -1, fmt.Errorf("failed to delete comment: %w", err)
	}

	return rowsAffected, nil
}
