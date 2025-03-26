package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"go-chat/apps/im/rpc/imclient"
	"go-chat/apps/social/api/internal/config"
	"go-chat/apps/social/rpc/socialclient"
	"go-chat/apps/user/rpc/userclient"
	"go-chat/pkg/interceptor"
	"go-chat/pkg/middleware"
)

var retryPolicy = `{
"methodConfig":[{
	"name":[{
		"service":"social.social"
	}],
	"waitForReady": true,
	"retryPolicy":{
		"maxAttempts":5,
		"initialBackoff":"0.001s",
		"maxBackoff":"0.002s",
		"backoffMultiplier":1.0,
		"retryableStatusCodes":["UNKNOWN", "DEADLINE_EXCEEDED"]
		}
	}]
}`

type ServiceContext struct {
	Config                config.Config
	IdempotenceMiddleware rest.Middleware
	LimitMiddleware       rest.Middleware
	socialclient.Social
	userclient.User
	imclient.Im
	Redis *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                c,
		IdempotenceMiddleware: middleware.NewIdempotenceMiddleware().Handler,
		LimitMiddleware:       middleware.NewLimitMiddleware(c.Redisx).TokenLimitHandler(1, 100, "SOCIAL_LIMIT"),
		Social: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc,
			//zrpc.WithDialOption(grpc.WithDefaultServiceConfig(retryPolicy)),
			zrpc.WithUnaryClientInterceptor(interceptor.DefaultIdempotentClient))),
		User:  userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		Im:    imclient.NewIm(zrpc.MustNewClient(c.ImRpc)),
		Redis: redis.MustNewRedis(c.Redisx),
	}
}
