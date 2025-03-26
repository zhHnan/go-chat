package middleware

import (
	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"net/http"
)

type LimitMiddleware struct {
	redisCfg redis.RedisConf
	*limit.TokenLimiter
}

func NewLimitMiddleware(redisCfg redis.RedisConf) *LimitMiddleware {
	return &LimitMiddleware{
		redisCfg: redisCfg,
	}
}

// TokenLimitHandler 创建一个限流中间件，用于限制请求速率。
// 参数：
//   - rate: 每秒允许的最大请求数。
//   - burst: 允许的突发请求数。
//   - redisKey: 用于存储限流状态的 Redis Key。
//
// 返回值：
//   - rest.Middleware: 符合 go-zero 框架的中间件函数。
func (m *LimitMiddleware) TokenLimitHandler(rate, burst int, redisKey string) rest.Middleware {
	// 初始化 Redis 客户端，并处理可能的错误
	redisClient, err := redis.NewRedis(m.redisCfg)
	if err != nil {
		// 如果 Redis 初始化失败，记录错误并 panic（或根据需求改为其他处理方式）
		panic("failed to initialize redis client: " + err.Error())
	}

	// 初始化令牌桶限流器
	m.TokenLimiter = limit.NewTokenLimiter(rate, burst, redisClient, redisKey)

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 检查当前请求是否允许通过
			if m.TokenLimiter.AllowCtx(r.Context()) {
				// 请求通过，调用下一个处理器
				next(w, r)
				return
			}

			// 请求被限流，返回 429 Too Many Requests 响应
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("too many requests, please try again later"))
		}
	}
}
