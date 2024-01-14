package userservicelogic

import (
	"context"

	"github.com/kilicmu/user-service/internal/svc"
	"github.com/kilicmu/user-service/pb/github.com/kilicmu/user-service"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsUserExistLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsUserExistLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsUserExistLogic {
	return &IsUserExistLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsUserExistLogic) IsUserExist(in *user_service.UserInfoReq) (*user_service.ResultResp, error) {
	// todo: add your logic here and delete this line

	return &user_service.ResultResp{}, nil
}
