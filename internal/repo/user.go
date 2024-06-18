package repo

import (
	"context"

	"kratos-starter/internal/data/model"
)

type UserRepo interface {
	ListAll(ctx context.Context) ([]*model.User, error)
}
