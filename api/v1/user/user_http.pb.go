// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.5.0
// - protoc             v3.20.1
// source: v1/user/user.proto

package user

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationUserServiceListAllUsers = "/api.v1.user.UserService/ListAllUsers"

type UserServiceHTTPServer interface {
	ListAllUsers(context.Context, *ListALlUserRequest) (*ListALlUserResponse, error)
}

func RegisterUserServiceHTTPServer(s *http.Server, srv UserServiceHTTPServer) {
	r := s.Route("/")
	r.GET("/v1/users/all", _UserService_ListAllUsers0_HTTP_Handler(srv))
}

func _UserService_ListAllUsers0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListALlUserRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceListAllUsers)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListAllUsers(ctx, req.(*ListALlUserRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListALlUserResponse)
		return ctx.Result(200, reply)
	}
}

type UserServiceHTTPClient interface {
	ListAllUsers(ctx context.Context, req *ListALlUserRequest, opts ...http.CallOption) (rsp *ListALlUserResponse, err error)
}

type UserServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewUserServiceHTTPClient(client *http.Client) UserServiceHTTPClient {
	return &UserServiceHTTPClientImpl{client}
}

func (c *UserServiceHTTPClientImpl) ListAllUsers(ctx context.Context, in *ListALlUserRequest, opts ...http.CallOption) (*ListALlUserResponse, error) {
	var out ListALlUserResponse
	pattern := "/v1/users/all"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationUserServiceListAllUsers))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}