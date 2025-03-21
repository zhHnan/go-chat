package conversation

import (
	"github.com/mitchellh/mapstructure"
	"go-chat/apps/im/ws/internal/svc"
	"go-chat/apps/im/ws/websocket"
	"go-chat/apps/im/ws/ws"
	"go-chat/apps/task/mq/mq"
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
			err := svc.MsgChatTransferClient.Push(&mq.MsgChatTransfer{
				ChatType:       data.ChatType,
				ConversationId: data.ConversationId,
				SendId:         conn.Uid,
				ReceiveId:      data.ReceiveId,
				SendTime:       time.Now().UnixNano(),
				MType:          data.Msg.MType,
				Content:        data.Msg.Content,
			})
			if err != nil {
				srv.Send(websocket.NewErrMessage(err), conn)
				return
			}
			//err := logic.NewConversation(context.Background(), srv, svc).SingleChat(&data, conn.Uid)
			//if err != nil {
			//	srv.Send(websocket.NewErrMessage(err), conn)
			//	return
			//}
			//srv.SendByUserId(websocket.NewMessage(conn.Uid, ws.Chat{
			//	ChatType:       data.ChatType,
			//	ConversationId: data.ConversationId,
			//	SendId:         conn.Uid,
			//	ReceiveId:      data.ReceiveId,
			//	SendTime:       time.Now().UnixMilli(),
			//	Msg:            data.Msg,
			//}), data.ReceiveId)
		}
	}
}
