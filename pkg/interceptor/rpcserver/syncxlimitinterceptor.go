package rpcserver

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/syncx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func SyncxLimitInterceptor(maxCount int) grpc.UnaryServerInterceptor {
	l := syncx.NewLimit(maxCount)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{},
		err error) {
		if l.TryBorrow() {
			defer func() {
				if err := l.Return(); err != nil {
					logx.WithContext(ctx).Errorf("【RPC SERVER ERROR】syncxLimitInterceptor err:%v", err)
				}
			}()
			return handler(ctx, req)
		} else {
			logx.Errorf("【RPC SERVER ERROR】syncxLimitInterceptor err:【%v】:too many requests", maxCount)
			return nil, status.Errorf(codes.ResourceExhausted, "too many requests")
		}
	}
}
