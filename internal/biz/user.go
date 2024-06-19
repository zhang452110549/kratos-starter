package biz

import (
	"context"
	"fmt"

	"kratos-starter/internal/data/model"
	"kratos-starter/internal/repo"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
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

func (uc *UserUsecase) CreateUser(ctx context.Context, user *model.User) error {
	en, err := uc.userRepo.GetByName(ctx, user.UserName)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if en.ID != 0 {
		return fmt.Errorf("user %s already exists", user.UserName)
	}

	return uc.userRepo.Create(ctx, user)
}

func (uc *UserUsecase) GetUser(ctx context.Context, id uint) (*model.User, error) {
	return uc.userRepo.GetById(ctx, id)
}
