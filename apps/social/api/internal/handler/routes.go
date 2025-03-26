// Code generated by goctl. DO NOT EDIT.
// goctl 1.8.1

package handler

import (
	"net/http"

	friend "go-chat/apps/social/api/internal/handler/friend"
	group "go-chat/apps/social/api/internal/handler/group"
	"go-chat/apps/social/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				// 好友申请
				Method:  http.MethodPost,
				Path:    "/friend/putIn",
				Handler: friend.FriendPutInHandler(serverCtx),
			},
			{
				// 好友申请处理
				Method:  http.MethodPut,
				Path:    "/friend/putIn",
				Handler: friend.FriendPutInHandleHandler(serverCtx),
			},
			{
				// 好友申请列表
				Method:  http.MethodGet,
				Path:    "/friend/putIns",
				Handler: friend.FriendPutInListHandler(serverCtx),
			},
			{
				// 好友列表
				Method:  http.MethodGet,
				Path:    "/friends",
				Handler: friend.FriendListHandler(serverCtx),
			},
			{
				// 查询好友在线状态
				Method:  http.MethodGet,
				Path:    "/friends/online",
				Handler: friend.FriendsOnlineHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.JwtAuth.AccessSecret),
		rest.WithPrefix("/v1/social"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.IdempotenceMiddleware, serverCtx.LimitMiddleware},
			[]rest.Route{
				{
					// 创建群组
					Method:  http.MethodPost,
					Path:    "/create",
					Handler: group.CreateGroupHandler(serverCtx),
				},
				{
					// 申请入群
					Method:  http.MethodPost,
					Path:    "/group/putIn",
					Handler: group.GroupPutInHandler(serverCtx),
				},
				{
					// 申请进群处理
					Method:  http.MethodPut,
					Path:    "/group/putIn",
					Handler: group.GroupPutInHandleHandler(serverCtx),
				},
				{
					// 群组申请列表
					Method:  http.MethodGet,
					Path:    "/group/putIns",
					Handler: group.GroupPutInListHandler(serverCtx),
				},
				{
					// 群组用户列表
					Method:  http.MethodGet,
					Path:    "/group/users",
					Handler: group.GroupUserListHandler(serverCtx),
				},
				{
					// 查询好友在线状态
					Method:  http.MethodGet,
					Path:    "/group/users/online",
					Handler: group.GroupUserOnlineHandler(serverCtx),
				},
				{
					// 用户申请列表
					Method:  http.MethodGet,
					Path:    "/groups",
					Handler: group.GroupListHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.JwtAuth.AccessSecret),
		rest.WithPrefix("/v1/social"),
	)
}
