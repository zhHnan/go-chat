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
	"sync"
	"time"
)

var (
	GroupMsgReadRecordDelayTime  = time.Second
	GroupMsgReadRecordDelayCount = 10
)

const (
	GroupMsgReadHandlerAtTransfer = iota
	GroupMsgReadHandlerDelayTransfer
)

type MsgReadTransfer struct {
	*baseMsgTransfer
	cache.Cache
	mu sync.Mutex
	// 用于存储和快速访问不同群组的消息信息。
	groupMsgs map[string]*groupMsgRead
	// 该通道用于向客户端发送推送通知或消息，支持异步通信。
	push chan *ws.Push
}

func NewMsgReadTransfer(svc *svc.ServiceContext) kq.ConsumeHandler {
	m := &MsgReadTransfer{
		baseMsgTransfer: NewBaseMsgTransfer(svc),
		groupMsgs:       make(map[string]*groupMsgRead, 1),
		push:            make(chan *ws.Push, 1),
	}
	if svc.Config.MsgReadHandler.GroupMsgReadHandler != GroupMsgReadHandlerAtTransfer {
		if svc.Config.MsgReadHandler.GroupMsgReadRecordDelayCount > 0 {
			GroupMsgReadRecordDelayCount = svc.Config.MsgReadHandler.GroupMsgReadRecordDelayCount
		}
		if svc.Config.MsgReadHandler.GroupMsgReadRecordDelayTime > 0 {
			GroupMsgReadRecordDelayTime = time.Duration(svc.Config.MsgReadHandler.GroupMsgReadRecordDelayTime)
		}
	}

	go m.transfer()
	return m
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
	push := &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		ReceiveId:      data.ReceiveId,
		ContentType:    constants.ContentMakeRead,
		ReadRecords:    readRecords,
	}
	switch push.ChatType {
	case constants.PrivateChatType:
		m.push <- push
	case constants.GroupChatType:
		// 判断是否开启合并消息的处理
		if m.svcCtx.Config.MsgReadHandler.GroupMsgReadHandler == GroupMsgReadHandlerAtTransfer {
			m.push <- push
		}

		if _, ok := m.groupMsgs[push.ConversationId]; ok {
			m.Infof("merge push: %v", push.ConversationId)
			// 合并请求
			m.groupMsgs[push.ConversationId].mergePush(push)
		} else {
			m.Infof("new push: %v", push.ConversationId)
			m.groupMsgs[push.ConversationId] = newGroupMsgRead(push, m.push)
		}
	}
	return nil

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
func (m *MsgReadTransfer) transfer() {

	for push := range m.push {
		if push.ReceiveId != "" || len(push.ReceiveIds) > 0 {
			if err := m.Transfer(context.Background(), push); err != nil {
				m.Errorf("Failed to transfer message: 【%v】push:【%v】", err, push)
			}
		}

		if push.ChatType == constants.PrivateChatType {
			continue
		}

		if m.svcCtx.Config.MsgReadHandler.GroupMsgReadHandler == GroupMsgReadHandlerAtTransfer {
			continue
		}
		// 清空数据
		m.mu.Lock()
		if _, ok := m.groupMsgs[push.ConversationId]; ok && m.groupMsgs[push.ConversationId].IsIdle() {
			m.groupMsgs[push.ConversationId].clear()
			delete(m.groupMsgs, push.ConversationId)
		}
		m.mu.Unlock()
	}
}
