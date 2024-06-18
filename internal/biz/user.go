package biz

import (
	"context"

	"kratos-starter/internal/data/model"
	"kratos-starter/internal/repo"

	"github.com/go-kratos/kratos/v2/log"
)

type UserUsecase struct {
	userRepo repo.UserRepo
	log      *log.Helper
}

func NewUserUsecase(repo repo.UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{
		userRepo: repo,
		log:      log.NewHelper(log.With(logger, "module", "biz/user")),
	}
}

func (uc *UserUsecase) ListAllUsers(ctx context.Context) ([]*model.User, error) {
	return uc.userRepo.ListAll(ctx)
}
