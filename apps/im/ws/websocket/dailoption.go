package websocket

import "net/http"

// DailOptions 创建一个DailOption
type DailOptions func(option *DailOption)
type DailOption struct {
	header http.Header
	patten string
}

// NewDialOption 创建一个DailOption
func NewDialOption(opts ...DailOptions) DailOption {
	option := DailOption{
		patten: "/ws",
		header: nil,
	}

	for _, opt := range opts {
		opt(&option)
	}
	return option
}

// WithClientPatten 设置客户端的patten
func WithClientHeader(header http.Header) DailOptions {
	return func(option *DailOption) {
		option.header = header
	}
}

// WithClientPatten 设置客户端的patten
func WithClientPatten(patten string) DailOptions {
	return func(option *DailOption) {
		option.patten = patten
	}
}
