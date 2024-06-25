package biz

import (
	"context"
	"fmt"
	"time"

	"kratos-starter/internal/constant"
	"kratos-starter/internal/data/model"
	"kratos-starter/internal/repo"
	"kratos-starter/pkg/utils"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
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
	uid, err := utils.GetUid(ctx)
	if err != nil {
		return nil, err
	}

	uc.log.WithContext(ctx).Infof("uid: %d", uid)

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

	pwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(pwd)

	return uc.userRepo.Create(ctx, user)
}

func (uc *UserUsecase) GetUser(ctx context.Context, id uint) (*model.User, error) {
	return uc.userRepo.GetById(ctx, id)
}

func (uc *UserUsecase) Login(ctx context.Context, user *model.User) (string, error) {
	en, err := uc.userRepo.GetByName(ctx, user.UserName)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(en.Password), []byte(user.Password)); err != nil {
		uc.log.WithContext(ctx).Error(err.Error())
		return "", fmt.Errorf("password is incorrect")
	}

	sec := lo.RandomString(10, lo.AlphanumericCharset)

	if err := uc.userRepo.SetUserToken(ctx, en.ID, sec, time.Hour*24); err != nil {
		return "", err
	}

	return jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, jwtv5.RegisteredClaims{
		Issuer:    sec,
		Subject:   "Subject",
		Audience:  []string{"aud1", "aud2"},
		ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(time.Hour * 24)),
		NotBefore: jwtv5.NewNumericDate(time.Now()),
		IssuedAt:  jwtv5.NewNumericDate(time.Now()),
		ID:        fmt.Sprintf("%d", en.ID),
	}).SignedString([]byte(constant.JwtSignKey))
}
