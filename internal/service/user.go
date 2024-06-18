package service

import (
	"context"

	userPb "kratos-starter/api/v1/user"
	"kratos-starter/internal/biz"
)

type UserService struct {
	userPb.UnimplementedUserServiceServer

	userUc *biz.UserUsecase
}

func (s *UserService) ListAllUsers(ctx context.Context, request *userPb.ListALlUserRequest) (*userPb.ListALlUserResponse, error) {
	users, err := s.userUc.ListAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	var list []*userPb.User

	for _, user := range users {
		list = append(list, &userPb.User{
			Id:       uint32(user.ID),
			UserName: user.UserName,
		})
	}

	return &userPb.ListALlUserResponse{
		Users: list,
	}, nil
}

func NewUserService(_userUc *biz.UserUsecase) *UserService {
	return &UserService{userUc: _userUc}
}
