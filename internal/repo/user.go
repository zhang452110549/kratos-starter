package repo

import (
	"context"
	"time"

	"kratos-starter/internal/data/model"
)

type UserRepo interface {
	ListAll(ctx context.Context) ([]*model.User, error)
	Create(ctx context.Context, user *model.User) error
	GetByName(ctx context.Context, name string) (*model.User, error)
	GetById(ctx context.Context, id uint) (*model.User, error)
	SetUserToken(ctx context.Context, userId uint, sec string, ttl time.Duration) error
}
