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
