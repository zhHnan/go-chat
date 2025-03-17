package conversation

import (
	"context"
	"github.com/mitchellh/mapstructure"
	"go-chat/apps/im/ws/internal/logic"
	"go-chat/apps/im/ws/internal/svc"
	"go-chat/apps/im/ws/websocket"
	"go-chat/apps/im/ws/ws"
	"go-chat/pkg/constants"
	"time"
)

func Chat(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		// todo: 私聊
		var data ws.Chat
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}
		switch data.ChatType {
		case constants.PrivateChatType:
			err := logic.NewConversation(context.Background(), srv, svc).SingleChat(&data, conn.Uid)
			if err != nil {
				srv.Send(websocket.NewErrMessage(err), conn)
				return
			}
			srv.SendByUserId(websocket.NewMessage(conn.Uid, ws.Chat{
				ChatType:       data.ChatType,
				ConversationId: data.ConversationId,
				SendId:         conn.Uid,
				ReceiveId:      data.ReceiveId,
				SendTime:       time.Now().UnixMilli(),
				Msg:            data.Msg,
			}), data.ReceiveId)
		}
	}
}
