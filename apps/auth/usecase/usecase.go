package usecase

import (
	"context"
	"mohhefni/go-blog-app/apps/auth/repository"
	"mohhefni/go-blog-app/apps/auth/request"
)

type Usecase interface {
	RegisterUser(ctx context.Context, req request.RegisterRequestPayload) (email string, err error)
}

type usecase struct {
	repo repository.Repository
}

func NewUsecase(repo repository.Repository) Usecase {
	return &usecase{
		repo: repo,
	}
}

func (u *usecase) RegisterUser(ctx context.Context, req request.RegisterRequestPayload) (email string, err error) {
	return
}
