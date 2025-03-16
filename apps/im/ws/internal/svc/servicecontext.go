package svc

import (
	"go-chat/apps/im/immodels"
	"go-chat/apps/im/ws/internal/config"
)

type ServiceContext struct {
	Config config.Config
	immodels.ChatLogModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		ChatLogModel: immodels.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
	}
}
