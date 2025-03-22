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

type PutConvesationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新会话列表
func NewPutConvesationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PutConvesationsLogic {
	return &PutConvesationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PutConvesationsLogic) PutConvesations(req *types.PutConvesationsReq) (resp *types.PutConvesationsResp, err error) {
	//flowchart TD
	//A[开始] --> B[获取用户ID]
	//B --> C[初始化会话列表变量]
	//C --> D[复制请求中的会话列表]
	//D --> E[调用服务更新会话]
	//E --> F[返回结果]

	uid := ctxdata.GetId(l.ctx)
	var conversationList map[string]*imclient.Conversation

	copier.Copy(&conversationList, req.ConversationList)
	_, err = l.svcCtx.PutConversations(l.ctx, &imclient.PutConversationsReq{
		UserId:           uid,
		ConversationList: conversationList,
	})

	return
}
