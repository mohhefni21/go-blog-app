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
	Comment     []CommentResponse           `json:"comment"`
}

type GetDetailPostAuthorResponse struct {
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Picture  string `json:"picture"`
}

type CommentResponse struct {
	CommentId int               `json:"comment_id"`
	PostId    int               `json:"post_id"`
	UserId    int               `json:"user_id"`
	ParentId  *int              `json:"parent_id,omitempty"`
	Content   string            `json:"content"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Level     int               `json:"-"`
	Replies   []CommentResponse `json:"replies,omitempty"`
}

func ConvertToCommentResponse(comments []entity.Comment) []CommentResponse {
	commentMap := make(map[int]*CommentResponse)

	// Masukkan semua komentar ke dalam map
	for _, comment := range comments {
		fmt.Printf("Processing comment ID: %d, ParentId: %v\n", comment.CommentId, comment.ParentId)
		commentMap[comment.CommentId] = &CommentResponse{
			CommentId: comment.CommentId,
			PostId:    comment.PostId,
			UserId:    comment.UserId,
			ParentId:  comment.ParentId,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
			Level:     comment.Level,
		}
	}

	// Proses nested map
	for _, comment := range comments {
		// Mengecek apakah ParentId tidak bernilai nol, artinya ini adalah reply dari komentar lain
		if comment.ParentId != nil && *comment.ParentId != 0 {
			// Jika iya, maka ambil komentar induk (parent comment) berdasarkan ParentId dari reply comment tersebut
			parentComment, exists := commentMap[*comment.ParentId]
			// Jika parent comment ditemukan
			if exists {
				// Tambahkan reply comment ke dalam array Replies dari parent comment
				parentComment.Replies = append(parentComment.Replies, *commentMap[comment.CommentId])
			}
		}
	}

	// Ubah map menjadi slice untuk dikembalikan
	var result []CommentResponse
	for _, comment := range commentMap {
		if comment.ParentId == nil || *comment.ParentId == 0 {
			fmt.Printf("Appending root comment ID: %d\n", comment.CommentId)
			result = append(result, *comment)
		}
	}

	fmt.Printf("Final result size: %d\n", len(result))
	return result
}

func NewDetailPostResponse(posts entity.GetDetailPostResponseEntity, comments []entity.Comment) (detailPost GetDetailPostResponse) {
	commentResponses := ConvertToCommentResponse(comments)

	fmt.Print(commentResponses)

	return GetDetailPostResponse{
		PostId:      posts.PostId,
		Cover:       posts.Cover.String,
		Title:       posts.Title,
		Content:     posts.Content,
		PublishedAt: posts.PublishedAt,
		Author: GetDetailPostAuthorResponse{
			Username: posts.Author.Username,
			Fullname: posts.Author.Fullname,
			Picture:  posts.Author.Picture.String,
		},
		Comment: commentResponses,
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
