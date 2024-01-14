package logic

import (
	"context"
	"fmt"

	user_service "github.com/kilicmu/user-service/github.com/kilicmu/user-service"
	"github.com/kilicmu/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type VarifyTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVarifyTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VarifyTokenLogic {
	return &VarifyTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *VarifyTokenLogic) VarifyToken(in *user_service.VarifyTokenReq) (*user_service.VarifyTokenResp, error) {
	// todo: add your logic here and delete this line
	accessToken := in.GetAccessToken()
	unsafedUserInfo, err := l.svcCtx.JWTService.UnsafeGetUserInfoPayload(accessToken)
	if err != nil {
		return nil, err
	}
	userInfo, err := l.svcCtx.UserManageService.GetUserByUid(unsafedUserInfo.UId)
	if err != nil {
		return nil, fmt.Errorf("get user info error: %s", err)
	}
	isAlive, err := l.svcCtx.JWTService.VarifyAccessTokenAlive(accessToken, userInfo)
	if err != nil {
		return nil, fmt.Errorf("varify token alive error: %s", err)
	}
	return &user_service.VarifyTokenResp{
		IsAlive: isAlive,
	}, nil
}
