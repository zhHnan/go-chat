package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-chat/apps/social/socialmodels"
	"go-chat/pkg/constants"
	"go-chat/pkg/xerr"

	"go-chat/apps/social/rpc/internal/svc"
	"go-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrFriendReqBeforePass   = xerr.NewMsg("好友申请已通过")
	ErrFriendReqBeforeReject = xerr.NewMsg("好友申请已拒绝")
)

type FriendPutInHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInHandleLogic {
	return &FriendPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutInHandleLogic) FriendPutInHandle(in *social.FriendPutInHandleReq) (*social.FriendPutInHandleResp, error) {
	// todo: add your logic here and delete this line
	// 获取好友申请记录
	friendReq, err := l.svcCtx.FriendRequestsModel.FindOne(l.ctx, uint64(in.FriendReqId))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friendsRequest by friendReqId err %v req %v", err, in.FriendReqId)
	}
	// 验证是否有处理
	switch constants.HandlerResult(friendReq.HandleResult.Int64) {
	case constants.PassHandlerResult:
		return nil, errors.WithStack(ErrFriendReqBeforePass)
	case constants.RefuseHandlerResult:
		return nil, errors.WithStack(ErrFriendReqBeforeReject)
	}
	friendReq.HandleResult.Int64 = int64(in.HandleResult)
	// 修改申请结果 -> 通过【建立两条好有关系记录】 -> 事务
	l.svcCtx.FriendRequestsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		if err := l.svcCtx.FriendRequestsModel.Update(ctx, session, friendReq); err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "update friendsRequest err %v, req %v", err, friendReq)
		}
		if constants.HandlerResult(in.HandleResult) != constants.PassHandlerResult {
			return nil
		}
		// 如果申请为通过，则建立两条好友关系记录
		friends := []*socialmodels.Friends{
			{
				UserId:    friendReq.UserId,
				FriendUid: friendReq.ReqUid,
			},
			{
				UserId:    friendReq.ReqUid,
				FriendUid: friendReq.UserId,
			},
		}
		_, err = l.svcCtx.FriendsModel.InsertTrans(ctx, session, friends...)
		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "insert friends err %v, req %v", err, friends)
		}
		return nil
	})
	return &social.FriendPutInHandleResp{}, nil
}
