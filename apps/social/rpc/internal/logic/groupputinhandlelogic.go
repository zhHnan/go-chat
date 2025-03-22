package logic

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-chat/apps/social/rpc/internal/svc"
	"go-chat/apps/social/rpc/social"
	"go-chat/apps/social/socialmodels"
	"go-chat/pkg/constants"
	"go-chat/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

var (
	ErrGroupReqBeforePass   = xerr.NewMsg("请求已通过")
	ErrGroupReqBeforeRefuse = xerr.NewMsg("请求已拒绝")
)

func NewGroupPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInHandleLogic {
	return &GroupPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupPutInHandleLogic) GroupPutInHandle(in *social.GroupPutInHandleReq) (*social.GroupPutInHandleResp, error) {
	one, err := l.svcCtx.GroupRequestsModel.FindOne(l.ctx, uint64(in.GroupReqId))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find groupRequest by groupReqId err %v req %v ", err, in.GroupReqId)
	}
	switch constants.HandlerResult(one.HandleResult.Int64) {
	case constants.PassHandlerResult:
		return nil, errors.WithStack(ErrGroupReqBeforePass)
	case constants.RefuseHandlerResult:
		return nil, errors.WithStack(ErrGroupReqBeforeRefuse)
	}
	one.HandleResult = sql.NullInt64{
		Int64: int64(in.HandleResult),
		Valid: true,
	}
	l.svcCtx.GroupRequestsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		if err := l.svcCtx.GroupRequestsModel.Update(l.ctx, session, one); err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "update group request err %v, req %v", err, one)
		}
		if constants.HandlerResult(one.HandleResult.Int64) != constants.PassHandlerResult {
			return nil
		}
		groupMember := &socialmodels.GroupMembers{
			GroupId:     one.GroupId,
			UserId:      one.ReqId,
			RoleLevel:   int64(constants.AtLargeGroupRoleLevel),
			OperatorUid: sql.NullString{String: in.HandleUid, Valid: true},
		}
		_, err = l.svcCtx.GroupMembersModel.Insert(l.ctx, session, groupMember)
		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "insert friend err %v req %v", err, groupMember)
		}

		return nil
	})

	if constants.HandlerResult(one.HandleResult.Int64) != constants.PassHandlerResult {
		return &social.GroupPutInHandleResp{}, nil
	}
	return &social.GroupPutInHandleResp{
		GroupId: one.GroupId,
	}, nil
}
