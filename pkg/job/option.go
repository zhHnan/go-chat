package job

import "time"

// RetryOptions定义了重试策略的配置接口。
// 它是一个函数类型，接收一个指向retryOptions的指针，用于修改重试策略的属性。
type (
	RetryOptions func(opts *retryOptions)

	// retryOptions包含了重试策略的具体配置。
	retryOptions struct {
		timeout     time.Duration   // timeout定义了重试的总时间限制。
		retryCount  int             // retryCount定义了最大重试次数。
		isRetryFunc IsRetryFunc     // isRetryFunc定义了是否需要重试的判断逻辑。
		retryJetLag RetryJetLagFunc // retryJetLag定义了重试间隔的计算逻辑。
	}
)

// newOptions创建并返回一个带有默认重试策略的retryOptions实例。
// 它接受一个或多个RetryOptions作为参数，这些参数用于定制重试策略。
// 参数opts是可变参数，每个参数都是一个配置函数，用于应用特定的配置。
func newOptions(opts ...RetryOptions) *retryOptions {
	// 初始化一个带有默认值的retryOptions实例。
	opt := &retryOptions{
		timeout:     DefaultRetryTimeout,  // 使用默认的重试超时时间。
		retryCount:  DefaultRetryMaxCount, // 使用默认的最大重试次数。
		isRetryFunc: RetryAlways,          // 默认情况下总是进行重试。
		retryJetLag: RetryJetLagAlways,    // 默认情况下立即进行重试。
	}
	// 遍历所有提供的配置函数，将它们应用于retryOptions实例。
	for _, options := range opts {
		options(opt)
	}
	// 返回配置完成的retryOptions实例。
	return opt
}

func WithTimeout(timeout time.Duration) RetryOptions {
	return func(opts *retryOptions) {
		if timeout > 0 {
			opts.timeout = timeout
		}
	}
}
