package logic

import (
	"github.com/zeromicro/go-zero/core/conf"
	"go-chat/apps/user/rpc/internal/config"
	"go-chat/apps/user/rpc/internal/svc"
	"path/filepath"
)

var svcCtx *svc.ServiceContext

func init() {
	var c config.Config
	conf.MustLoad(filepath.Join("../../etc/dev/user.yaml"), &c)
	svcCtx = svc.NewServiceContext(c)
}
