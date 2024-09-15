package usecase

import (
	"context"
	"mohhefni/go-blog-app/apps/interaction/entity"
	"mohhefni/go-blog-app/apps/interaction/repository"
	"mohhefni/go-blog-app/apps/interaction/request"
	"mohhefni/go-blog-app/utility"
	"strconv"
)

type Usecase interface {
	CreateInteractionLike(ctx context.Context, req request.AddInteractionRequestPayload, idUser string) (idInteraction int, err error)
	CreateInteractionShare(ctx context.Context, req request.AddInteractionRequestPayload, publicId string) (idInteraction int, err error)
	CreateInteractionBookmark(ctx context.Context, req request.AddInteractionRequestPayload, publicId string) (idInteraction int, err error)
	DeleteInteraction(ctx context.Context, interactionId string) (err error)
}

type usecase struct {
	repo repository.Repository
}

func NewUsecase(repo repository.Repository) Usecase {
	return &usecase{
		repo: repo,
	}
}

func (u *usecase) CreateInteractionLike(ctx context.Context, req request.AddInteractionRequestPayload, publicId string) (idInteraction int, err error) {
	interactionEntity := entity.NewFromAddInteractionLikeRequest(req, string(entity.LIKE_TYPE))
	publicIdUuid, err := utility.ParseUUID(publicId)
	if err != nil {
		return
	}

	var userEntity entity.UserEntity
	userEntity, err = u.repo.GetUserByPublicId(ctx, publicIdUuid)
	if err != nil {
		return
	}

	interactionEntity.UserId = userEntity.UserId

	idInteraction, err = u.repo.AddInteractions(ctx, interactionEntity)
	if err != nil {
		return
	}

	return
}

func (u *usecase) CreateInteractionShare(ctx context.Context, req request.AddInteractionRequestPayload, publicId string) (idInteraction int, err error) {
	interactionEntity := entity.NewFromAddInteractionLikeRequest(req, string(entity.SHARE_TYPE))
	publicIdUuid, err := utility.ParseUUID(publicId)
	if err != nil {
		return
	}

	var userEntity entity.UserEntity
	userEntity, err = u.repo.GetUserByPublicId(ctx, publicIdUuid)
	if err != nil {
		return
	}

	interactionEntity.UserId = userEntity.UserId

	idInteraction, err = u.repo.AddInteractions(ctx, interactionEntity)
	if err != nil {
		return
	}

	return
}

func (u *usecase) CreateInteractionBookmark(ctx context.Context, req request.AddInteractionRequestPayload, publicId string) (idInteraction int, err error) {
	interactionEntity := entity.NewFromAddInteractionLikeRequest(req, string(entity.BOOKMARK_TYPE))
	publicIdUuid, err := utility.ParseUUID(publicId)
	if err != nil {
		return
	}

	var userEntity entity.UserEntity
	userEntity, err = u.repo.GetUserByPublicId(ctx, publicIdUuid)
	if err != nil {
		return
	}

	interactionEntity.UserId = userEntity.UserId

	idInteraction, err = u.repo.AddInteractions(ctx, interactionEntity)
	if err != nil {
		return
	}

	return
}

func (u *usecase) DeleteInteraction(ctx context.Context, interactionId string) (err error) {
	interactionIdInt, err := strconv.Atoi(interactionId)
	if err != nil {
		return
	}

	err = u.repo.DeleteInteractionById(ctx, interactionIdInt)
	if err != nil {
		return
	}

	return
}
