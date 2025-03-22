package logic

import (
	"context"

	"go-chat/apps/im/api/internal/svc"
	"go-chat/apps/im/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetUpUserConvesationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 建立会话列表
func NewSetUpUserConvesationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUpUserConvesationLogic {
	return &SetUpUserConvesationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetUpUserConvesationLogic) SetUpUserConvesation(req *types.SetUpUserConvesationReq) (resp *types.SetUpUserConvesationResp, err error) {
	// todo: add your logic here and delete this line

	return
}
