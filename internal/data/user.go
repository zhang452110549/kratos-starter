package data

import (
	"context"

	"kratos-starter/internal/data/model"
	"kratos-starter/internal/repo"

	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func (u userRepo) ListAll(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	return users, u.db.WithContext(ctx).Find(&users).Error
}

func NewUserRepo(data *Data) repo.UserRepo {
	return &userRepo{
		db: data.db,
	}
}
