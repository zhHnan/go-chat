package msgTransfer

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"go-chat/apps/im/ws/websocket"
	"go-chat/apps/im/ws/ws"
	"go-chat/apps/social/rpc/socialclient"
	"go-chat/apps/task/mq/internal/svc"
	"go-chat/pkg/constants"
)

type baseMsgTransfer struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewBaseMsgTransfer(svc *svc.ServiceContext) *baseMsgTransfer {
	return &baseMsgTransfer{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svc,
	}
}

func (m *baseMsgTransfer) Transfer(ctx context.Context, data *ws.Push) error {
	var err error
	switch data.ChatType {
	case constants.PrivateChatType:
		err = m.single(ctx, data)
	case constants.GroupChatType:
		err = m.group(ctx, data)
	}
	return err
}

// single 私聊
func (m *baseMsgTransfer) single(ctx context.Context, data *ws.Push) error {
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
	err := m.svcCtx.WsClient.Send(message)
	if err != nil {
		m.Logger.Errorf("推送消息失败: %v", err)
		return err
	}
	m.Logger.Info("推送消息成功")
	return nil
}

// group 群聊
func (m *baseMsgTransfer) group(ctx context.Context, data *ws.Push) error {
	// 查询群的用户
	users, err := m.svcCtx.Social.GroupUsers(ctx, &socialclient.GroupUsersReq{
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
	return m.svcCtx.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FormId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})
}
