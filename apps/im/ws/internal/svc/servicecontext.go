package svc

import (
	"go-chat/apps/im/immodels"
	"go-chat/apps/im/ws/internal/config"
	"go-chat/apps/task/mq/mqclient"
)

type ServiceContext struct {
	Config config.Config
	immodels.ChatLogModel
	mqclient.MsgChatTransferClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                c,
		MsgChatTransferClient: mqclient.NewMsgChatTransferClient(c.MsgChatTransfer.Addrs, c.MsgChatTransfer.Topic),
		ChatLogModel:          immodels.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
	}
}
