package logic

import (
	"context"
	"fmt"

	user_service "github.com/kilicmu/user-service/github.com/kilicmu/user-service"
	"github.com/kilicmu/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

type LoginLogicParams struct {
	nameOrPhoneOrEmail string `validate:"required"`
	Password           string `validate:"required,gt=4"`
}

func (l *LoginLogic) Login(in *user_service.LoginReq) (*user_service.AccessTokenPair, error) {
	nameOrPhoneOrEmail := in.GetNameOrPhoneOrEmail()
	password := in.GetPassword()

	if err := validate.Struct(LoginLogicParams{
		nameOrPhoneOrEmail: nameOrPhoneOrEmail,
		Password:           password,
	}); err != nil {
		return nil, err
	}

	userinfo, err := l.svcCtx.UserManageService.GetUserInfoByEmailOrPhoneOrName(nameOrPhoneOrEmail)
	if err != nil {
		return nil, fmt.Errorf("login error: %s", err)
	}
	if l.svcCtx.UserManageService.ValidateUserPassword(password, userinfo) {
		return nil, fmt.Errorf("unmatched password")
	}
	fmt.Println(userinfo)
	if pair, err := l.svcCtx.JWTService.GenerateJwtTokenPairByUserInfo(userinfo); err == nil {
		return &user_service.AccessTokenPair{
			AccessToken:  pair.AccessToken,
			RefreshToken: pair.RefreshToken,
		}, nil
	}
	return nil, fmt.Errorf("login error: %s", err)
}
