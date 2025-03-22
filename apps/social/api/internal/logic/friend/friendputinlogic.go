package friend

import (
	"context"
	"github.com/pkg/errors"
	"go-chat/apps/social/rpc/socialclient"
	"go-chat/pkg/ctxdata"
	"go-chat/pkg/xerr"

	"go-chat/apps/social/api/internal/svc"
	"go-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友申请
func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendPutInLogic) FriendPutIn(req *types.FriendPutInReq) (resp *types.FriendPutInResp, err error) {

	// 获取用户的uid
	userId := ctxdata.GetId(l.ctx)
	_, err = l.svcCtx.Social.FriendPutIn(l.ctx, &socialclient.FriendPutInReq{
		UserId:  userId,
		ReqUid:  req.UserId,
		ReqMsg:  req.ReqMsg,
		ReqTime: req.ReqTime,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "friend put in err %v, req %v", err, req)
	}
	return
}
