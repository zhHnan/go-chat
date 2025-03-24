package main

import (
	"flag"
	"fmt"
	"go-chat/pkg/interceptor"
	"go-chat/pkg/interceptor/rpcserver"

	"go-chat/apps/social/rpc/internal/config"
	"go-chat/apps/social/rpc/internal/server"
	"go-chat/apps/social/rpc/internal/svc"
	"go-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/dev/social.yaml", "the config file")

//var configFile = flag.String("f", "apps/social/rpc/etc/dev/social.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		social.RegisterSocialServer(grpcServer, server.NewSocialServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	s.AddUnaryInterceptors(rpcserver.LoginInterceptor)
	s.AddUnaryInterceptors(interceptor.NewIdempotentServer(interceptor.NewDefaultIdempotent(c.Cache[0].RedisConf)))
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
