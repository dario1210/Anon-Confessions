package posts

import (
	"anon-confessions/cmd/internal/models"
	"context"
	"log"

	"gorm.io/gorm"
)

//! TODO: FIX ALL THE LOGGING AND RESPONSE HANDLING.

type PostsRepository interface {
	CreatePosts(context.Context, models.PostDBModel) error
	GetPost(context.Context, int) (*models.GetPostWithComments, error)
	GetPostsCollection(context.Context, int) (*models.GetPostsCollection, error)
	UpdatePosts(context.Context, int, int, models.PostRequest) (int64, error)
	DeletePost(int, int) (int64, error)
	UpdateLikes(context.Context, int, int, int) (int64, error)
}

type SQLitePostsRepository struct {
	db *gorm.DB
}

func NewSQLitePostsRepository(db *gorm.DB) *SQLitePostsRepository {
	return &SQLitePostsRepository{db: db}
}

func (repo *SQLitePostsRepository) CreatePosts(ctx context.Context, post models.PostDBModel) error {
	result := repo.db.WithContext(ctx).Create(&post)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *SQLitePostsRepository) GetPost(ctx context.Context, id int) (*models.GetPostWithComments, error) {
	var post models.GetPostWithComments
	err := repo.db.Preload("Comments").First(&post, id).Error
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &post, nil
}

func (repo *SQLitePostsRepository) GetPostsCollection(ctx context.Context, userId int) (*models.GetPostsCollection, error) {
	var postCollection models.GetPostsCollection

	result := repo.db.WithContext(ctx).
		Model(&models.PostDBModel{}).
		Select(`
			posts.id,
			posts.content,
			posts.created_at,
			posts.total_likes,
        	posts_likes.user_id IS NOT NULL AS IsLiked
		`).
		Joins("LEFT JOIN posts_likes ON posts.id = posts_likes.post_id AND posts_likes.user_id = ?", userId).
		Scan(&postCollection)

	if result.Error != nil {
		log.Println(result.Error)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &postCollection, nil
}

func (repo *SQLitePostsRepository) UpdatePosts(ctx context.Context, id int, userId int, post models.PostRequest) (int64, error) {
	result := repo.db.WithContext(ctx).Model(&models.PostDBModel{}).Where("id = ? AND user_id = ?", id, userId).Update("content", post.Content)

	if result.Error != nil {
		return -1, result.Error
	}

	return result.RowsAffected, nil
}

func (repo *SQLitePostsRepository) DeletePost(id, userId int) (int64, error) {
	result := repo.db.Where("id = ? AND user_id = ?", id, userId).Delete(&models.PostDBModel{})

	if result.Error != nil {
		return -1, result.Error
	}

	return result.RowsAffected, nil
}

// UpdateLikes may seem complex at first glance, but it is actually straightforward.
// Transactions are used here because the operation involves two different tables. If one operation fails,
// the entire transaction is rolled back to ensure data consistency.
// We return rowsAffected so the handler can check it and provide a better response to the API user.
func (repo *SQLitePostsRepository) UpdateLikes(ctx context.Context, postId, userId int, sign int) (int64, error) {
	var rowsAffected int64

	err := repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.PostDBModel{}).
			Where("id = ? AND user_id = ?", postId, userId).
			Update("total_likes", gorm.Expr("total_likes + ?", sign)).Error; err != nil {
			return err
		}

		// The use of raw SQL here is intentional to handle scenarios where the unique constraint on this table might be violated.
		// Instead of performing an additional query to check if the post was already liked, raw SQL with "INSERT OR IGNORE"
		// allows us to handle this case efficiently in a single operation.Ignoring the insert if it already exists.
		if sign > 0 {
			rawSQL := `
			INSERT OR IGNORE INTO posts_likes (post_id, user_id) 
			VALUES (?, ?);
			`
			result := tx.Exec(rawSQL, postId, userId)
			if result.Error != nil {
				return result.Error
			}
			rowsAffected = result.RowsAffected
		} else {
			// For the delete operation, raw SQL is not required because the `Delete` method,
			// it will return `rowsAffected` as 0 if no matching record exists,
			result := tx.Where("post_id = ? AND user_id = ?", postId, userId).Delete(&models.PostsLikesDBModel{})
			if result.Error != nil {
				return result.Error
			}
			rowsAffected = result.RowsAffected
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
