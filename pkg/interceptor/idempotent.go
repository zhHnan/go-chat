package interceptor

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/utils"
	"go-chat/pkg/xerr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

type Idempotent interface {
	// 幂等标识
	Identify(ctx context.Context, method string) string
	// 判断是否是幂等方法
	IsIdempotentMethod(fullMethod string) bool
	// 幂等性的验证
	TryAcquire(ctx context.Context, id string) (resp interface{}, isAcquire bool)
	// 执行之后结果的保存
	SaveResp(ctx context.Context, id string, resp interface{}, respErr error) error
}

var (
	// 请求任务标识
	TKey = "hnz-chat-idempotent-task-id"
	// rpc调度中rpc请求的标识
	DKey = "hnz-chat-idempotent-dispatch-id"
)

func ContextWithVal(ctx context.Context) context.Context {
	// 设置请求的id
	return context.WithValue(ctx, TKey, utils.NewUuid())
}

// NewIdempotent 定义客户端的拦截器
func NewIdempotentClient(idempotent Idempotent) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any,
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// 获取唯一的key
		identify := idempotent.Identify(ctx, method)
		metadata.NewOutgoingContext(ctx, map[string][]string{
			DKey: {identify},
		})

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func NewIdempotentServer(idempotent Idempotent) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		// 获取请求的id
		identify := metadata.ValueFromIncomingContext(ctx, DKey)
		if len(identify) == 0 || idempotent.IsIdempotentMethod(info.FullMethod) {
			// 不进行幂等性验证
			return handler(ctx, req)
		}
		fmt.Println("幂等性验证...", identify)
		r, isAcquire := idempotent.TryAcquire(ctx, identify[0])
		if isAcquire {
			resp, err = handler(ctx, req)
			fmt.Println("执行任务...", identify)
			if err := idempotent.SaveResp(ctx, identify[0], resp, err); err != nil {
				return resp, err
			}
			return resp, err
		}
		// 任务已经执行过了
		fmt.Println("任务在执行...", identify)
		if r != nil {
			fmt.Println("任务执行完毕...", identify)
			return r, nil
		}
		// 任务执行失败
		return r, errors.WithStack(xerr.New(int(codes.DeadlineExceeded), fmt.Sprintf("存在其他任务在执行: %s", identify[0])))
	}
}

var (
	DefaultIdempotent       = new(defaultIdempotent)
	DefaultIdempotentClient = NewIdempotentClient(DefaultIdempotent)
)

type defaultIdempotent struct {
	// 获取和设置请求的id
	*redis.Redis
	// 注意存储
	*collection.Cache
	// 设置方法对幂等性的支持
	method map[string]bool
}

func NewDefaultIdempotent(c redis.RedisConf) Idempotent {
	cache, err := collection.NewCache(60 * 60)
	if err != nil {
		panic(err)
	}

	return &defaultIdempotent{
		Redis: redis.MustNewRedis(c),
		Cache: cache,
		method: map[string]bool{
			"/social.social/GroupCreate": true,
		},
	}
}

// Identify 幂等标识
func (d *defaultIdempotent) Identify(ctx context.Context, method string) string {
	// 获取请求的id
	id := ctx.Value(TKey)
	rpcId := fmt.Sprintf("%v.%v", id, method)
	return rpcId
}

// IsIdempotentMethod `
func (d *defaultIdempotent) IsIdempotentMethod(fullMethod string) bool {
	return d.method[fullMethod]
}

// TryAcquire 幂等性的验证
func (d *defaultIdempotent) TryAcquire(ctx context.Context, id string) (resp interface{}, isAcquire bool) {
	// 基于redis实现幂等性
	retry, err := d.SetnxEx(id, "1", 60*60)
	if err != nil {
		return nil, false
	}
	if retry {
		return nil, true
	}
	resp, _ = d.Cache.Get(id)
	return resp, false
}

// SaveResp 执行之后结果的保存
func (d *defaultIdempotent) SaveResp(ctx context.Context, id string, resp interface{}, respErr error) error {
	d.Cache.Set(id, resp)
	return nil
}
