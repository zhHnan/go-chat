package mqclient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-queue/kq"
	"go-chat/apps/task/mq/mq"
)

type MsgChatTransferClient interface {
	Push(msg *mq.MsgChatTransfer) error
}

type msgChatTransferClient struct {
	pusher *kq.Pusher
}

func NewMsgChatTransferClient(addr []string, topic string, opts ...kq.PushOption) MsgChatTransferClient {
	fmt.Println("addrs", addr, "topic", topic)
	return &msgChatTransferClient{
		pusher: kq.NewPusher(addr, topic),
	}
}
func (m *msgChatTransferClient) Push(msg *mq.MsgChatTransfer) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return m.pusher.Push(context.Background(), string(body))
}
