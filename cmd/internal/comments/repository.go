package comments

import (
	"context"
	"log"

	"gorm.io/gorm"
)

type CommentsRepository interface {
	CreateComments(context.Context, CommentsDbModel) error
	GetCommentsCollection(context.Context, int) (*GetCommentsCollection, error)
	UpdateComments(context.Context, int, int, int, CreateCommentRequest) (int64, error)
	DeleteComments(context.Context, int, int, int) (int64, error)
}

type SQLiteCommentsRepository struct {
	db *gorm.DB
}

func NewSQLiteCommentsRepository(db *gorm.DB) *SQLiteCommentsRepository {
	return &SQLiteCommentsRepository{db: db}
}

func (repo *SQLiteCommentsRepository) CreateComments(ctx context.Context, commentsDbModel CommentsDbModel) error {
	result := repo.db.WithContext(ctx).Create(&commentsDbModel)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *SQLiteCommentsRepository) GetCommentsCollection(ctx context.Context, postId int) (*GetCommentsCollection, error) {
	var commentsCollection GetCommentsCollection

	result := repo.db.WithContext(ctx).Where("post_id = ?", postId).Find(&commentsCollection)

	if result.Error != nil {
		log.Println(result.Error)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &commentsCollection, nil
}

func (repo *SQLiteCommentsRepository) UpdateComments(ctx context.Context, commentId, postId, userId int, comment CreateCommentRequest) (int64, error) {
	result := repo.db.WithContext(ctx).Model(&CommentsDbModel{}).Where("id = ? AND post_id = ? AND user_id = ?", commentId, postId, userId).Update("content", comment.Content)

	if result.Error != nil {
		return -1, result.Error
	}

	return result.RowsAffected, nil
}

func (repo *SQLiteCommentsRepository) DeleteComments(ctx context.Context, postId, userId, commentId int) (int64, error) {
	result := repo.db.Where("post_id = ? AND user_id = ? AND id = ?", postId, userId, commentId).Delete(&CommentsDbModel{})

	if result.Error != nil {
		return -1, result.Error
	}

	return result.RowsAffected, nil
}
