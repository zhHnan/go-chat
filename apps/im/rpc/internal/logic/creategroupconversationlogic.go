package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-chat/apps/im/immodels"
	"go-chat/apps/im/rpc/im"
	"go-chat/apps/im/rpc/internal/svc"
	"go-chat/pkg/constants"
	"go-chat/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateGroupConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupConversationLogic {
	return &CreateGroupConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateGroupConversationLogic) CreateGroupConversation(in *im.CreateGroupConversationReq) (*im.CreateGroupConversationResp, error) {
	// todo: add your logic here and delete this line
	res := &im.CreateGroupConversationResp{}

	_, err := l.svcCtx.ConversationModel.FindOne(l.ctx, in.GroupId)

	if err == nil {
		return res, nil
	}

	if err != immodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "get conversation err 【%v】, req 【%v】", err, in.GroupId)
	}

	err = l.svcCtx.ConversationModel.Insert(l.ctx, &immodels.Conversation{
		ChatType:       constants.GroupChatType,
		ConversationId: in.GroupId,
	})

	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "conversation.insert err 【%v】", err)
	}

	_, err = NewSetUpUserConversationLogic(l.ctx, l.svcCtx).SetUpUserConversation(&im.SetUpUserConversationReq{
		SendId:   in.CreateId,
		RecvId:   in.GroupId,
		ChatType: int32(constants.GroupChatType),
	})
	return res, err
}
