package group

import (
	"context"
	"github.com/jinzhu/copier"
	"go-chat/apps/social/rpc/socialclient"

	"go-chat/apps/social/api/internal/svc"
	"go-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 群组申请列表
func NewGroupPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInListLogic {
	return &GroupPutInListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupPutInListLogic) GroupPutInList(req *types.GroupPutInListReq) (resp *types.GroupPutInListResp, err error) {
	// 调用GroupPutInList方法，将指定群组放入列表中
	// 该方法需要上下文和群组ID作为参数
	list, err := l.svcCtx.Social.GroupPutInList(l.ctx, &socialclient.GroupPutInListReq{
		GroupId: req.GroupId,
	})

	// 初始化一个类型为types.GroupRequests切片的变量respList
	var respList []*types.GroupRequests

	// 使用copier.Copy将list.List中的数据复制到respList中
	// 这里之所以使用copier库，是因为它能够方便地在不同结构体类型之间进行数据复制
	copier.Copy(&respList, list.List)

	// 返回GroupPutInListResp响应，其中包含复制后的群组请求列表respList
	// 错误处理被省略，正常情况下返回nil
	return &types.GroupPutInListResp{
		List: respList,
	}, nil
}
