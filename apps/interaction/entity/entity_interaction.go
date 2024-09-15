package entity

import (
	"mohhefni/go-blog-app/apps/interaction/request"
	"time"
)

type Type string

var (
	LIKE_TYPE     Type = "like"
	SHARE_TYPE    Type = "share"
	BOOKMARK_TYPE Type = "bookmark"
)

type InteractionEntity struct {
	InteractionId int       `db:"interaction_id"`
	PostId        int       `db:"post_id"`
	UserId        int       `db:"user_id"`
	Type          string    `db:"type"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func NewFromAddInteractionLikeRequest(req request.AddInteractionRequestPayload, typeInteraction string) InteractionEntity {
	return InteractionEntity{
		PostId:    req.PostId,
		Type:      typeInteraction,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
