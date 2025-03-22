package mq

import "go-chat/pkg/constants"

type MsgChatTransfer struct {
	ConversationId     string `json:"conversationId"`
	constants.ChatType `json:"chatType"`
	constants.MType    `json:"mtype"`
	SendId             string   `json:"sendId"`
	SendTime           int64    `json:"sendTime"`
	ReceiveId          string   `json:"receiveId"`
	ReceiveIds         []string `json:"receiveIds"`
	Content            string   `json:"content"`
}

// 增加一个消费者处理已读消息
type MsgMarkRead struct {
	constants.ChatType `json:"chatType"`
	SendId             string   `json:"sendId"`
	ReceiveId          string   `json:"receiveId"`
	ConversationId     string   `json:"conversationId"`
	MsgIds             []string `json:"msgIds"`
}
