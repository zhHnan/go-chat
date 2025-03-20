package main

import (
	"flag"
	"fmt"
	"go-chat/apps/im/ws/internal/config"
	"go-chat/apps/im/ws/internal/handler"
	"go-chat/apps/im/ws/internal/svc"
	"go-chat/apps/im/ws/websocket"

	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/dev/im.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	if err := c.SetUp(); err != nil {
		panic(err)
	}
	ctx := svc.NewServiceContext(c)
	srv := websocket.NewServer(c.ListenOn, websocket.WithServerAuthentication(handler.NewJwtAuth(ctx)),
		websocket.WithServerAck(websocket.OnlyAck))
	defer srv.Stop()

	handler.RegisterHandlers(srv, ctx)

	fmt.Printf("Starting server at %s...\n", c.ListenOn)
	srv.Start()
}
