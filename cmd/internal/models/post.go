package models

import (
	"time"
)

type PostDBModel struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Content    string    `json:"content" gorm:"type:text;not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UserId     int       `json:"user_id" gorm:"not null"`
	TotalLikes int       `json:"total_likes" gorm:"default:0"`
}

type GetPostWithComments struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"createdAt"`
	TotalLikes int       `json:"totalLikes"`
	UserId     int       `json:"userId"`
	Comments   []Comment `json:"comments" gorm:"foreignKey:PostID;references:ID"`
}

// PostRequest is used for creating or updating a post.
// Needed in both in POST & PATCH requests.
type PostRequest struct {
	Content string `json:"content" binding:"required,min=2"`
}

type GetPost struct {
	ID         int       `json:"id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"createdAt"`
	TotalLikes int       `json:"totalLikes"`
	IsLiked    int       `json:"isLiked"`
}

type GetPostsCollection []GetPost

type UpdateLikesRequest struct {
	Action string `json:"action" binding:"required,oneof=Like Unlike"`
}

type PostsLikesDBModel struct {
	PostId int `json:"post_id"`
	UserId int `json:"user_id"`
}

// By implementing the interface we can override the name of the table, that GORM uses when querying the database.
func (PostDBModel) TableName() string {
	return "posts"
}

func (GetPost) TableName() string {
	return "posts"
}

func (GetPostWithComments) TableName() string {
	return "posts"
}

func (GetPostsCollection) TableName() string {
	return "posts"
}

func (PostsLikesDBModel) TableName() string {
	return "posts_likes"
}
