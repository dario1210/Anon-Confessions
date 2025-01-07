package comments

import (
	"anon-confessions/cmd/internal/models"
	"context"
	"log/slog"

	"gorm.io/gorm"
)

type CommentsRepository interface {
	CreateComments(context.Context, models.CommentsDbModel) error
	GetCommentsCollection(context.Context, int) (*models.GetCommentsCollection, error)
	UpdateComments(context.Context, int, int, int, models.CreateCommentRequest) (int64, error)
	DeleteComments(context.Context, int, int, int) (int64, error)
}

type SQLiteCommentsRepository struct {
	db *gorm.DB
}

func NewSQLiteCommentsRepository(db *gorm.DB) *SQLiteCommentsRepository {
	return &SQLiteCommentsRepository{db: db}
}

func (repo *SQLiteCommentsRepository) CreateComments(ctx context.Context, commentsDbModel models.CommentsDbModel) error {
	slog.Debug("Creating a new comment in the database", slog.Int("postId", commentsDbModel.PostId), slog.Int("userId", commentsDbModel.UserId))

	result := repo.db.WithContext(ctx).Create(&commentsDbModel)
	if result.Error != nil {
		slog.Error("Failed to create comment", slog.String("error", result.Error.Error()), slog.Int("postId", commentsDbModel.PostId), slog.Int("userId", commentsDbModel.UserId))
		return result.Error
	}

	slog.Info("Comment created successfully", slog.Int("commentId", commentsDbModel.ID))
	return nil
}

func (repo *SQLiteCommentsRepository) GetCommentsCollection(ctx context.Context, postId int) (*models.GetCommentsCollection, error) {
	slog.Debug("Retrieving comments collection for post", slog.Int("postId", postId))

	var commentsCollection models.GetCommentsCollection
	result := repo.db.WithContext(ctx).Where("post_id = ?", postId).Find(&commentsCollection)

	if result.Error != nil {
		slog.Error("Failed to retrieve comments collection", slog.String("error", result.Error.Error()), slog.Int("postId", postId))
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	slog.Info("Comments collection retrieved successfully", slog.Int("postId", postId), slog.Int64("count", result.RowsAffected))
	return &commentsCollection, nil
}

func (repo *SQLiteCommentsRepository) UpdateComments(ctx context.Context, commentId, postId, userId int, comment models.CreateCommentRequest) (int64, error) {
	slog.Debug("Updating comment in the database", slog.Int("commentId", commentId), slog.Int("postId", postId), slog.Int("userId", userId))

	result := repo.db.WithContext(ctx).Model(&models.CommentsDbModel{}).
		Where("id = ? AND post_id = ? AND user_id = ?", commentId, postId, userId).
		Update("content", comment.Content)

	if result.Error != nil {
		slog.Error("Failed to update comment", slog.String("error", result.Error.Error()), slog.Int("commentId", commentId), slog.Int("postId", postId), slog.Int("userId", userId))
		return -1, result.Error
	}

	return result.RowsAffected, nil
}

func (repo *SQLiteCommentsRepository) DeleteComments(ctx context.Context, postId, userId, commentId int) (int64, error) {
	slog.Debug("Deleting comment from the database", slog.Int("commentId", commentId), slog.Int("postId", postId), slog.Int("userId", userId))

	result := repo.db.WithContext(ctx).Where("post_id = ? AND user_id = ? AND id = ?", postId, userId, commentId).Delete(&models.CommentsDbModel{})

	if result.Error != nil {
		slog.Error("Failed to delete comment", slog.String("error", result.Error.Error()), slog.Int("commentId", commentId), slog.Int("postId", postId), slog.Int("userId", userId))
		return -1, result.Error
	}

	return result.RowsAffected, nil
}
