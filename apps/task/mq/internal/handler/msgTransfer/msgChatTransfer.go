package msgTransfer

import (
	"context"
	"encoding/json"
	"fmt"
	"go-chat/apps/im/immodels"
	"go-chat/apps/im/ws/ws"
	"go-chat/apps/task/mq/internal/svc"
	"go-chat/apps/task/mq/mq"
	"go-chat/pkg/bitmap"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type MsgChatTransfer struct {
	*baseMsgTransfer
}

func NewMsgChatTransfer(svc *svc.ServiceContext) *MsgChatTransfer {
	return &MsgChatTransfer{
		NewBaseMsgTransfer(svc),
	}
}

func (m *MsgChatTransfer) Consume(ctx context.Context, key, value string) error {
	// 检查上下文是否被取消或超时
	if err := ctx.Err(); err != nil {
		m.Logger.Errorf("Context error: %v", err)
		return err
	}
	fmt.Println("key : ", key, " value : ", value)
	var (
		data  mq.MsgChatTransfer
		msgId = primitive.NewObjectID()
	)
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		m.Logger.Errorf("Failed to unmarshal message: %v", err)
		return err
	}
	// 记录消息
	if err := m.addChatLog(ctx, msgId, &data); err != nil {
		m.Logger.Errorf("Failed to add chat log: %v", err)
		return err
	}
	return m.Transfer(ctx, &ws.Push{
		MsgId:          msgId.Hex(),
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		MType:          data.MType,
		SendId:         data.SendId,
		SendTime:       data.SendTime,
		ReceiveId:      data.ReceiveId,
		ReceiveIds:     data.ReceiveIds,
		Content:        data.Content,
	})

}

func (m *MsgChatTransfer) addChatLog(ctx context.Context, msgId primitive.ObjectID, data *mq.MsgChatTransfer) error {
	chatLog := immodels.ChatLog{
		ID:             msgId,
		ConversationId: data.ConversationId,
		SendId:         data.SendId,
		ReceiveId:      data.ReceiveId,
		MsgFrom:        0,
		ChatType:       data.ChatType,
		MsgType:        data.MType,
		MsgContent:     data.Content,
		SendTime:       time.Now().UnixNano(),
	}
	// 设置已读消息
	readRecords := bitmap.NewBitmap(0)
	readRecords.Set(chatLog.SendId)
	chatLog.ReadRecords = readRecords.Export()

	err := m.svcCtx.ChatLogModel.Insert(ctx, &chatLog)
	if err != nil {
		return err
	}

	return m.svcCtx.ConversationModel.UpdateMsg(ctx, &chatLog)
}

//func (m *MsgChatTransfer) group(ctx context.Context, data *mq.MsgChatTransfer) error {
//	// 查询群的用户
//	users, err := m.svc.Social.GroupUsers(ctx, &socialclient.GroupUsersReq{
//		GroupId: data.ReceiveId,
//	})
//	if err != nil {
//		return err
//	}
//	data.ReceiveIds = make([]string, 0, len(users.List))
//
//	for _, member := range users.List {
//		// 跳过发送者
//		if member.UserId == data.SendId {
//			continue
//		}
//		data.ReceiveIds = append(data.ReceiveIds, member.UserId)
//	}
//	return m.svc.WsClient.Send(websocket.Message{
//		FrameType: websocket.FrameData,
//		Method:    "push",
//		FormId:    constants.SYSTEM_ROOT_UID,
//		Data:      data,
//	})
//}
//func (m *MsgChatTransfer) single(data *mq.MsgChatTransfer) error {
//	// 记录发送的消息内容
//	message := websocket.Message{
//		FrameType: websocket.FrameData,
//		Method:    "push",
//		FormId:    constants.SYSTEM_ROOT_UID,
//		Data:      data,
//	}
//	messageJson, _ := json.Marshal(message)
//	m.Logger.Infof("推送消息: %s", string(messageJson))
//
//	// 确保接收者存在
//	m.Logger.Infof("推送给用户: %s", data.ReceiveId)
//
//	// 推送消息
//	err := m.svc.WsClient.Send(message)
//	if err != nil {
//		m.Logger.Errorf("推送消息失败: %v", err)
//		return err
//	}
//	m.Logger.Info("推送消息成功")
//	return nil
//}
