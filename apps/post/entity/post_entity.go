package entity

import (
	"mohhefni/go-blog-app/apps/post/request"
	"time"
)

type status string

var (
	POST_DRAFT   status = "draft"
	POST_PUBLISH status = "publish"
)

type PostEntity struct {
	PostId      int       `db:"post_id"`
	UserId      int       `db:"user_id"`
	Cover       string    `db:"cover"`
	Title       string    `db:"title"`
	Slug        string    `db:"slug"`
	Excerpt     string    `db:"excerpt"`
	Content     string    `db:"content"`
	Status      string    `db:"status"`
	PublishedAt time.Time `db:"published_at"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func NewFromRequestAddPostRequest(req request.AddPostRequestPayload) PostEntity {
	return PostEntity{
		UserId:    req.UserId,
		Cover:     req.Cover,
		Title:     req.Title,
		Excerpt:   req.Excerpt,
		Content:   req.Content,
		Status:    req.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (p *PostEntity) StrToTimestamp(stringTime string) (timeParse time.Time, err error) {
	timeParse, err = time.Parse("2006-01-02 15:04:05.999999", stringTime)
	if err != nil {
		return
	}

	return
}
