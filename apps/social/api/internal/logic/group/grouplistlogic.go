package group

import (
	"context"
	"github.com/jinzhu/copier"
	"go-chat/apps/social/rpc/socialclient"
	"go-chat/pkg/ctxdata"

	"go-chat/apps/social/api/internal/svc"
	"go-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户申请列表
func NewGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupListLogic {
	return &GroupListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupListLogic) GroupList(req *types.GroupListReq) (resp *types.GroupListResp, err error) {
	uid := ctxdata.GetId(l.ctx)
	list, err := l.svcCtx.Social.GroupList(l.ctx, &socialclient.GroupListReq{
		UserId: uid,
	})
	if err != nil {
		return nil, err
	}
	var respList []*types.Groups
	copier.Copy(&respList, list.List)
	return &types.GroupListResp{List: respList}, nil
}
