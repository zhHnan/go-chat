package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"go-chat/apps/im/rpc/imclient"
	"go-chat/pkg/ctxdata"

	"go-chat/apps/im/api/internal/svc"
	"go-chat/apps/im/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConvesationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取会话列表
func NewGetConvesationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConvesationsLogic {
	return &GetConvesationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetConvesationsLogic) GetConvesations(req *types.GetConvesationsReq) (resp *types.GetConvesationsResp, err error) {
	uid := ctxdata.GetId(l.ctx)
	data, err := l.svcCtx.GetConversations(l.ctx, &imclient.GetConversationsReq{
		UserId: uid,
	})
	if err != nil {
		return nil, err
	}

	var res types.GetConvesationsResp
	copier.Copy(&res, &data)
	return &res, err
}
