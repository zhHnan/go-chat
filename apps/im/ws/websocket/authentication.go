package websocket

import (
	"fmt"
	"net/http"
	"time"
)

// 鉴权

type Authentication interface {
	// 鉴权
	Auth(w http.ResponseWriter, r *http.Request) bool
	UserId(r *http.Request) string
}
type authentication struct {
}

func (authentication) Auth(w http.ResponseWriter, r *http.Request) bool {
	return true
}

func (authentication) UserId(r *http.Request) string {
	query := r.URL.Query()
	if query != nil && query["userId"] != nil {
		return fmt.Sprintf("%s", query["userId"])
	}
	return fmt.Sprintf("%s", time.Now().Unix())
}
