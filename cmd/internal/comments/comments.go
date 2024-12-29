package comments

import "time"

type CommentsDbModel struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UserId    int       `json:"user_id"`
	PostId    int       `json:"post_id"`
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=2"`
}

type Comments struct {
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

type GetCommentsCollection []Comments

func (CommentsDbModel) TableName() string {
	return "comments"
}
