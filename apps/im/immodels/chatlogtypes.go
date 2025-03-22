package immodels

import (
	"go-chat/pkg/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatLog struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ConversationId string             `bson:"conversationId"`
	SendId         string             `bson:"sendId"`
	ReceiveId      string             `bson:"receiveId"`
	MsgFrom        int                `bson:"msgFrom"`
	ChatType       constants.ChatType `bson:"chatType"`
	MsgType        constants.MType    `bson:"msgType"`
	MsgContent     string             `bson:"msgContent"`
	SendTime       int64              `bson:"sendTime"`
	Status         int                `bson:"status"`
	// 存储消息已读记录
	ReadRecords []byte `bson:"readRecords"`

	UpdateAt time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
