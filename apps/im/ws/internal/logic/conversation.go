package logic

import (
	"context"
	"go-chat/apps/im/immodels"
	"go-chat/apps/im/ws/internal/svc"
	"go-chat/apps/im/ws/websocket"
	"go-chat/apps/im/ws/ws"
	"go-chat/pkg/wuid"
	"time"
)

type Conversation struct {
	ctx context.Context
	srv *websocket.Server
	svc *svc.ServiceContext
}

func NewConversation(ctx context.Context, srv *websocket.Server, svc *svc.ServiceContext) *Conversation {
	return &Conversation{
		ctx: ctx,
		srv: srv,
		svc: svc,
	}
}

func (l *Conversation) SingleChat(data *ws.Chat, userId string) error {
	if data.ConversationId == "" {
		data.ConversationId = wuid.CombineId(userId, data.ReceiveId)
	}
	chatLog := &immodels.ChatLog{
		ConversationId: data.ConversationId,
		SendId:         userId,
		ReceiveId:      data.ReceiveId,
		MsgFrom:        0,
		ChatType:       data.ChatType,
		MsgType:        data.MType,
		MsgContent:     data.Content,
		SendTime:       time.Now().Unix(),
	}
	err := l.svc.ChatLogModel.Insert(l.ctx, chatLog)
	return err
}
