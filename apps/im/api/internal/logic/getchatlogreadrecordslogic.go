package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-chat/apps/im/rpc/im"
	"go-chat/apps/social/rpc/social"
	"go-chat/apps/user/rpc/user"
	"go-chat/pkg/bitmap"
	"go-chat/pkg/constants"
	"go-chat/pkg/xerr"

	"go-chat/apps/im/api/internal/svc"
	"go-chat/apps/im/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatLogReadRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取聊天记录已读未读记录
func NewGetChatLogReadRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatLogReadRecordsLogic {
	return &GetChatLogReadRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetChatLogReadRecordsLogic 获取聊天记录已读未读状态的逻辑处理
// 实现了根据消息ID获取聊天记录，并区分私聊和群聊来确定哪些用户已读或未读该消息
func (l *GetChatLogReadRecordsLogic) GetChatLogReadRecords(req *types.GetChatLogReadRecordsReq) (resp *types.GetChatLogReadRecordsResp, err error) {
	// 从IM服务获取聊天记录
	chatlogs, err := l.svcCtx.Im.GetChatLog(l.ctx, &im.GetChatLogReq{
		MsgId: req.MsgId,
	})
	if err != nil || len(chatlogs.List) == 0 {
		return nil, errors.Wrapf(xerr.NewDBErr(), "get chatlogs err 【%v】, req 【%v】", err, req)
	}

	var (
		chatlog = chatlogs.List[0]
		reads   = []string{chatlog.SendId}
		unreads []string
		ids     []string
	)
	// 分别设置已读未读消息
	switch constants.ChatType(chatlog.ChatType) {
	case constants.PrivateChatType:
		// 私聊情况下，根据读取记录设置已读和未读用户
		if len(chatlog.ReadRecords) == 0 || chatlog.ReadRecords[0] == 0 {
			unreads = []string{chatlog.RecvId}
		} else {
			reads = append(reads, chatlog.RecvId)
		}
		ids = []string{chatlog.SendId, chatlog.RecvId}
	case constants.GroupChatType:
		// 群聊情况下，获取群组成员并根据读取记录设置已读和未读用户
		groupUsers, err := l.svcCtx.Social.GroupUsers(l.ctx, &social.GroupUsersReq{
			GroupId: chatlog.ConversationId,
		})
		if err != nil {
			return nil, errors.Wrapf(xerr.NewDBErr(), "get group users err 【%v】, req 【%v】", err, req)
		}

		bitmaps := bitmap.Load(chatlog.ReadRecords)
		for _, member := range groupUsers.List {
			ids = append(ids, member.UserId)

			if member.UserId == chatlog.SendId {
				continue
			}
			if bitmaps.IsSet(member.UserId) {
				reads = append(reads, member.UserId)
			} else {
				unreads = append(unreads, member.UserId)
			}
		}
	}
	// 根据用户ID获取用户实体信息
	userEntities, err := l.svcCtx.User.FindUser(l.ctx, &user.FindUserReq{Ids: ids})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "get user err 【%v】, req 【%v】", err, req)
	}

	userEntitySet := make(map[string]*user.UserEntity, len(userEntities.Users))

	for i, entity := range userEntities.Users {
		userEntitySet[entity.Id] = userEntities.Users[i]
	}
	// 设置手机号码
	for i, read := range reads {
		if u := userEntitySet[read]; u != nil {
			reads[i] = u.Phone
		}
	}

	for i, unread := range unreads {
		if u := userEntitySet[unread]; u != nil {
			unreads[i] = u.Phone
		}
	}
	// 返回已读和未读用户的信息
	return &types.GetChatLogReadRecordsResp{
		Reads:   reads,
		UnReads: unreads,
	}, nil
}
