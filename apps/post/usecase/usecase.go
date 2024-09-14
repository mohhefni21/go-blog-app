package usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"mohhefni/go-blog-app/apps/post/entity"
	"mohhefni/go-blog-app/apps/post/repository"
	"mohhefni/go-blog-app/apps/post/request"
	"mohhefni/go-blog-app/internal/config"
	"mohhefni/go-blog-app/utility"
	"strconv"
)

type Usecase interface {
	CreatePost(ctx context.Context, req request.AddPostRequestPayload) (idPost int, err error)
	UploadCover(ctx context.Context, cover *multipart.FileHeader, idPost string) (err error)
	GetListPosts(ctx context.Context, req request.GetPostsRequestPayload) (postEntity []entity.GetListPostsEntity, err error)
	GetDetailPost(ctx context.Context, slug string) (DetailPostEntity entity.GetDetailPostResponseEntity, err error)
	GetListPostsByUsername(ctx context.Context, req request.GetPostsRequestPayload, username string) (postEntity []entity.GetListPostsEntity, err error)
	GetListPostsByUserLogin(ctx context.Context, publicId string) (post []entity.GetListPostsByUserLoginEntity, err error)
	DeletePost(ctx context.Context, slug string) (err error)
	UpdatePost(ctx context.Context, req request.UpdatePostRequestPayload, idPost string) (err error)
	UpdateImageContent(ctx context.Context, idPost string, image *multipart.FileHeader) (url string, err error)
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
	if req.PublishedAt != "" {
		timeConvert, err := postEntity.StrToTimestamp(req.PublishedAt)
		if err != nil {
			return 0, err
		}
		postEntity.PublishedAt = timeConvert
	}

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
	oldPost, err := u.repo.GetPostById(ctx, idPostInt)
	if err != nil {
		return err
	}

	var fileName string
	if cover != nil {
		if oldPost.Cover.String != "" {
			filePath := fmt.Sprintf("static/cover/%s", oldPost.Cover.String)

			err = utility.DeleteFile(filePath)
			if err != nil {
				return
			}
		}

		fileName, err = utility.UploadFile(cover, "static/cover")
		if err != nil {
			return
		}
	} else {
		fileName = oldPost.Cover.String
	}

	err = u.repo.UpdateCover(ctx, fileName, idPostInt)
	if err != nil {
		return
	}

	return
}

func (u *usecase) GetListPosts(ctx context.Context, req request.GetPostsRequestPayload) (postEntity []entity.GetListPostsEntity, err error) {
	pagination := entity.NewFromRequest(req)

	postEntity, err = u.repo.GetDataPosts(ctx, pagination)
	if err != nil {
		return
	}

	return
}

func (u *usecase) GetDetailPost(ctx context.Context, slug string) (DetailPostEntity entity.GetDetailPostResponseEntity, err error) {
	DetailPostEntity, err = u.repo.GetDetailPostBySLug(ctx, slug)
	if err != nil {
		return
	}

	return
}

func (u *usecase) GetListPostsByUsername(ctx context.Context, req request.GetPostsRequestPayload, username string) (postEntity []entity.GetListPostsEntity, err error) {
	pagination := entity.NewFromRequest(req)

	err = u.repo.VerifyAvailableUsername(ctx, username)
	if err != nil {
		return
	}

	postEntity, err = u.repo.GetDataPostsByUsername(ctx, pagination, username)
	if err != nil {
		return
	}

	return
}

func (u *usecase) GetListPostsByUserLogin(ctx context.Context, publicId string) (post []entity.GetListPostsByUserLoginEntity, err error) {
	publicIdUuid, err := utility.ParseUUID(publicId)
	if err != nil {
		return
	}

	post, err = u.repo.GetDataPostsByUserLogin(ctx, publicIdUuid)
	if err != nil {
		return
	}

	return
}

func (u *usecase) DeletePost(ctx context.Context, slug string) (err error) {

	post, err := u.repo.GetDetailPostBySLug(ctx, slug)
	if err != nil {
		return
	}

	// delete cover
	if post.Cover.String != "" {
		filePath := fmt.Sprintf("static/cover/%s", post.Cover.String)

		err = utility.DeleteFile(filePath)
		if err != nil {
			return
		}
	}

	// delete content image
	var contentImages []entity.ContentImage
	contentImages, err = u.repo.GetContentImageByPostId(ctx, post.PostId)
	if err != nil {
		return
	}

	if len(contentImages) > 0 {
		for _, contentImage := range contentImages {
			filePath := fmt.Sprintf("static/content-image/%s", contentImage.FileName)

			err = utility.DeleteFile(filePath)
			if err != nil {
				return
			}
		}
	}

	err = u.repo.DeletePostById(ctx, post.PostId)
	if err != nil {
		return
	}

	return
}

func (u *usecase) UpdatePost(ctx context.Context, req request.UpdatePostRequestPayload, idPost string) (err error) {
	idPostInt, err := strconv.Atoi(idPost)
	if err != nil {
		return
	}

	oldPost, err := u.repo.GetPostById(ctx, idPostInt)
	if err != nil {
		return
	}

	postEntity := entity.NewFromRequestUpdatePostRequest(req)
	postEntity.PostId = idPostInt

	if oldPost.Title != req.Title {
		err = u.repo.VerifyAvailableTitle(ctx, postEntity.Title)
		if err != nil {
			return
		}
		postEntity.Slug = utility.GenerateSlug(postEntity.Title)
	} else {
		postEntity.Slug = oldPost.Slug
	}

	if req.PublishedAt != "" {
		timeConvert, err := postEntity.StrToTimestamp(req.PublishedAt)
		if err != nil {
			return err
		}
		postEntity.PublishedAt = timeConvert
	}

	err = u.repo.UpdatePostById(ctx, postEntity)
	if err != nil {
		return
	}

	return
}

func (u *usecase) UpdateImageContent(ctx context.Context, idPost string, imageContent *multipart.FileHeader) (url string, err error) {
	idPostInt, err := strconv.Atoi(idPost)
	if err != nil {
		return
	}

	var fileNameSave string

	if imageContent != nil {
		fileNameSave, err = utility.UploadFile(imageContent, "static/content-image")
		if err != nil {
			return
		}
	}

	imageContentEntity := entity.NewFromUploadContentImageRequest(idPostInt, fileNameSave)

	fileName, err := u.repo.UploadImageContent(ctx, imageContentEntity)
	if err != nil {
		return
	}

	url = fmt.Sprintf("%s/api/v1/posts/content-image/%s", config.Cfg.AppConfig.BaseUrl, fileName)

	return
}
