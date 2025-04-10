package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-chat/apps/social/socialmodels"
	"go-chat/pkg/constants"
	"go-chat/pkg/wuid"
	"go-chat/pkg/xerr"

	"go-chat/apps/social/rpc/internal/svc"
	"go-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupCreateLogic {
	return &GroupCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupCreateLogic) GroupCreate(in *social.GroupCreateReq) (*social.GroupCreateResp, error) {
	// 创建群组
	groups := &socialmodels.Groups{
		Id:         wuid.GenUid(l.svcCtx.Config.Mysql.DataSource),
		Name:       in.Name,
		Icon:       in.Icon,
		CreatorUid: in.CreatorUid,
		IsVerify:   false,
	}
	// 修改这里：处理Trans返回的错误
	err := l.svcCtx.GroupsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		_, err := l.svcCtx.GroupsModel.Insert(ctx, session, groups)
		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "groups insert err %v, req %v", err, groups)
		}
		_, err = l.svcCtx.GroupMembersModel.Insert(l.ctx, session, &socialmodels.GroupMembers{
			GroupId:   groups.Id,
			UserId:    in.CreatorUid,
			RoleLevel: int64(constants.CreatorGroupRoleLevel),
		})
		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "group members insert err %v, req %v", err, groups)
		}
		return nil
	})
	// 添加错误处理
	if err != nil {
		l.Logger.Errorf("创建群组失败: %v", err)
		return nil, errors.Wrapf(xerr.NewDBErr(), "create group failed: %v", err)
	}
	// 返回创建的群组ID
	return &social.GroupCreateResp{
		Id: groups.Id,
	}, nil
}
