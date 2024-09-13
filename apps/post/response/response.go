package response

import (
	"mohhefni/go-blog-app/apps/post/entity"
	"time"
)

type GetListPostsResponse struct {
	PostId      int                        `json:"post_id"`
	Cover       string                     `json:"cover"`
	Title       string                     `json:"title"`
	Slug        string                     `json:"slug"`
	Excerpt     string                     `json:"excerpt"`
	PublishedAt time.Time                  `json:"published_at"`
	Author      GetListPostsAuthorResponse `json:"author"`
}

type GetListPostsAuthorResponse struct {
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Picture  string `json:"picture"`
}

func NewListPostsResponse(posts []entity.GetListPostsEntity) (postsList []GetListPostsResponse) {
	postsList = []GetListPostsResponse{}

	for _, post := range posts {
		postsList = append(postsList, GetListPostsResponse{
			PostId:      post.PostId,
			Cover:       post.Cover.String,
			Title:       post.Title,
			Slug:        post.Slug,
			Excerpt:     post.Excerpt,
			PublishedAt: post.PublishedAt.Time,
			Author: GetListPostsAuthorResponse{
				Username: post.Username,
				Fullname: post.Fullname,
				Picture:  post.Picture.String,
			},
		})
	}

	return postsList
}

type GetDetailPostResponse struct {
	PostId      int                         `json:"post_id"`
	Cover       string                      `json:"cover"`
	Title       string                      `json:"title"`
	Content     string                      `json:"content"`
	PublishedAt time.Time                   `json:"published_at"`
	Author      GetDetailPostAuthorResponse `json:"author"`
}

type GetDetailPostAuthorResponse struct {
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Picture  string `json:"picture"`
}

func NewDetailPostResponse(post entity.GetDetailPostResponseEntity) (detailPost GetDetailPostResponse) {
	return GetDetailPostResponse{
		PostId:      post.PostId,
		Cover:       post.Cover.String,
		Title:       post.Title,
		Content:     post.Content,
		PublishedAt: post.PublishedAt,
		Author: GetDetailPostAuthorResponse{
			Username: post.Author.Username,
			Fullname: post.Author.Fullname,
			Picture:  post.Author.Picture.String,
		},
	}
}

type GetListPostsByUserLoginResponse struct {
	PostId      int       `db:"post_id"`
	Cover       string    `db:"cover"`
	Title       string    `db:"title"`
	Slug        string    `db:"slug"`
	Status      string    `db:"status"`
	PublishedAt time.Time `db:"published_at"`
	CreatedAt   time.Time `db:"created_at"`
}

func NewListPostsByUserLoginResponse(posts []entity.GetListPostsByUserLoginEntity) (postsList []GetListPostsByUserLoginResponse) {
	postsList = []GetListPostsByUserLoginResponse{}

	for _, post := range posts {
		postsList = append(postsList, GetListPostsByUserLoginResponse{
			PostId:      post.PostId,
			Cover:       post.Cover.String,
			Title:       post.Title,
			Slug:        post.Slug,
			Status:      post.Status,
			PublishedAt: post.PublishedAt.Time,
			CreatedAt:   post.CreatedAt,
		})
	}

	return postsList
}
