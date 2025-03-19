package svc

import (
	"go-chat/apps/user/models"
	"go-chat/apps/user/rpc/internal/config"
	"go-chat/pkg/constants"
	"go-chat/pkg/ctxdata"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// ServiceContext 注册服务上下文
type ServiceContext struct {
	Config config.Config
	*redis.Redis
	models.UsersModel
	logx.Logger
}

func NewServiceContext(c config.Config) *ServiceContext {
	// sql 连接
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)

	ctx := &ServiceContext{
		Config:     c,
		Redis:      redis.MustNewRedis(c.Redisx),
		UsersModel: models.NewUsersModel(sqlConn, c.Cache),
	}
	return ctx
}

func (svc *ServiceContext) SetRootToken() error {
	// 生成jwt
	token, err := ctxdata.GetJwtToken(svc.Config.Jwt.AccessSecret, time.Now().Unix(), 999999999, constants.SYSTEM_ROOT_UID)
	if err != nil {
		return err
	}
	return svc.Redis.Set(constants.REDIS_SYSTEM_ROOT_TOKEN, token)
}
