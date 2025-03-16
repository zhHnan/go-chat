package websocket

// 定义服务器配置
type ServerOptions func(opt *serverOptions)

type serverOptions struct {
	patten string
	Authentication
}

func newServerOptions(opts ...ServerOptions) serverOptions {
	so := serverOptions{
		patten:         "/ws",
		Authentication: new(authentication),
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

func WithServerPatten(patten string) ServerOptions {
	return func(opt *serverOptions) {
		opt.patten = patten
	}
}
