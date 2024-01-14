package userservicelogic

import (
	"context"

	"github.com/kilicmu/user-service/internal/svc"
	"github.com/kilicmu/user-service/pb/github.com/kilicmu/user-service"

	"github.com/zeromicro/go-zero/core/logx"
)

type ActiveUserCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewActiveUserCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ActiveUserCountLogic {
	return &ActiveUserCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ActiveUserCountLogic) ActiveUserCount(in *user_service.VoidData) (*user_service.ActiveUserCountResp, error) {
	// todo: add your logic here and delete this line

	return &user_service.ActiveUserCountResp{}, nil
}
