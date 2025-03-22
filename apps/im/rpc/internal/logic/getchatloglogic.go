package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-chat/pkg/xerr"

	"go-chat/apps/im/rpc/im"
	"go-chat/apps/im/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetChatLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatLogLogic {
	return &GetChatLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取会话记录
func (l *GetChatLogLogic) GetChatLog(in *im.GetChatLogReq) (*im.GetChatLogResp, error) {
	// 根据id
	if in.MsgId != "" {
		// 根据id获取会话记录
		chatLog, err := l.svcCtx.ChatLogModel.FindOne(l.ctx, in.MsgId)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewDBErr(), "find chatlog by msgId err【%v】, req【%v】", err, in.MsgId)
		}
		return &im.GetChatLogResp{
			List: []*im.ChatLog{
				{
					Id:             chatLog.ID.Hex(),
					ConversationId: chatLog.ConversationId,
					SendId:         chatLog.SendId,
					RecvId:         chatLog.ReceiveId,
					MsgType:        int32(chatLog.MsgType),
					MsgContent:     chatLog.MsgContent,
					ChatType:       int32(chatLog.ChatType),
					SendTime:       chatLog.SendTime,
					ReadRecords:    chatLog.ReadRecords,
				},
			},
		}, nil
	}
	// 	根据时间段查询
	data, err := l.svcCtx.ChatLogModel.ListBySendTime(l.ctx, in.ConversationId, in.StartSendTime, in.EndSendTime, in.Count)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find chatlog by sendTime err【%v】, req【%v】", err, in)
	}
	res := make([]*im.ChatLog, 0, len(data))
	for _, v := range data {
		res = append(res, &im.ChatLog{
			Id:             v.ID.Hex(),
			ConversationId: v.ConversationId,
			SendId:         v.SendId,
			RecvId:         v.ReceiveId,
			MsgType:        int32(v.MsgType),
			MsgContent:     v.MsgContent,
			ChatType:       int32(v.ChatType),
			SendTime:       v.SendTime,
			ReadRecords:    v.ReadRecords,
		})
	}
	return &im.GetChatLogResp{
		List: res,
	}, nil
}
