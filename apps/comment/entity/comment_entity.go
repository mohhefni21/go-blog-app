package entity

import (
	"database/sql"
	"mohhefni/go-blog-app/apps/comment/request"
	"time"
)

type CommentEntity struct {
	CommentId int           `db:"comment_id"`
	PostId    int           `db:"post_id"`
	UserId    int           `db:"user_id"`
	ParentId  sql.NullInt64 `db:"parent_id"`
	Content   string        `db:"content"`
	CreatedAt time.Time     `db:"created_at"`
	UpdatedAt time.Time     `db:"updated_at"`
}

func NewFromAddCommentRequest(req request.AddCommentPayload) CommentEntity {
	return CommentEntity{
		PostId:    req.PostId,
		Content:   req.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func NewFromUpdateCommentRequest(req request.UpdateCommentPayload, commentId int) CommentEntity {
	return CommentEntity{
		CommentId: commentId,
		Content:   req.Content,
		UpdatedAt: time.Now(),
	}
}
