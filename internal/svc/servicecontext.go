package svc

import (
	"github.com/kilicmu/user-service/internal/config"
	"github.com/kilicmu/user-service/internal/svc/jwt_svc"
	"github.com/kilicmu/user-service/internal/svc/oauth_svc"
	"github.com/kilicmu/user-service/internal/svc/user_manage_svc"
)

type ServiceContext struct {
	Config            config.Config
	OAuthService      *oauth_svc.OAuthService
	UserManageService *user_manage_svc.UserManageService
	JWTService        *jwt_svc.JWTService
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:            c,
		OAuthService:      oauth_svc.NewOAuthService(),
		UserManageService: user_manage_svc.NewUserManageService(),
		JWTService:        jwt_svc.NewJwtService(),
	}
}
