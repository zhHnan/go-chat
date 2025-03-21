package msgTransfer

import (
	"context"
	"encoding/json"
	"fmt"
	"go-chat/apps/im/immodels"
	"go-chat/apps/im/ws/websocket"
	"go-chat/apps/social/rpc/socialclient"
	"go-chat/apps/task/mq/internal/svc"
	"go-chat/apps/task/mq/mq"
	"go-chat/pkg/constants"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type MsgChatTransfer struct {
	logx.Logger
	svc *svc.ServiceContext
}

func NewMsgChatTransfer(svc *svc.ServiceContext) *MsgChatTransfer {
	return &MsgChatTransfer{
		Logger: logx.WithContext(context.Background()),
		svc:    svc,
	}
}

func (m *MsgChatTransfer) Consume(ctx context.Context, key, value string) error {
	// 检查上下文是否被取消或超时
	if err := ctx.Err(); err != nil {
		m.Logger.Errorf("Context error: %v", err)
		return err
	}
	fmt.Println("key : ", key, " value : ", value)
	var data mq.MsgChatTransfer
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		m.Logger.Errorf("Failed to unmarshal message: %v", err)
		return err
	}
	// 记录消息
	if err := m.addChatLog(ctx, &data); err != nil {
		m.Logger.Errorf("Failed to add chat log: %v", err)
		return err
	}
	switch data.ChatType {
	case constants.PrivateChatType:
		// 单聊
		return m.single(&data)
	case constants.GroupChatType:
		// 群聊
		return m.group(ctx, &data)
	}
	return nil

}

func (m *MsgChatTransfer) addChatLog(ctx context.Context, data *mq.MsgChatTransfer) error {
	chatLog := immodels.ChatLog{
		ConversationId: data.ConversationId,
		SendId:         data.SendId,
		ReceiveId:      data.ReceiveId,
		MsgFrom:        0,
		ChatType:       data.ChatType,
		MsgType:        data.MType,
		MsgContent:     data.Content,
		SendTime:       time.Now().UnixNano(),
	}
	err := m.svc.ChatLogModel.Insert(ctx, &chatLog)
	if err != nil {
		return err
	}

	return m.svc.ConversationModel.UpdateMsg(ctx, &chatLog)
}
func (m *MsgChatTransfer) group(ctx context.Context, data *mq.MsgChatTransfer) error {
	// 查询群的用户
	users, err := m.svc.Social.GroupUsers(ctx, &socialclient.GroupUsersReq{
		GroupId: data.ReceiveId,
	})
	if err != nil {
		return err
	}
	data.ReceiveIds = make([]string, 0, len(users.List))

	for _, member := range users.List {
		// 跳过发送者
		if member.UserId == data.SendId {
			continue
		}
		data.ReceiveIds = append(data.ReceiveIds, member.UserId)
	}
	return m.svc.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FormId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})
}
func (m *MsgChatTransfer) single(data *mq.MsgChatTransfer) error {
	// 记录发送的消息内容
	message := websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FormId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	}
	messageJson, _ := json.Marshal(message)
	m.Logger.Infof("推送消息: %s", string(messageJson))

	// 确保接收者存在
	m.Logger.Infof("推送给用户: %s", data.ReceiveId)

	// 推送消息
	err := m.svc.WsClient.Send(message)
	if err != nil {
		m.Logger.Errorf("推送消息失败: %v", err)
		return err
	}
	m.Logger.Info("推送消息成功")
	return nil
}
