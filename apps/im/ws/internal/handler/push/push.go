package push

import (
	"encoding/json"
	"go-chat/apps/im/ws/internal/svc"
	"go-chat/apps/im/ws/websocket"
	"go-chat/apps/im/ws/ws"
	"go-chat/pkg/constants"

	"github.com/mitchellh/mapstructure"
)

func Push(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		// 添加调试日志
		jsonData, _ := json.Marshal(msg.Data)
		srv.Infof("收到消息: %s", string(jsonData))

		var data ws.Push
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Errorf("解码消息失败: %v, 原始数据: %s", err, string(jsonData))
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}

		// 添加调试日志
		srv.Infof("解码后的消息: sendId=%s, receiveId=%s, content=%s",
			data.SendId, data.ReceiveId, data.Content)

		// 打印所有当前在线用户
		srv.Infof("当前在线用户: %v", srv.GetUsers())

		// 发送的目标
		switch data.ChatType {
		case constants.PrivateChatType:
			single(srv, &data, data.ReceiveId)
		case constants.GroupChatType:
			group(srv, &data)
		}
		//receiveConn := srv.GetConn(data.ReceiveId)
		//if receiveConn == nil {
		//	// 添加调试日志
		//	srv.Errorf("目标用户不在线: %s", data.ReceiveId)
		//	// 检查接收者ID格式是否有问题
		//	srv.Infof("检查接收者ID格式: [%s]", data.ReceiveId)
		//
		//	// 尝试使用不同格式查找接收者
		//	for _, uid := range srv.GetUsers() {
		//		srv.Infof("比较在线用户: [%s] vs [%s]", uid, data.ReceiveId)
		//	}
		//	// todo 目标离线
		//	return
		//}
		//
		//srv.Infof("找到接收者连接, 推送消息")
		//
		//message := websocket.NewMessage(data.SendId, ws.Chat{
		//	ConversationId: data.ConversationId,
		//	ChatType:       data.ChatType,
		//	SendTime:       data.SendTime,
		//	Msg: ws.Msg{
		//		Content: data.Content,
		//		MType:   data.MType,
		//	},
		//})
		//
		//// 添加调试日志
		//messageJSON, _ := json.Marshal(message)
		//srv.Infof("发送消息: %s", string(messageJSON))
		//
		//if err := srv.Send(message, receiveConn); err != nil {
		//	srv.Errorf("发送消息失败: %v", err)
		//} else {
		//	srv.Infof("发送消息成功")
		//}
	}
}

func single(srv *websocket.Server, data *ws.Push, receiveId string) error {
	// 发送的目标
	receiveConn := srv.GetConn(receiveId)
	if receiveConn == nil {
		// 添加调试日志
		srv.Errorf("目标用户不在线: %s", data.ReceiveId)
		// 检查接收者ID格式是否有问题
		srv.Infof("检查接收者ID格式: [%s]", data.ReceiveId)
		// 尝试使用不同格式查找接收者
		for _, uid := range srv.GetUsers() {
			srv.Infof("比较在线用户: [%s] vs [%s]", uid, data.ReceiveId)
		}
		// todo 目标离线
		return nil
	}
	srv.Infof("找到接收者连接, 推送消息")
	message := websocket.NewMessage(data.SendId, ws.Chat{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendTime:       data.SendTime,
		Msg: ws.Msg{
			Content: data.Content,
			MType:   data.MType,
		},
	})
	// 添加调试日志
	messageJSON, _ := json.Marshal(message)
	srv.Infof("发送消息: %s", string(messageJSON))

	return srv.Send(message, receiveConn)
}

func group(srv *websocket.Server, data *ws.Push) error {
	for _, uid := range data.ReceiveIds {
		func(id string) {
			srv.Schedule(func() {
				single(srv, data, id)
			})
		}(uid)
	}
	return nil
}
