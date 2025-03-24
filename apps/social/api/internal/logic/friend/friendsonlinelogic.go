package friend

import (
	"context"
	"go-chat/apps/social/rpc/social"
	"go-chat/pkg/constants"
	"go-chat/pkg/ctxdata"

	"go-chat/apps/social/api/internal/svc"
	"go-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendsOnlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询好友在线状态
func NewFriendsOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendsOnlineLogic {
	return &FriendsOnlineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendsOnlineLogic) FriendsOnline(req *types.FriendOnlineReq) (resp *types.FriendOnlineResp, err error) {
	// todo: add your logic here and delete this line
	uid := ctxdata.GetId(l.ctx)
	list, err := l.svcCtx.Social.FriendList(l.ctx, &social.FriendListReq{UserId: uid})
	if err != nil || len(list.List) == 0 {
		return &types.FriendOnlineResp{}, err
	}
	ids := make([]string, 0, len(list.List))
	for _, friend := range list.List {
		ids = append(ids, friend.FriendUid)
	}

	onlineUsers, err := l.svcCtx.Redis.Hgetall(constants.REDIS_ONLINE_USER)
	if err != nil {
		return &types.FriendOnlineResp{}, err
	}
	resOnlineList := make(map[string]bool, len(ids))
	for _, s := range ids {
		if _, ok := onlineUsers[s]; ok {
			resOnlineList[s] = true
		} else {
			resOnlineList[s] = false
		}
	}
	return &types.FriendOnlineResp{
		OnlineList: resOnlineList,
	}, err
}
