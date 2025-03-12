package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go-chat/pkg/xerr"

	"go-chat/apps/social/rpc/internal/svc"
	"go-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupListLogic {
	return &GroupListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupListLogic) GroupList(in *social.GroupListReq) (*social.GroupListResp, error) {
	userGroup, err := l.svcCtx.GroupMembersModel.ListByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "groupmembers list by userId error %v, req %v", err, in.UserId)
	}
	if len(userGroup) == 0 {
		return &social.GroupListResp{}, nil
	}
	ids := make([]string, 0, len(userGroup))

	for _, v := range userGroup {
		ids = append(ids, v.GroupId)
	}
	groups, err := l.svcCtx.ListByGroupIds(l.ctx, ids)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "groups list by groupIds error %v, req %v", err, ids)
	}
	var resp []*social.Groups
	copier.Copy(&resp, groups)
	return &social.GroupListResp{
		List: resp,
	}, nil
}
