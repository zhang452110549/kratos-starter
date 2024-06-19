package data

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"kratos-starter/internal/data/model"
	"kratos-starter/internal/repo"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type userRepo struct {
	log *log.Helper
	db  *gorm.DB
	rdb redis.UniversalClient
}

func (u userRepo) GetById(ctx context.Context, id uint) (user *model.User, err error) {
	if str, err := u.rdb.Get(ctx, fmt.Sprintf("user:%d", id)).Result(); err != nil {
		u.log.WithContext(ctx).Errorf("redis get user:%d failed: %v", id, err)
	} else {
		if err := json.Unmarshal([]byte(str), &user); err != nil {
			u.log.WithContext(ctx).Errorf("json unmarshal user:%d failed: %v", id, err)
		} else {
			return user, nil
		}
	}

	if err := u.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}

	if data, err := json.Marshal(user); err != nil {
		u.log.WithContext(ctx).Errorf("json marshal user:%d failed: %v", id, err)
	} else {
		if err := u.rdb.Set(ctx, fmt.Sprintf("user:%d", id), data, 10*time.Minute).Err(); err != nil {
			u.log.WithContext(ctx).Errorf("redis set user:%d failed: %v", id, err)
		}
	}

	return user, nil
}

func (u userRepo) GetByName(ctx context.Context, name string) (*model.User, error) {
	var user *model.User
	return user, u.db.WithContext(ctx).Where("f_user_name", name).First(&user).Error
}

func (u userRepo) Create(ctx context.Context, user *model.User) error {
	return u.db.WithContext(ctx).Create(user).Error
}

func (u userRepo) ListAll(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	return users, u.db.WithContext(ctx).Find(&users).Error
}

func NewUserRepo(data *Data, logger log.Logger) repo.UserRepo {
	return &userRepo{
		log: log.NewHelper(log.With(logger, "module", "data/user")),
		db:  data.db,
		rdb: data.rdb,
	}
}
