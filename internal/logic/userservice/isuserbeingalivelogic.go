package userservicelogic

import (
	"context"

	"github.com/kilicmu/user-service/internal/svc"
	"github.com/kilicmu/user-service/pb/github.com/kilicmu/user-service"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsUserBeingAliveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsUserBeingAliveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsUserBeingAliveLogic {
	return &IsUserBeingAliveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsUserBeingAliveLogic) IsUserBeingAlive(in *user_service.ValidateUserIsAliveReq) (*user_service.ResultResp, error) {
	// todo: add your logic here and delete this line

	return &user_service.ResultResp{}, nil
}
