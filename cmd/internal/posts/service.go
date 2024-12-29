package posts

import (
	"context"
	"fmt"
	"log"
	"time"
)

type PostsService struct {
	PostsRepo PostsRepository
}

func NewPostsService(PostsRepo PostsRepository) *PostsService {
	return &PostsService{PostsRepo: PostsRepo}
}

func (s *PostsService) CreatePosts(ctx context.Context, post PostRequest, userID int) error {
	postDBModel := PostDBModel{
		Content:   post.Content,
		CreatedAt: time.Now(),
		UserId:    userID,
	}

	err := s.PostsRepo.CreatePosts(ctx, postDBModel)
	if err != nil {
		log.Println("Create Posts Service:", err)
		return err
	}
	return nil
}

func (s *PostsService) GetPost(ctx context.Context, postID int) (*GetPost, error) {
	post, err := s.PostsRepo.GetPost(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve post: %w", err)
	}

	return post, nil
}

func (s *PostsService) GetPostsCollection(ctx context.Context) (*GetPostsCollection, error) {
	postCollection, err := s.PostsRepo.GetPostsCollection(ctx)
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

func (s *PostsService) UpdatePosts(ctx context.Context, postId, userId int, post PostRequest) (int64, error) {
	rowsAffected, err := s.PostsRepo.UpdatePosts(ctx, postId, userId, post)
	if err != nil {
		return -1, fmt.Errorf("failed to update post: %w", err)
	}

	return rowsAffected, nil
}

func (s *PostsService) UpdateLikes(ctx context.Context, postId, userId int, postsLikes UpdateLikesRequest) (int64, error) {
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

	return rowsAffected, nil
}