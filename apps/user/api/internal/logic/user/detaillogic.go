package user

import (
	"context"
	"github.com/jinzhu/copier"
	"go-chat/apps/user/rpc/user"
	"go-chat/pkg/ctxdata"

	"go-chat/apps/user/api/internal/svc"
	"go-chat/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户信息
func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.UserInfoReq) (resp *types.UserInfoResp, err error) {
	uid := ctxdata.GetId(l.ctx)

	infoResp, err := l.svcCtx.User.GetUserInfo(l.ctx, &user.GetUserInfoReq{
		Id: uid,
	})
	if err != nil {
		return nil, err
	}
	var res types.User
	copier.Copy(&res, infoResp.User)
	return &types.UserInfoResp{
		Info: res,
	}, nil
}
