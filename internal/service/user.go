package service

import (
	"context"

	userPb "kratos-starter/api/v1/user"
	"kratos-starter/internal/biz"
	"kratos-starter/internal/data/model"
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

func (s *UserService) CreateUser(ctx context.Context, request *userPb.CreateUserRequest) (*userPb.CreateUserResponse, error) {
	user := &model.User{
		UserName: request.GetUserName(),
	}

	return &userPb.CreateUserResponse{
		Id: uint32(user.ID),
	}, s.userUc.CreateUser(ctx, user)
}

func NewUserService(_userUc *biz.UserUsecase) *UserService {
	return &UserService{userUc: _userUc}
}

func (s *UserService) GetUser(ctx context.Context, request *userPb.GetUserRequest) (*userPb.GetUserResponse, error) {
	user, err := s.userUc.GetUser(ctx, uint(request.GetId()))
	if err != nil {
		return nil, err
	}

	return &userPb.GetUserResponse{
		User: &userPb.User{
			Id:       uint32(user.ID),
			UserName: user.UserName,
		},
	}, nil
}
