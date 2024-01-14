package logic

import (
	"context"
	"fmt"
	"log"

	user_service "github.com/kilicmu/user-service/github.com/kilicmu/user-service"
	"github.com/kilicmu/user-service/internal/database"
	"github.com/kilicmu/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

type UpdateUserInfoLogicParams struct {
	UId      string `validate:"required"`
	Name     string `validate:"omitempty,gt=4"`
	Email    string `validate:"omitempty,email"`
	Picture  string `validate:"omitempty,url"`
	Password string `validate:"omitempty,gt=6,lt=20"`
	Phone    string `validate:"omitempty,numeric"`
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(in *user_service.UpdatableUserInfo) (*user_service.UserInfoDTO, error) {
	UId := in.GetUid()
	Name := in.GetName()
	Email := in.GetEmail()
	Picture := in.GetPicture()
	Password := in.GetPassword()
	Phone := in.GetPhone()

	if err := validate.Struct(UpdateUserInfoLogicParams{
		UId,
		Name,
		Email,
		Picture,
		Password,
		Phone,
	}); err != nil {
		return nil, err
	}

	shouldUpdataUserinfo := &database.UserInfo{
		UId:      UId,
		Name:     Name,
		Email:    Email,
		Picture:  Picture,
		Phone:    Phone,
		Password: Password,
	}

	updatedUserinfo, err := l.svcCtx.UserManageService.UpdateUserInfo(shouldUpdataUserinfo)
	if err != nil {
		log.Printf("update user info failed: err: %s", err)
		return nil, fmt.Errorf("update user info failed")
	}
	return &user_service.UserInfoDTO{
		Uid:           updatedUserinfo.UId,
		Name:          updatedUserinfo.Name,
		Email:         updatedUserinfo.Email,
		Picture:       updatedUserinfo.Picture,
		Channel:       updatedUserinfo.Channel,
		Identify:      updatedUserinfo.Identify,
		EmailVerified: updatedUserinfo.EmailVerified,
		Phone:         updatedUserinfo.Phone,
	}, nil
}
