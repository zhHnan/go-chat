package logic

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"go-chat/apps/user/models"
	"go-chat/apps/user/rpc/internal/svc"
	"go-chat/apps/user/rpc/user"
	"go-chat/pkg/ctxdata"
	"go-chat/pkg/encrypt"
	"go-chat/pkg/wuid"
	"go-chat/pkg/xerr"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrPhoneRegister = xerr.New(xerr.SERVER_COMMON_ERROR, "手机号已注册")
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	// 验证用户是否注册， 根据手机号验证
	userEntity, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
	if err != nil && err != models.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find user by phone error %v, req %v", err, in.Phone)
	}
	if userEntity != nil {
		return nil, errors.WithStack(ErrPhoneRegister)
	}
	// 定义用户数据
	userEntity = &models.Users{
		Id:       wuid.GenUid(l.svcCtx.Config.Mysql.DataSource),
		Avatar:   in.Avatar,
		Nickname: in.Nickname,
		Phone:    in.Phone,
		Sex: sql.NullInt64{
			Int64: int64(in.Sex),
			Valid: true,
		},
	}
	if len(in.Password) > 0 {
		genPwd, err := encrypt.GenPasswordHash([]byte(in.Password))
		if err != nil {
			return nil, errors.Wrapf(xerr.NewDBErr(), "gen password hash error %v", err)
		}
		userEntity.Password = sql.NullString{
			String: string(genPwd),
			Valid:  true,
		}
	}
	_, err = l.svcCtx.UsersModel.Insert(l.ctx, userEntity)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert user error %v", err)
	}
	// 生成token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userEntity.Id)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "generate token error %v", err)
	}
	return &user.RegisterResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
