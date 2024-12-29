package models

import "time"

type CommentsDbModel struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UserId    int       `json:"user_id" gorm:"not null"`
	PostId    int       `json:"post_id" gorm:"not null"`
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=2"`
}

type Comments struct {
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

type Comment struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Content   string    `json:"content"`
	PostID    int       `json:"postId" gorm:"column:post_id;not null"`
	CreatedAt time.Time `json:"createdAt"`
}

type GetCommentsCollection []Comments

func (CommentsDbModel) TableName() string {
	return "comments"
}
