package logic

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"go-chat/apps/social/socialmodels"
	"go-chat/pkg/constants"
	"go-chat/pkg/xerr"
	"time"

	"go-chat/apps/social/rpc/internal/svc"
	"go-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutInLogic) FriendPutIn(in *social.FriendPutInReq) (*social.FriendPutInResp, error) {
	// 申请人是否与目标是好友关系
	friends, err := l.svcCtx.FriendsModel.FindByUidAndFid(l.ctx, in.ReqUid, in.UserId)
	if err != nil && err != socialmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find user by uid and fid error %v, req %v", err, in)
	}
	if friends != nil {
		return &social.FriendPutInResp{}, err
	}
	// 是否已经有过申请，申请是不成功，没有完成
	friendsReq, err := l.svcCtx.FriendRequestsModel.FindByUidAndReqUid(l.ctx, in.UserId, in.ReqUid)
	if err != nil && err != socialmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find user by uid and fid error %v, req %v", err, in)
	}
	if friendsReq != nil {
		return &social.FriendPutInResp{}, err
	}
	// 创建好友申请记录
	_, err = l.svcCtx.FriendRequestsModel.Insert(l.ctx, &socialmodels.FriendRequests{
		UserId:       in.UserId,
		ReqUid:       in.ReqUid,
		ReqMsg:       sql.NullString{String: in.ReqMsg, Valid: true},
		ReqTime:      time.Unix(in.ReqTime, 0),
		HandleResult: sql.NullInt64{Int64: int64(constants.NoHandlerResult), Valid: true},
		HandleMsg:    sql.NullString{String: "", Valid: true},
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert friend request error %v, req %v", err, in)
	}
	return &social.FriendPutInResp{}, nil
}
