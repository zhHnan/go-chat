package job

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

// 定义重试的时间策略
// RetryJetLagFunc 定义了重试时间间隔的策略函数类型。
// 参数:
// - ctx: 上下文，用于控制超时或取消操作。
// - retryCount: 当前已重试的次数。
// - lastTime: 上一次重试的时间间隔。
// 返回值: 下一次重试的时间间隔。
type RetryJetLagFunc func(ctx context.Context, retryCount int, lastTime time.Duration) time.Duration

// RetryJetLagAlways 是一个默认的重试时间间隔策略函数。
// 无论重试次数或上次时间间隔如何，始终返回固定的默认重试间隔 DefaultRetryJetLag。
// 参数:
// - ctx: 上下文，用于控制超时或取消操作。
// - retryCount: 当前已重试的次数。
// - lastTime: 上一次重试的时间间隔。
// 返回值: 固定的默认重试间隔 DefaultRetryJetLag。
func RetryJetLagAlways(ctx context.Context, retryCount int, lastTime time.Duration) time.Duration {
	return DefaultRetryJetLag
}

// IsRetryFunc 定义了一个判断是否需要重试的函数类型。
// 参数:
// - ctx: 上下文，用于控制超时或取消操作。
// - retryCount: 当前已重试的次数。
// - err: 上一次执行时发生的错误。
// 返回值: 是否需要继续重试。
type IsRetryFunc func(ctx context.Context, retryCount int, err error) bool

// RetryAlways 是一个始终返回 true 的重试判断函数。
// 无论错误内容或重试次数如何，都会继续重试。
// 参数:
// - ctx: 上下文，用于控制超时或取消操作。
// - retryCount: 当前已重试的次数。
// - err: 上一次执行时发生的错误。
// 返回值: 始终返回 true，表示需要重试。
func RetryAlways(ctx context.Context, retryCount int, err error) bool {
	return true
}

// WithRetry 提供了一个带有重试逻辑的函数执行机制。
// 根据提供的上下文、处理器函数和重试选项，尝试多次执行处理器函数。
// 参数:
// - ctx: 上下文，用于控制超时或取消操作。
// - handler: 需要执行的目标函数，接收上下文并返回错误。
// - opts: 可选的重试配置参数。
// 返回值: 如果成功执行则返回 nil，否则返回最后一次执行的错误。
func WithRetry(ctx context.Context, handler func(ctx context.Context) error, opts ...RetryOptions) error {
	opt := newOptions(opts...)
	// 检查上下文中是否设置了截止时间，如果没有，则添加一个超时限制以防止无限重试。
	_, ok := ctx.Deadline()
	if !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opt.timeout)
		defer cancel()
	}

	var (
		handlerErr  error                 // 保存处理器函数的执行结果。
		retryJetLag time.Duration         // 当前的重试时间间隔。
		ch          = make(chan error, 1) // 用于异步接收处理器函数的执行结果。
	)

	// 按照配置的重试次数进行循环重试。
	for i := 0; i < opt.retryCount; i++ {
		go func() {
			ch <- handler(ctx) // 异步执行处理器函数并将结果发送到通道。
		}()

		select {
		case handlerErr = <-ch: // 接收到处理器函数的执行结果。
			if handlerErr == nil { // 如果执行成功，则直接返回。
				return nil
			}
			// 判断是否需要继续重试。
			if !opt.isRetryFunc(ctx, i, handlerErr) {
				return handlerErr
			}
			// 计算下一次重试的时间间隔并休眠。
			retryJetLag := opt.retryJetLag(ctx, i, retryJetLag)
			time.Sleep(retryJetLag)
		case <-ctx.Done(): // 如果上下文被取消或超时，则返回错误。
			return errors.New("retry timeout")
		}
	}
	// 返回最后一次执行的错误。
	return handlerErr
}
