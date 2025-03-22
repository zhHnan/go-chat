package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go-chat/apps/im/rpc/imclient"
	"go-chat/apps/social/api/internal/config"
	"go-chat/apps/social/rpc/socialclient"
	"go-chat/apps/user/rpc/userclient"
)

type ServiceContext struct {
	Config config.Config
	socialclient.Social
	userclient.User
	imclient.Im
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Social: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
		User:   userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		Im:     imclient.NewIm(zrpc.MustNewClient(c.ImRpc)),
	}
}
