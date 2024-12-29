package posts

import "time"

type PostDBModel struct {
	ID         int       `json:"id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UserId     int       `json:"user_id"`
	TotalLikes int       `json:"total_likes"`
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
}

type GetPostsCollection []GetPost

type UpdateLikesRequest struct {
	Action string `json:"action" binding:"required,oneof=Like Unlike"`
}

type PostsLikesDBModel struct {
	PostId int `json:"post_id"`
	UserId int `json:"user_id"`
}

// By implementing the interface we can change the name of the table.
func (PostDBModel) TableName() string {
	return "posts"
}

func (GetPost) TableName() string {
	return "posts"
}

func (GetPostsCollection) TableName() string {
	return "posts"
}

func (PostsLikesDBModel) TableName() string {
	return "posts_likes"
}
