package friend

import (
	"context"
	"go-chat/apps/social/rpc/socialclient"
	"go-chat/apps/user/rpc/userclient"
	"go-chat/pkg/ctxdata"
	"strconv"

	"go-chat/apps/social/api/internal/svc"
	"go-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友列表
func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendListLogic) FriendList(req *types.FriendListReq) (resp *types.FriendListResp, err error) {
	// 好友列表查询
	uid := ctxdata.GetId(l.ctx)
	friends, err := l.svcCtx.Social.FriendList(l.ctx, &socialclient.FriendListReq{
		UserId: uid,
	})
	if err != nil {
		return nil, err
	}
	if len(friends.List) == 0 {
		return &types.FriendListResp{}, nil
	}
	// 根据id 获取好友的信息
	uids := make([]string, 0, len(friends.List))
	for _, v := range friends.List {
		uids = append(uids, v.FriendUid)
	}
	// 获取用户信息
	users, err := l.svcCtx.User.FindUser(l.ctx, &userclient.FindUserReq{
		Ids: uids,
	})
	if err != nil {
		return &types.FriendListResp{}, nil
	}
	userRecords := make(map[string]*userclient.UserEntity, len(users.Users))
	for i, _ := range users.Users {
		userRecords[users.Users[i].Id] = users.Users[i]
	}

	respList := make([]*types.Friends, 0, len(friends.List))
	for _, v := range friends.List {
		fid, _ := strconv.ParseInt(v.FriendUid, 10, 32)
		friend := &types.Friends{
			Id:        v.Id,
			FriendUid: int32(fid),
		}

		if u, ok := userRecords[v.FriendUid]; ok {
			friend.Nickname = u.Nickname
			friend.Avatar = u.Avatar
		}
		respList = append(respList, friend)
	}
	return &types.FriendListResp{
		List: respList,
	}, nil
}
