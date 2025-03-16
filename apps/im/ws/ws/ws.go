package ws

import "go-chat/pkg/constants"

// 定义发送消息的格式
// Msg定义了一个消息结构体，包含了消息类型和内容。
// MType字段通过mapstructure标签指定，用于在结构体与map之间进行转换时映射到正确的字段。
// Content字段保存了消息的具体内容。
type (
	Msg struct {
		constants.MType `mapstructure:"mType"`
		Content         string `mapstructure:"content"`
	}

	// Chat 定义了一个聊天消息的结构体，包括会话ID、发送者ID、接收者ID、消息内容和发送时间。
	// ConversationId字段用于唯一标识一次会话。
	// SendId和ReceiveId字段分别标识了消息的发送者和接收者。
	// Msg字段是一个匿名嵌入的Msg结构体，用于存储消息的具体内容。
	// SendIme字段记录了消息的发送时间戳。
	Chat struct {
		ConversationId     string `mapstructure:"conversationId"`
		constants.ChatType `mapstructure:"chatType"`
		SendId             string `mapstructure:"sendId"`
		ReceiveId          string `mapstructure:"receiveId"`
		Msg                `mapstructure:"msg"`
		SendTime           int64 `mapstructure:"sendTime"`
	}
)
