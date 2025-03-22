package constants

// MType 定义自定义整数类型以表示不同的消息类型。
type MType int

// TextMtype 表示文本的 message 类型.
const (
	TextMtype MType = iota
)

type ChatType int

const (
	GroupChatType ChatType = iota + 1
	PrivateChatType
)

type ContentType int

const (
	ContentChatMsg ContentType = iota
	ContentMakeRead
)
