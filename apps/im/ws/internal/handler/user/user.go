package user

import (
	"go-chat/apps/im/ws/internal/svc"
	"go-chat/apps/im/ws/websocket"
)

func Online(svc *svc.ServiceContext) websocket.HandlerFunc {
	// todo
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		uids := srv.GetUsers(conn)
		u := srv.GetUsers(conn)
		err := srv.Send(websocket.NewMessage(u[0], uids), conn)
		srv.Info("err", err)
	}
}
