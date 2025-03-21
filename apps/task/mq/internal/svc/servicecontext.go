package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"go-chat/apps/im/immodels"
	"go-chat/apps/im/ws/websocket"
	"go-chat/apps/social/rpc/socialclient"
	"go-chat/apps/task/mq/internal/config"
	"go-chat/pkg/constants"
	"net/http"
)

// ServiceContext 引入配置文件
type ServiceContext struct {
	Config   config.Config
	WsClient websocket.Client
	*redis.Redis
	socialclient.Social
	immodels.ChatLogModel
	immodels.ConversationModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	svc := ServiceContext{
		Config:            c,
		Redis:             redis.MustNewRedis(c.Redisx),
		ConversationModel: immodels.MustConversationModel(c.Mongo.Url, c.Mongo.Db),
		ChatLogModel:      immodels.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
		Social:            socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
	}
	token, err := svc.GetSystemToken()
	if err != nil {
		panic(err)
	}
	header := http.Header{}
	header.Set("Authorization", token)
	svc.WsClient = websocket.NewClient(c.Ws.Host, websocket.WithClientHeader(header))
	return &svc
}

func (svc *ServiceContext) GetSystemToken() (string, error) {
	return svc.Redis.Get(constants.REDIS_SYSTEM_ROOT_TOKEN)
}
