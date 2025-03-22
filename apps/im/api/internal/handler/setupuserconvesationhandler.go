package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-chat/apps/im/api/internal/logic"
	"go-chat/apps/im/api/internal/svc"
	"go-chat/apps/im/api/internal/types"
)

// 建立会话列表
func setUpUserConvesationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SetUpUserConvesationReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewSetUpUserConvesationLogic(r.Context(), svcCtx)
		resp, err := l.SetUpUserConvesation(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
