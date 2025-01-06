package models

import "time"

// CommentsDbModel is used by GORM to represent a comment in the database.
type CommentsDbModel struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UserId    int       `json:"user_id" gorm:"not null"`
	PostId    int       `json:"post_id" gorm:"not null"`
}

// CreateCommentRequest is used to validate incoming requests for creating a comment.
// It ensures that the `content` field is present and meets the minimum length requirement.
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=2"`
}

// Comments represents a lightweight structure for comments.
// Used for collections 
type Comments struct {
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

// Comment is used for single comment responses
type Comment struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Content   string    `json:"content"`
	PostID    int       `json:"postId" gorm:"column:post_id;not null"`
	CreatedAt time.Time `json:"createdAt"`
}

type GetCommentsCollection []Comments

// TableName overrides the default table name for GORM for CommentsDbModel.
func (CommentsDbModel) TableName() string {
	return "comments"
}
