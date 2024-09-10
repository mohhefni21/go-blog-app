package response

import (
	"fmt"
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
		if post.Picture.String != "" {
			post.Picture.String = fmt.Sprintf("/api/v1/auth/profile/%s", post.Picture.String)
		}
		if post.Cover.String != "" {
			post.Cover.String = fmt.Sprintf("/api/v1/posts/cover/%s", post.Cover.String)
		}

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
