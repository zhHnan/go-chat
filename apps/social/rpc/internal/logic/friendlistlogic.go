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

type FriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendListLogic) FriendList(in *social.FriendListReq) (*social.FriendListResp, error) {
	// 获取好友列表
	friendList, err := l.svcCtx.FriendsModel.ListByUserId(l.ctx, in.UserId)
	if err != nil {
		l.Logger.Errorf("获取好友列表失败: %v, userId: %s", err, in.UserId)
		return nil, errors.Wrapf(xerr.NewDBErr(), "list friends by userId err %v req %v", err, in.UserId)
	}
	// 转换结果
	var res []*social.Friends
	if len(friendList) > 0 {
		err = copier.Copy(&res, friendList)
		if err != nil {
			l.Logger.Errorf("复制好友列表失败: %v", err)
			return nil, errors.Wrapf(xerr.NewDBErr(), "copy friends list err %v", err)
		}
	}
	return &social.FriendListResp{
		List: res,
	}, nil
}
