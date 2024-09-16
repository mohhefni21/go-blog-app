package entity

import (
	"database/sql"
	"mohhefni/go-blog-app/apps/post/request"
	"time"
)

type status string

var (
	POST_DRAFT   status = "draft"
	POST_PUBLISH status = "publish"
)

type PostEntity struct {
	PostId      int            `db:"post_id"`
	UserId      int            `db:"user_id"`
	Cover       sql.NullString `db:"cover"`
	Title       string         `db:"title"`
	Slug        string         `db:"slug"`
	Excerpt     string         `db:"excerpt"`
	Content     string         `db:"content"`
	Status      string         `db:"status"`
	PublishedAt time.Time      `db:"published_at"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
}

type PostsPaginationEntity struct {
	Cursor int
	Limit  int
	Search string
}

// List posts
type GetListPostsEntity struct {
	PostId      int            `db:"post_id"`
	Cover       sql.NullString `db:"cover"`
	Title       string         `db:"title"`
	Slug        string         `db:"slug"`
	Excerpt     string         `db:"excerpt"`
	PublishedAt sql.NullTime   `db:"published_at"`
	Username    string         `db:"username"`
	Fullname    string         `db:"fullname"`
	Picture     sql.NullString `db:"picture"`
}

// Detail Posts
type GetDetailPostResponseEntity struct {
	PostId      int                               `db:"post_id"`
	Cover       sql.NullString                    `db:"cover"`
	Title       string                            `db:"title"`
	Content     string                            `db:"content"`
	PublishedAt time.Time                         `db:"published_at"`
	Author      GetDetailPostAuthorResponseEntity `db:"author"`
}

type ContentImage struct {
	IdPost    int       `db:"id_post"`
	FileName  string    `db:"filename"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type GetDetailPostAuthorResponseEntity struct {
	Username string         `db:"username"`
	Fullname string         `db:"fullname"`
	Picture  sql.NullString `db:"picture"`
}

type GetListPostsByUserLoginEntity struct {
	PostId      int            `db:"post_id"`
	Cover       sql.NullString `db:"cover"`
	Title       string         `db:"title"`
	Slug        string         `db:"slug"`
	Status      string         `db:"status"`
	PublishedAt sql.NullTime   `db:"published_at"`
	CreatedAt   time.Time      `db:"created_at"`
}

func NewFromRequestAddPostRequest(req request.AddPostRequestPayload) PostEntity {
	return PostEntity{
		Title:     req.Title,
		Excerpt:   req.Excerpt,
		Content:   req.Content,
		Status:    req.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func NewFromRequest(req request.GetPostsRequestPayload) PostsPaginationEntity {
	req.DefaultValuePagination()
	return PostsPaginationEntity{
		Cursor: req.Cursor,
		Limit:  req.Limit,
		Search: req.Search,
	}
}

func NewFromRequestUpdatePostRequest(req request.UpdatePostRequestPayload) PostEntity {
	return PostEntity{
		Title:     req.Title,
		Excerpt:   req.Excerpt,
		Content:   req.Content,
		Status:    req.Status,
		UpdatedAt: time.Now(),
	}
}

func NewFromUploadContentImageRequest(idPost int, fileName string) ContentImage {
	return ContentImage{
		IdPost:    idPost,
		FileName:  fileName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (p *PostEntity) StrToTimestamp(stringTime string) (timeParse time.Time, err error) {
	timeParse, err = time.Parse("2006-01-02 15:04:05", stringTime)
	if err != nil {
		return
	}

	return
}

type Comment struct {
	CommentId int       `db:"comment_id"`
	PostId    int       `db:"post_id"`
	UserId    int       `db:"user_id"`
	ParentId  *int      `db:"parent_id,omitempty"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Level     int       `db:"-"`
	Replies   []Comment `db:"replies,omitempty"`
}
