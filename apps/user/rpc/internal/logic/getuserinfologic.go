package logic

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	"go-chat/apps/user/models"
	"go-chat/apps/user/rpc/internal/svc"
	"go-chat/apps/user/rpc/user"

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

func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoReq) (*user.GetUserInfoResp, error) {
	// todo: add your logic here and delete this line
	userEntity, err := l.svcCtx.UsersModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == models.ErrNotFound {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	var res user.UserEntity
	copier.Copy(&res, userEntity)
	return &user.GetUserInfoResp{
		User: &res,
	}, nil
}
