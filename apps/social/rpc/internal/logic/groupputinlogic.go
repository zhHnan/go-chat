package logic

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"go-chat/apps/social/socialmodels"
	"go-chat/pkg/constants"
	"go-chat/pkg/xerr"
	"time"

	"go-chat/apps/social/rpc/internal/svc"
	"go-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInLogic {
	return &GroupPutInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupPutInLogic) GroupPutIn(in *social.GroupPutInReq) (*social.GroupPutInResp, error) {
	/*
		1. 普通用户申请加入群聊 若无群验证，则直接加入群聊
		2. 群成员邀请 若无群验证，则直接加入群聊
		3. 群管理员、群主邀请 则直接进群
	*/
	var (
		inviteGroupMember *socialmodels.GroupMembers
		userGroupMember   *socialmodels.GroupMembers
		groupInfo         *socialmodels.Groups
		err               error
	)
	userGroupMember, err = l.svcCtx.GroupMembersModel.FindByGroudIdAndUserId(l.ctx, in.GroupId, in.ReqId)
	userGroupMember, err = l.svcCtx.GroupMembersModel.FindByGroudIdAndUserId(l.ctx, in.ReqId, in.GroupId)
	if err != nil && err != socialmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group member by groud id and  req id err %v, req %v, %v", err,
			in.GroupId, in.ReqId)
	}
	if userGroupMember != nil {
		return &social.GroupPutInResp{}, nil
	}
	groupReq, err := l.svcCtx.GroupRequestsModel.FindByGroupIdAndReqId(l.ctx, in.GroupId, in.ReqId)
	if err != nil && err != socialmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group req by groud id and user id err %v, req %v, %v", err,
			in.GroupId, in.ReqId)
	}
	if groupReq != nil {
		return &social.GroupPutInResp{}, nil
	}
	groupReq = &socialmodels.GroupRequests{
		ReqId:   in.ReqId,
		GroupId: in.GroupId,
		ReqMsg: sql.NullString{
			String: in.ReqMsg,
			Valid:  true,
		},
		ReqTime: sql.NullTime{
			Time:  time.Unix(in.ReqTime, 0),
			Valid: true,
		},
		JoinSource: sql.NullInt64{
			Int64: int64(in.JoinSource),
			Valid: true,
		},
		InviterUserId: sql.NullString{
			String: in.InviterUid,
			Valid:  true,
		},
		HandleResult: sql.NullInt64{
			Int64: int64(constants.NoHandlerResult),
			Valid: true,
		},
	}
	createGroupMember := func() {
		if err != nil {
			return
		}
		err = l.createGroupMember(in)
	}
	groupInfo, err = l.svcCtx.GroupsModel.FindOne(l.ctx, in.GroupId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group by groud id err %v, req %v", err, in.GroupId)
	}
	// 验证是否需要处理
	if !groupInfo.IsVerify {
		// 不需要处理 直接加入群聊
		defer createGroupMember()
		groupReq.HandleResult = sql.NullInt64{
			Int64: int64(constants.PassHandlerResult),
			Valid: true,
		}
		return l.createGroupReq(groupReq, true)
	}
	// 验证进群方式
	if constants.GroupJoinSource(in.JoinSource) == constants.PutInGroupJoinSource {
		// 申请
		return l.createGroupReq(groupReq, false)
	}
	inviteGroupMember, err = l.svcCtx.GroupMembersModel.FindByGroudIdAndUserId(l.ctx, in.GroupId, in.ReqId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group member by groud id and user id err %v, req %v",
			in.InviterUid, in.GroupId)
	}

	if constants.GroupRoleLevel(inviteGroupMember.RoleLevel) == constants.CreatorGroupRoleLevel ||
		constants.GroupRoleLevel(inviteGroupMember.RoleLevel) == constants.ManagerGroupRoleLevel {
		// 是管理员或群主
		defer createGroupMember()
		groupReq.HandleResult = sql.NullInt64{
			Int64: int64(constants.PassHandlerResult),
			Valid: true,
		}
		groupReq.HandleUserId = sql.NullString{
			String: in.InviterUid,
			Valid:  true,
		}
		return l.createGroupReq(groupReq, true)
	}
	return l.createGroupReq(groupReq, false)
}
func (l *GroupPutInLogic) createGroupReq(groupReq *socialmodels.GroupRequests, isPass bool) (*social.GroupPutInResp, error) {
	_, err := l.svcCtx.GroupRequestsModel.Insert(l.ctx, groupReq)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert group req err %v req %v", err, groupReq)
	}
	if isPass {
		return &social.GroupPutInResp{GroupId: groupReq.GroupId}, nil
	}
	return &social.GroupPutInResp{}, nil
}
func (l *GroupPutInLogic) createGroupMember(in *social.GroupPutInReq) error {
	groupMember := &socialmodels.GroupMembers{
		GroupId:     in.GroupId,
		UserId:      in.ReqId,
		RoleLevel:   int64(constants.AtLargeGroupRoleLevel),
		OperatorUid: sql.NullString{String: in.InviterUid, Valid: true},
	}
	_, err := l.svcCtx.GroupMembersModel.Insert(l.ctx, nil, groupMember)
	if err != nil {
		return errors.Wrapf(xerr.NewDBErr(), "insert friend err %v req %v", err, groupMember)
	}
	return nil
}
