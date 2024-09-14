package usecase

import (
	"context"
	"mohhefni/go-blog-app/apps/comment/entity"
	"mohhefni/go-blog-app/apps/comment/repository"
	"mohhefni/go-blog-app/apps/comment/request"
	"mohhefni/go-blog-app/utility"
	"strconv"
)

type Usecase interface {
	CreateComment(ctx context.Context, req request.AddCommentPayload, publicId string) (idComment int, err error)
	UpdateComment(ctx context.Context, req request.UpdateCommentPayload, commentId string) (err error)
	DeleteComment(ctx context.Context, commentId string) (err error)
}

type usecase struct {
	repo repository.Repository
}

func NewUsecase(repo repository.Repository) Usecase {
	return &usecase{
		repo: repo,
	}
}

func (u *usecase) CreateComment(ctx context.Context, req request.AddCommentPayload, publicId string) (idComment int, err error) {
	commentEntity := entity.NewFromAddCommentRequest(req)

	publicIdUuid, err := utility.ParseUUID(publicId)
	if err != nil {
		return
	}

	var userEntity entity.UserEntity
	userEntity, err = u.repo.GetUserByPublicId(ctx, publicIdUuid)
	if err != nil {
		return
	}
	commentEntity.UserId = userEntity.UserId
	commentEntity.ParentId.Int64 = req.ParentId
	commentEntity.ParentId.Valid = true

	idComment, err = u.repo.AddComment(ctx, commentEntity)
	if err != nil {
		return
	}

	return
}

func (u *usecase) UpdateComment(ctx context.Context, req request.UpdateCommentPayload, commentId string) (err error) {
	idCommentInt, err := strconv.Atoi(commentId)
	if err != nil {
		return
	}

	commentEntity := entity.NewFromUpdateCommentRequest(req, idCommentInt)

	err = u.repo.UpdateCommentById(ctx, commentEntity)
	if err != nil {
		return
	}

	return
}
func (u *usecase) DeleteComment(ctx context.Context, commentId string) (err error) {
	idCommentInt, err := strconv.Atoi(commentId)
	if err != nil {
		return
	}

	err = u.repo.DeleteCommentById(ctx, idCommentInt)
	if err != nil {
		return
	}

	return
}
