package group

import (
	"context"
	"go-chat/apps/im/rpc/imclient"
	"go-chat/apps/social/rpc/socialclient"
	"go-chat/pkg/constants"
	"go-chat/pkg/ctxdata"

	"go-chat/apps/social/api/internal/svc"
	"go-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 申请入群
func NewGroupPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInLogic {
	return &GroupPutInLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupPutInLogic) GroupPutIn(req *types.GroupPutInReq) (resp *types.GroupPutInResp, err error) {
	uid := ctxdata.GetId(l.ctx)

	res, err := l.svcCtx.Social.GroupPutIn(l.ctx, &socialclient.GroupPutInReq{
		GroupId:    req.GroupId,
		ReqId:      uid,
		ReqMsg:     req.ReqMsg,
		ReqTime:    req.ReqTime,
		JoinSource: int32(req.JoinSource),
	})
	// 建立会话
	if res.GroupId == "" {
		return nil, err
	}
	// 建立会话
	_, err = l.svcCtx.Im.SetUpUserConversation(l.ctx, &imclient.SetUpUserConversationReq{
		SendId:   uid,
		RecvId:   res.GroupId,
		ChatType: int32(constants.GroupChatType),
	})
	return nil, err
}
