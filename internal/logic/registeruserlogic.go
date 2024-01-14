package logic

import (
	"context"

	user_service "github.com/kilicmu/user-service/github.com/kilicmu/user-service"
	"github.com/kilicmu/user-service/internal/database"
	"github.com/kilicmu/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterUserLogic {
	return &RegisterUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

type RegisterUserInfoParams struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,gt=6,lt=20"`
}

func (l *RegisterUserLogic) RegisterUser(in *user_service.RegisterUserReq) (*user_service.UserInfoDTO, error) {
	if err := validate.Struct(RegisterUserInfoParams{
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	}); err != nil {
		return nil, err
	}
	userinfo := &database.UserInfo{
		Email:    in.GetEmail(),
		Password: l.svcCtx.UserManageService.EncodePassword(in.GetPassword()),
	}

	userinfo, err := l.svcCtx.UserManageService.CreateUserIfNotExist(userinfo)
	if err != nil {
		return nil, err
	}
	return &user_service.UserInfoDTO{
		Uid:           userinfo.UId,
		Name:          userinfo.Name,
		Email:         userinfo.Email,
		Picture:       userinfo.Picture,
		Channel:       userinfo.Channel,
		Identify:      userinfo.Identify,
		EmailVerified: userinfo.EmailVerified,
	}, nil
}
