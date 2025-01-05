package comments

import (
	"anon-confessions/cmd/internal/models"
	"anon-confessions/cmd/internal/websocket"
	"context"
	"fmt"
	"log"
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
	commentsDbModel := models.CommentsDbModel{
		Content:   comment.Content,
		CreatedAt: time.Now(),
		UserId:    userId,
		PostId:    postId,
	}

	err := s.CommentsRepo.CreateComments(ctx, commentsDbModel)
	if err != nil {
		log.Println("Create Comments Service:", err)
		return err
	}

	s.hub.Broadcast <- []byte("New Comment Created")

	return nil
}

func (s *CommentsService) GetCommentsCollection(ctx context.Context, postId int) (*models.GetCommentsCollection, error) {
	commentsCollection, err := s.CommentsRepo.GetCommentsCollection(ctx, postId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve comments: %w", err)
	}

	return commentsCollection, nil
}

func (s *CommentsService) UpdateComments(ctx context.Context, commentId, postId, userId int, comment models.CreateCommentRequest) (int64, error) {
	rowsAffected, err := s.CommentsRepo.UpdateComments(ctx, commentId, postId, userId, comment)
	if err != nil {
		return -1, fmt.Errorf("failed to update post: %w", err)
	}

	return rowsAffected, nil

}

func (s *CommentsService) DeleteComments(ctx context.Context, postId, userId, commentId int) (int64, error) {
	rowsAffected, err := s.CommentsRepo.DeleteComments(ctx, postId, userId, commentId)
	if err != nil {
		return -1, fmt.Errorf("failed to delete post: %w", err)
	}

	return rowsAffected, nil
}
