package handler

import (
	"go-chat/apps/im/ws/internal/handler/conversation"
	"go-chat/apps/im/ws/internal/handler/user"
	"go-chat/apps/im/ws/internal/svc"
	"go-chat/apps/im/ws/websocket"
)

func RegisterHandlers(srv *websocket.Server, svc *svc.ServiceContext) {
	srv.AddRoute([]websocket.Route{
		{
			Method:  "user.online",
			Handler: user.Online(svc),
		},
		{
			Method:  "conversation.chat",
			Handler: conversation.Chat(svc),
		},
	})
}
