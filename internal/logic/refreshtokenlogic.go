package logic

import (
	"context"
	"fmt"
	"log"

	user_service "github.com/kilicmu/user-service/github.com/kilicmu/user-service"
	"github.com/kilicmu/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefreshTokenLogic) RefreshToken(in *user_service.RefreshTokenReq) (*user_service.GoogleOAuthGetAccessResp, error) {
	refreshToken := in.GetRefreshToken()
	unsafedUserInfo, err := l.svcCtx.JWTService.UnsafeGetUserInfoPayload(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("get user info by refresh token error: %s", err)
	}
	uid := unsafedUserInfo.UId
	userInfo, err := l.svcCtx.UserManageService.GetUserByUid(uid)
	if err != nil {
		return nil, fmt.Errorf("get user info by uid error: %s", err)
	}

	result, err := l.svcCtx.JWTService.VarifyRefreshTokenAlive(refreshToken, userInfo)
	if !result {
		log.Printf("varify token alive error: %s", err)
		return nil, fmt.Errorf("refresh token unvalid")
	}

	tokenPair, err := l.svcCtx.JWTService.GenerateJwtTokenPairByUserInfo(userInfo)
	return &user_service.GoogleOAuthGetAccessResp{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}
