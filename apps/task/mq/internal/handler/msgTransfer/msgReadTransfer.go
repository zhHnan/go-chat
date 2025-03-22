package msgTransfer

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-queue/kq"
	"go-chat/apps/im/ws/ws"
	"go-chat/apps/task/mq/internal/svc"
	"go-chat/apps/task/mq/mq"
	"go-chat/pkg/bitmap"
	"go-chat/pkg/constants"
)

type MsgReadTransfer struct {
	*baseMsgTransfer
}

func NewMsgReadTransfer(svc *svc.ServiceContext) kq.ConsumeHandler {
	return &MsgReadTransfer{
		baseMsgTransfer: NewBaseMsgTransfer(svc),
	}
}
func (m *MsgReadTransfer) Consume(ctx context.Context, key, value string) error {
	fmt.Println("....")
	var data mq.MsgMarkRead

	if err := json.Unmarshal([]byte(value), &data); err != nil {
		m.Logger.Errorf("Failed to unmarshal message: %v", err)
		return err
	}
	// 业务处理
	readRecords, err := m.UpdateChatLogRead(ctx, &data)
	if err != nil {
		return err
	}
	// map[msgId]已读记录
	return m.Transfer(ctx, &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		ReceiveId:      data.ReceiveId,
		ContentType:    constants.ContentMakeRead,
		ReadRecords:    readRecords,
	})
}
func (m *MsgReadTransfer) UpdateChatLogRead(ctx context.Context, data *mq.MsgMarkRead) (map[string]string, error) {
	res := make(map[string]string)
	chatLogs, err := m.svcCtx.ChatLogModel.ListByMsgIds(ctx, data.MsgIds)
	if err != nil {
		return nil, err
	}
	// 处理已读
	for _, chatLog := range chatLogs {
		switch chatLog.ChatType {
		case constants.PrivateChatType:
			chatLog.ReadRecords = []byte{}
		case constants.GroupChatType:
			// 设置已读消息
			readRecords := bitmap.Load(chatLog.ReadRecords)
			readRecords.Set(data.SendId)
			chatLog.ReadRecords = readRecords.Export()
		}
		// 保障数据在传输中的精度
		res[chatLog.ID.Hex()] = base64.StdEncoding.EncodeToString(chatLog.ReadRecords)
		err := m.svcCtx.ChatLogModel.UpdateMakeRead(ctx, chatLog.ID, chatLog.ReadRecords)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
