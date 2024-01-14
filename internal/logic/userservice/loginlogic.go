package userservicelogic

import (
	"context"

	"github.com/kilicmu/user-service/internal/svc"
	"github.com/kilicmu/user-service/pb/github.com/kilicmu/user-service"

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

func (l *LoginLogic) Login(in *user_service.LoginReq) (*user_service.LoginResp, error) {
	// todo: add your logic here and delete this line

	return &user_service.LoginResp{}, nil
}
