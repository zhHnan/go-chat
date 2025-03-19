package handler

import (
	"context"
	"go-chat/apps/im/ws/internal/svc"
	"go-chat/pkg/ctxdata"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/token"
)

type JwtAuth struct {
	svc    *svc.ServiceContext
	parser *token.TokenParser
	logx.Logger
}

func (j *JwtAuth) Auth(w http.ResponseWriter, r *http.Request) bool {
	j.Infof("开始处理认证请求: %s", r.URL.String())
	tok, err := j.parser.ParseToken(r, j.svc.Config.JwtAuth.AccessSecret, "")
	if err != nil {
		j.Errorf("解析token失败: %v", err)
		return false
	}
	if !tok.Valid {
		j.Errorf("token无效")
		return false
	}
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		j.Errorf("token claims不是map类型")
		return false
	}

	// 添加详细日志
	j.Infof("token解析成功, claims: %v", claims)
	j.Infof("找到用户ID: %v", claims[ctxdata.Identify])

	*r = *r.WithContext(context.WithValue(r.Context(), ctxdata.Identify, claims[ctxdata.Identify]))
	return true
}

func (j *JwtAuth) UserId(r *http.Request) string {
	userId := ctxdata.GetId(r.Context())
	j.Infof("获取用户ID: %s", userId)
	return userId
}

func NewJwtAuth(svc *svc.ServiceContext) *JwtAuth {
	return &JwtAuth{
		svc:    svc,
		parser: token.NewTokenParser(),
		Logger: logx.WithContext(context.Background()),
	}
}
