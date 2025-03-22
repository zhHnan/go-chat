package group

import (
	"context"
	"go-chat/apps/social/rpc/social"
	"go-chat/apps/user/rpc/userclient"

	"go-chat/apps/social/api/internal/svc"
	"go-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 群组用户列表
func NewGroupUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUserListLogic {
	return &GroupUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupUserListLogic) GroupUserList(req *types.GroupUserListReq) (resp *types.GroupUserListResp, err error) {
	groupUsers, err := l.svcCtx.Social.GroupUsers(l.ctx, &social.GroupUsersReq{
		GroupId: req.GroupId,
	})
	// 查询用户的信息
	uids := make([]string, 0, len(groupUsers.List))
	for _, v := range groupUsers.List {
		uids = append(uids, v.UserId)
	}
	userList, err := l.svcCtx.User.FindUser(l.ctx, &userclient.FindUserReq{Ids: uids})
	if err != nil {
		return nil, err
	}

	userRecords := make(map[string]*userclient.UserEntity, len(userList.Users))
	for i, _ := range userList.Users {
		userRecords[userList.Users[i].Id] = userList.Users[i]
	}
	respList := make([]*types.GroupMembers, 0, len(groupUsers.List))

	for _, v := range groupUsers.List {
		member := &types.GroupMembers{
			Id:        int64(v.Id),
			UserId:    v.UserId,
			GroupId:   v.GroupId,
			RoleLevel: int(v.RoleLevel),
		}

		if u, ok := userRecords[v.UserId]; ok {
			member.UserAvatarUrl = u.Avatar
			member.Nickname = u.Nickname
		}
		respList = append(respList, member)
	}
	return &types.GroupUserListResp{List: respList}, err
}
