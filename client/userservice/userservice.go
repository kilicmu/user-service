// Code generated by goctl. DO NOT EDIT.
// Source: user-service.proto

package userservice

import (
	"context"

	"github.com/kilicmu/user-service/pb/github.com/kilicmu/user-service"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	ActiveUserCountResp    = user_service.ActiveUserCountResp
	LoginReq               = user_service.LoginReq
	LoginResp              = user_service.LoginResp
	LoginResp_Data         = user_service.LoginResp_Data
	ResultResp             = user_service.ResultResp
	UserInfoDTO            = user_service.UserInfoDTO
	UserInfoReq            = user_service.UserInfoReq
	UserInfoResp           = user_service.UserInfoResp
	ValidateUserIsAliveReq = user_service.ValidateUserIsAliveReq
	VoidData               = user_service.VoidData

	UserService interface {
		Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginResp, error)
		RegisterUser(ctx context.Context, in *UserInfoDTO, opts ...grpc.CallOption) (*ResultResp, error)
		UserInfo(ctx context.Context, in *UserInfoReq, opts ...grpc.CallOption) (*UserInfoResp, error)
		IsUserBeingAlive(ctx context.Context, in *ValidateUserIsAliveReq, opts ...grpc.CallOption) (*ResultResp, error)
		IsUserExist(ctx context.Context, in *UserInfoReq, opts ...grpc.CallOption) (*ResultResp, error)
		ActiveUserCount(ctx context.Context, in *VoidData, opts ...grpc.CallOption) (*ActiveUserCountResp, error)
		UpdateUserInfo(ctx context.Context, in *UserInfoDTO, opts ...grpc.CallOption) (*ResultResp, error)
	}

	defaultUserService struct {
		cli zrpc.Client
	}
)

func NewUserService(cli zrpc.Client) UserService {
	return &defaultUserService{
		cli: cli,
	}
}

func (m *defaultUserService) Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginResp, error) {
	client := user_service.NewUserServiceClient(m.cli.Conn())
	return client.Login(ctx, in, opts...)
}

func (m *defaultUserService) RegisterUser(ctx context.Context, in *UserInfoDTO, opts ...grpc.CallOption) (*ResultResp, error) {
	client := user_service.NewUserServiceClient(m.cli.Conn())
	return client.RegisterUser(ctx, in, opts...)
}

func (m *defaultUserService) UserInfo(ctx context.Context, in *UserInfoReq, opts ...grpc.CallOption) (*UserInfoResp, error) {
	client := user_service.NewUserServiceClient(m.cli.Conn())
	return client.UserInfo(ctx, in, opts...)
}

func (m *defaultUserService) IsUserBeingAlive(ctx context.Context, in *ValidateUserIsAliveReq, opts ...grpc.CallOption) (*ResultResp, error) {
	client := user_service.NewUserServiceClient(m.cli.Conn())
	return client.IsUserBeingAlive(ctx, in, opts...)
}

func (m *defaultUserService) IsUserExist(ctx context.Context, in *UserInfoReq, opts ...grpc.CallOption) (*ResultResp, error) {
	client := user_service.NewUserServiceClient(m.cli.Conn())
	return client.IsUserExist(ctx, in, opts...)
}

func (m *defaultUserService) ActiveUserCount(ctx context.Context, in *VoidData, opts ...grpc.CallOption) (*ActiveUserCountResp, error) {
	client := user_service.NewUserServiceClient(m.cli.Conn())
	return client.ActiveUserCount(ctx, in, opts...)
}

func (m *defaultUserService) UpdateUserInfo(ctx context.Context, in *UserInfoDTO, opts ...grpc.CallOption) (*ResultResp, error) {
	client := user_service.NewUserServiceClient(m.cli.Conn())
	return client.UpdateUserInfo(ctx, in, opts...)
}
