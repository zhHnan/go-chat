package websocket

import "time"

// 定义服务器配置
type ServerOptions func(opt *serverOptions)

type serverOptions struct {
	patten     string
	ack        AckType
	ackTimeout time.Duration
	Authentication
	maxConnectionIdle time.Duration
}

func newServerOptions(opts ...ServerOptions) serverOptions {
	so := serverOptions{
		patten:            "/ws",
		maxConnectionIdle: defaultMaxConnectionIdle,
		ackTimeout:        defaultAckTimeout,
		Authentication:    new(authentication),
	}
	for _, opt := range opts {
		opt(&so)
	}
	return so
}

func WithServerAuthentication(auth Authentication) ServerOptions {
	return func(opt *serverOptions) {
		opt.Authentication = auth
	}
}

func WithServerAck(ack AckType) ServerOptions {
	return func(opt *serverOptions) {
		opt.ack = ack
	}
}

func WithServerPatten(patten string) ServerOptions {
	return func(opt *serverOptions) {
		opt.patten = patten
	}
}

func WithServerMaxConnectionIdle(maxConnectionIdle time.Duration) ServerOptions {
	return func(opt *serverOptions) {
		if maxConnectionIdle > 0 {
			opt.maxConnectionIdle = maxConnectionIdle
		}
	}
}
