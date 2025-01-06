package models

import "time"

// PostDBModel is used by GORM to represent a post in the database.
type PostDBModel struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Content    string    `json:"content" gorm:"type:text;not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UserId     int       `json:"user_id" gorm:"not null"`
	TotalLikes int       `json:"total_likes" gorm:"default:0"`
}

// GetPostWithComments represents a post along with its associated comments.
// Used in API responses to fetch posts and their related comments.
type GetPostWithComments struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"createdAt"`
	TotalLikes int       `json:"totalLikes"`
	UserId     int       `json:"userId"`
	Comments   []Comment `json:"comments" gorm:"foreignKey:PostID;references:ID"`
}

// PostRequest is used for creating or updating a post.
// This is validated in POST or PATCH requests to ensure valid content.
type PostRequest struct {
	Content string `json:"content" binding:"required,min=2"`
}

// GetPost represents a minimal view of a post with metadata and user interaction details.
type GetPost struct {
	ID         int       `json:"id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"createdAt"`
	TotalLikes int       `json:"totalLikes"`
	IsLiked    int       `json:"isLiked"`
}

// GetPostsCollection is a slice of GetPost, used for paginated responses or post collections.
type GetPostsCollection []GetPost

// PostQueryParams defines query parameters for fetching posts.
// Includes pagination, sorting, and filtering options.
type PostQueryParams struct {
	Page               int    `form:"page" binding:"omitempty,min=1"`
	Limit              int    `form:"limit" binding:"omitempty,min=1"`
	SortByCreationDate string `form:"creation_date" binding:"omitempty,oneof=asc desc"`
	SortByLikes        string `form:"sort_by_likes" binding:"omitempty,oneof=asc desc"`
}

// UpdateLikesRequest is used for updating likes on a post.
type UpdateLikesRequest struct {
	Action string `json:"action" binding:"required,oneof=Like Unlike"`
}

// PostsLikesDBModel represents a many-to-many relationship between users and liked posts.
// Used by GORM for likes functionality.
type PostsLikesDBModel struct {
	PostId int `json:"post_id"`
	UserId int `json:"user_id"`
}

// TableName overrides the default table name for GORM for various models.
func (PostDBModel) TableName() string         { return "posts" }
func (GetPost) TableName() string             { return "posts" }
func (GetPostWithComments) TableName() string { return "posts" }
func (GetPostsCollection) TableName() string  { return "posts" }
func (PostsLikesDBModel) TableName() string   { return "posts_likes" }
