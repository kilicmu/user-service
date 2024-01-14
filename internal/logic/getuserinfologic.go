package logic

import (
	"context"
	"fmt"

	user_service "github.com/kilicmu/user-service/github.com/kilicmu/user-service"
	"github.com/kilicmu/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *user_service.UserInfoReq) (*user_service.UserInfoDTO, error) {
	// todo: add your logic here and delete this line
	accessToken := in.GetAccessToken()
	unsafedUserInfo, err := l.svcCtx.JWTService.UnsafeGetUserInfoPayload(accessToken)
	if err != nil {
		return nil, err
	}
	userInfo, err := l.svcCtx.UserManageService.GetUserByUid(unsafedUserInfo.UId)
	if err != nil {
		return nil, err
	}

	if result, _ := l.svcCtx.JWTService.VarifyAccessTokenAlive(accessToken, userInfo); !result {
		return nil, fmt.Errorf("token unvalid")
	}

	return &user_service.UserInfoDTO{
		Uid:           userInfo.UId,
		Name:          userInfo.Name,
		Email:         userInfo.Email,
		Picture:       userInfo.Picture,
		Channel:       userInfo.Channel,
		Identify:      userInfo.Identify,
		EmailVerified: userInfo.EmailVerified,
		Phone:         userInfo.Phone,
	}, nil
}
