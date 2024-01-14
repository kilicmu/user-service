package logic

import (
	"context"
	"fmt"

	user_service "github.com/kilicmu/user-service/github.com/kilicmu/user-service"
	"github.com/kilicmu/user-service/internal/database"
	"github.com/kilicmu/user-service/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	IDENTIFY_USER  = "user"
	IDENTIFY_ADMIN = "admin"
)

const (
	GOOGLE_CHANNEL_MARK = "google"
)

type GoogleOAuthGetAccessTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGoogleOAuthGetAccessTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GoogleOAuthGetAccessTokenLogic {
	return &GoogleOAuthGetAccessTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GoogleOAuthGetAccessTokenLogic) GoogleOAuthGetAccessToken(in *user_service.GoogleOAuthGetAccessTokenReq) (*user_service.GoogleOAuthGetAccessResp, error) {
	code := in.GetCode()

	gUserInfo, err := l.svcCtx.OAuthService.GoogleOAuthService.GetGoogleUserInfoByAccessCode(code)
	if err != nil {
		return nil, fmt.Errorf("get user info by code %s error: %s", code, err)
	}

	userInfo, err := l.svcCtx.UserManageService.CreateUserIfNotExist(&database.UserInfo{
		UId:           gUserInfo.UId,
		Email:         gUserInfo.Email,
		Name:          gUserInfo.Name,
		EmailVerified: gUserInfo.EmailVerified,
		Picture:       gUserInfo.Picture,
		Channel:       GOOGLE_CHANNEL_MARK,
		Identify:      IDENTIFY_USER,
		Password:      "",
	})
	if err != nil {
		return nil, fmt.Errorf("google oauth get user info error: %s", err)
	}

	tokenPair, err := l.svcCtx.JWTService.GenerateJwtTokenPairByUserInfo(userInfo)
	if err != nil {
		return nil, fmt.Errorf("generate jwt error: %s", err)
	}

	return &user_service.GoogleOAuthGetAccessResp{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}
