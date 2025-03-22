package immodels

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Conversations struct {
	ID               primitive.ObjectID       `bson:"_id,omitempty" json:"id,omitempty"`
	UserId           string                   `bson:"userId,omitempty" json:"userId,omitempty"`
	ConversationList map[string]*Conversation `bson:"conversationList,omitempty" json:"conversationList,omitempty"`
	UpdateAt         time.Time                `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt         time.Time                `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
