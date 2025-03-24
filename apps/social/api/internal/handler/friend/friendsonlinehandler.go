package friend

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-chat/apps/social/api/internal/logic/friend"
	"go-chat/apps/social/api/internal/svc"
	"go-chat/apps/social/api/internal/types"
)

// 查询好友在线状态
func FriendsOnlineHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FriendOnlineReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := friend.NewFriendsOnlineLogic(r.Context(), svcCtx)
		resp, err := l.FriendsOnline(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
