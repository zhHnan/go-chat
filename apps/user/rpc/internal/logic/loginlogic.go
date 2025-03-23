package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-chat/apps/user/models"
	"go-chat/apps/user/rpc/internal/svc"
	"go-chat/apps/user/rpc/user"
	"go-chat/pkg/ctxdata"
	"go-chat/pkg/encrypt"
	"go-chat/pkg/xerr"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrPhoneNotRegister = xerr.New(xerr.SERVER_COMMON_ERROR, "手机号未注册")
	ErrorUserPwdError   = xerr.New(xerr.SERVER_COMMON_ERROR, "密码不正确")
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	// 验证用户是否注册， 根据手机号验证
	userEntity, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
	if err != nil {
		if err == models.ErrNotFound {
			return nil, errors.WithStack(ErrPhoneNotRegister)
		}
		return nil, errors.Wrapf(xerr.NewDBErr(), "find user by phone error %v, req %v", err, in.Phone)
	}
	// 密码验证
	if !encrypt.ValidatePasswordHash(in.Password, userEntity.Password.String) {
		return nil, errors.WithStack(ErrorUserPwdError)
	}
	// 生成token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userEntity.Id)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "generate token error %v", err)
	}
	return &user.LoginResp{
		Id:     userEntity.Id,
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
