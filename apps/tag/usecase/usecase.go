package usecase

import (
	"context"
	"mohhefni/go-blog-app/apps/tag/entity"
	"mohhefni/go-blog-app/apps/tag/repository"
)

type Usecase interface {
	GetTags(ctx context.Context, tagSearch string) (nameTag []entity.TagsList, err error)
}

type usecase struct {
	repo repository.Repository
}

func NewUsecase(repo repository.Repository) Usecase {
	return &usecase{
		repo: repo,
	}
}

func (u *usecase) GetTags(ctx context.Context, tagSearch string) (nameTag []entity.TagsList, err error) {
	nameTag, err = u.repo.GetTagByName(ctx, tagSearch)
	if err != nil {
		return
	}

	return
}
