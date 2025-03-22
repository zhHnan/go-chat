package logic

import (
	"context"
	"go-chat/apps/im/immodels"
	"go-chat/pkg/xerr"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"go-chat/apps/im/rpc/im"
	"go-chat/apps/im/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsLogic {
	return &GetConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取会话
func (l *GetConversationsLogic) GetConversations(in *im.GetConversationsReq) (*im.GetConversationsResp, error) {
	// 添加调试日志
	l.Logger.Infof("尝试查询用户会话: userId=%s", in.UserId)
	list, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, in.UserId)
	if err != nil {
		if err == immodels.ErrNotFound {
			return &im.GetConversationsResp{}, nil
		}
		return nil, errors.Wrapf(xerr.NewDBErr(), "find chatlog by msgId err【%v】, req【%v】", err, in.UserId)
	}
	// 根据会话列表，查询会话详情
	var res im.GetConversationsResp
	copier.Copy(&res, &list)
	ids := make([]string, 0, len(list.ConversationList))
	// 获取id
	for _, v := range list.ConversationList {
		ids = append(ids, v.ConversationId)
	}
	conversations, err := l.svcCtx.ConversationModel.ListByConversationIds(l.ctx, ids)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find chatlog by msgId err【%v】, req【%v】", err, ids)
	}
	// 计算是否存在未读消息
	for _, conversation := range conversations {
		// 获取会话未读消息
		// 会话为空时，跳过
		if _, ok := res.ConversationList[conversation.ConversationId]; !ok {
			continue
		}
		// 用户读取的消息量
		total := res.ConversationList[conversation.ConversationId].Total
		if total < int32(conversation.Total) {
			// 有新的消息, 更新总的消息量
			res.ConversationList[conversation.ConversationId].Total = int32(conversation.Total)
			// 设置未读消息
			res.ConversationList[conversation.ConversationId].ToRead = int32(conversation.Total) - total
			// 更改当前会话为显示状态
			res.ConversationList[conversation.ConversationId].IsShow = true
		}
	}
	return &res, nil
}
