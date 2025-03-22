package logic

import (
	"context"
	"go-chat/apps/im/immodels"
	"go-chat/pkg/constants"
	"go-chat/pkg/wuid"
	"go-chat/pkg/xerr"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go-chat/apps/im/rpc/im"
	"go-chat/apps/im/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetUpUserConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetUpUserConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUpUserConversationLogic {
	return &SetUpUserConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 建立会话：群聊、私聊
func (l *SetUpUserConversationLogic) SetUpUserConversation(in *im.SetUpUserConversationReq) (*im.SetUpUserConversationResp, error) {
	// 根据会话的类型创建会话
	switch constants.ChatType(in.ChatType) {
	case constants.PrivateChatType:
		// 生成会话的id
		conversationId := wuid.CombineId(in.SendId, in.RecvId)
		// 验证是否建立过会话
		conRes, err := l.svcCtx.ConversationModel.FindOne(l.ctx, conversationId)
		if err != nil {
			if err == immodels.ErrNotFound {
				// 创建会话
				err = l.svcCtx.ConversationModel.Insert(l.ctx, &immodels.Conversation{
					ChatType:       constants.PrivateChatType,
					ConversationId: conversationId,
				})
				if err != nil {
					return nil, errors.Wrapf(xerr.NewDBErr(), "insert conversation error %v", err)
				}
			} else {
				return nil, errors.Wrapf(xerr.NewDBErr(), "find conversation error %v", err)
			}
		} else if conRes != nil {
			// 存在会话，直接返回
			return nil, nil
		}
		// 建立两者的会话
		//调用 setUpUserConversation 方法为发送方和接收方设置会话，
		//参数 true 表示发送方。
		//再次调用该方法，参数 false 表示接收方。如果任意一步出错，则返回错误信息。
		err = l.setUpUserConversation(conversationId, in.SendId, in.RecvId, constants.PrivateChatType, true)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewDBErr(), "set up user conversation error %v", err)
		}
		err = l.setUpUserConversation(conversationId, in.RecvId, in.SendId, constants.PrivateChatType, false)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewDBErr(), "set up user conversation error %v", err)
		}
	case constants.GroupChatType:
		err := l.setUpUserConversation(in.RecvId, in.SendId, in.RecvId, constants.GroupChatType, true)

		if err != nil {
			return nil, errors.Wrapf(xerr.NewDBErr(), "set up user conversation error %v", err)
		}
	}

	return &im.SetUpUserConversationResp{}, nil
}

// setUpUserConversation 创建会话
func (l *SetUpUserConversationLogic) setUpUserConversation(conversationId, userId, recvId string, chatType constants.ChatType, isShow bool) error {
	conversations, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, userId)
	if err != nil {
		if err == immodels.ErrNotFound {
			conversations = &immodels.Conversations{
				ID:               primitive.NewObjectID(),
				UserId:           userId,
				ConversationList: make(map[string]*immodels.Conversation),
			}
		} else {
			return errors.Wrapf(xerr.NewDBErr(), "find conversation error 【%v】, userId【%v】", err, userId)
		}
	}
	// 更新会话记录
	if _, ok := conversations.ConversationList[conversationId]; !ok {
		// 添加会话记录
		conversations.ConversationList[conversationId] = &immodels.Conversation{
			ChatType:       chatType,
			ConversationId: conversationId,
			IsShow:         isShow,
		}
	}
	// 更新
	_, err = l.svcCtx.ConversationsModel.Update(l.ctx, conversations)
	if err != nil {
		return errors.Wrapf(xerr.NewDBErr(), "update conversation error 【%v】, userId【%v】", err, userId)
	}
	return nil
}
