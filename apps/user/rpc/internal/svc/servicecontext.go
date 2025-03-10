package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-chat/apps/user/models"
	"go-chat/apps/user/rpc/internal/config"
)

// ServiceContext 注册服务上下文
type ServiceContext struct {
	Config config.Config
	models.UsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// sql 连接
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config:     c,
		UsersModel: models.NewUsersModel(sqlConn, c.Cache),
	}
}
