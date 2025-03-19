package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/url"
)

type Client interface {
	Close() error
	// Send 发送消息
	Send(v any) error
	// Read 接收消息
	Read(v any) error
}

type client struct {
	*websocket.Conn
	host string
	opt  DailOption
}

func NewClient(host string, opts ...DailOptions) *client {
	opt := NewDialOption(opts...)
	c := client{
		host: host,
		opt:  opt,
		Conn: nil,
	}
	conn, err := c.dail()
	if err != nil {
		panic(err)
	}
	c.Conn = conn
	return &c
}

// dail 建立与websocket的连接
func (c *client) dail() (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: c.host, Path: c.opt.patten}
	dial, _, err := websocket.DefaultDialer.Dial(u.String(), c.opt.header)
	return dial, err
}

func (c *client) Send(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	// 如果连接丢失，尝试重新连接
	if c.Conn == nil {
		conn, err := c.dail()
		if err != nil {
			return err
		}
		c.Conn = conn
	}

	err = c.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		// 发生错误，可能是连接断开，尝试重新连接
		conn, dialErr := c.dail()
		if dialErr != nil {
			return dialErr
		}
		c.Conn = conn
		// 使用新连接重试一次
		return c.WriteMessage(websocket.TextMessage, data)
	}
	return nil
}

func (c *client) Read(v any) error {
	_, message, err := c.Conn.ReadMessage()
	if err != nil {
		return err
	}
	return json.Unmarshal(message, v)
}
