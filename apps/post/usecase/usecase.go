package usecase

import (
	"context"
	"mime/multipart"
	"mohhefni/go-blog-app/apps/post/entity"
	"mohhefni/go-blog-app/apps/post/repository"
	"mohhefni/go-blog-app/apps/post/request"
	"mohhefni/go-blog-app/utility"
	"strconv"
)

type Usecase interface {
	CreatePost(ctx context.Context, req request.AddPostRequestPayload) (idPost int, err error)
	UploadCover(ctx context.Context, cover *multipart.FileHeader, idPost string) (err error)
	GetDataPosts(ctx context.Context, req request.GetPostsRequestPayload) (postEntity []entity.GetListPostsEntity, err error)
}

type usecase struct {
	repo repository.Repository
}

func NewUsecase(repo repository.Repository) Usecase {
	return &usecase{
		repo: repo,
	}
}

func (u *usecase) CreatePost(ctx context.Context, req request.AddPostRequestPayload) (idPost int, err error) {
	postEntity := entity.NewFromRequestAddPostRequest(req)

	err = u.repo.VerifyAvailableTitle(ctx, postEntity.Title)
	if err != nil {
		return
	}

	postEntity.Slug = utility.GenerateSlug(postEntity.Title)
	timeConvert, err := postEntity.StrToTimestamp(req.PublishedAt)
	if err != nil {
		return
	}

	postEntity.PublishedAt = timeConvert
	idPost, err = u.repo.AddPost(ctx, postEntity)
	if err != nil {
		return
	}

	return
}

func (u *usecase) UploadCover(ctx context.Context, cover *multipart.FileHeader, idPost string) (err error) {
	idPostInt, err := strconv.Atoi(idPost)
	if err != nil {
		return
	}

	var fileName string
	if cover != nil {
		fileName, err = utility.UploadFile(cover, "static/cover")
		if err != nil {
			return
		}
	}

	err = u.repo.UpdateCover(ctx, fileName, idPostInt)
	if err != nil {
		return
	}

	return
}

func (u *usecase) GetDataPosts(ctx context.Context, req request.GetPostsRequestPayload) (postEntity []entity.GetListPostsEntity, err error) {
	pagination := entity.NewFromRequest(req)

	postEntity, err = u.repo.GetDataPosts(ctx, pagination)
	if err != nil {
		return
	}

	return
}
