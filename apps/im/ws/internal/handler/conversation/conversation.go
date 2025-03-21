package conversation

import (
	"github.com/mitchellh/mapstructure"
	"go-chat/apps/im/ws/internal/svc"
	"go-chat/apps/im/ws/websocket"
	"go-chat/apps/im/ws/ws"
	"go-chat/apps/task/mq/mq"
	"go-chat/pkg/constants"
	"go-chat/pkg/wuid"
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
		if data.ConversationId == "" {
			switch data.ChatType {
			case constants.PrivateChatType:
				data.ConversationId = wuid.CombineId(conn.Uid, data.ReceiveId)
			case constants.GroupChatType:
				data.ConversationId = data.ReceiveId
			}
		}
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
	}
}
