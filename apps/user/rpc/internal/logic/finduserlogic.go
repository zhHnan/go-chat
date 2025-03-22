package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"go-chat/apps/user/models"

	"go-chat/apps/user/rpc/internal/svc"
	"go-chat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLogic {
	return &FindUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindUserLogic) FindUser(in *user.FindUserReq) (*user.FindUserResp, error) {
	var (
		userEntities []*models.Users
		err          error
		userEntity   *models.Users
	)
	if in.Phone != "" {
		userEntity, err = l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
		userEntities = append(userEntities, userEntity)
	} else if len(in.Ids) > 0 {
		userEntities, err = l.svcCtx.UsersModel.ListByIds(l.ctx, in.Ids)
	} else {
		userEntities, err = l.svcCtx.UsersModel.ListByName(l.ctx, in.Name)
	}
	if err != nil {
		return nil, err
	}
	var res []*user.UserEntity
	copier.Copy(&res, userEntities)
	return &user.FindUserResp{
		Users: res,
	}, nil
}
