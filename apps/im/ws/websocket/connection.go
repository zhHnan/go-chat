package websocket

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Conn struct {
	idleMu sync.Mutex
	Uid    string
	*websocket.Conn
	s *Server
	// 连接空闲时间
	idle time.Time
	// 连接最大空闲时间
	maxConnectionIdle time.Duration
	// 停止信号
	done chan struct{}
}

func NewConn(s *Server, w http.ResponseWriter, r *http.Request) *Conn {
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Errorf("Upgrade ws err %v", err)
		return nil
	}
	conn := &Conn{
		Conn:              c,
		s:                 s,
		idle:              time.Now(),
		maxConnectionIdle: defaultMaxConnectionIdle,
		done:              make(chan struct{}),
	}
	// 启动keepAlive协程来保持连接活跃
	go conn.keepAlive()
	return conn
}

// keepAlive 是一个维护连接活跃状态的方法，它属于 Conn 类型。
// 它通过定时检查连接的闲置时间来决定是否关闭连接，以避免资源泄露。
// 此方法使用一个定时器来监控连接的闲置状态，并在连接闲置时间超过最大闲置时间时关闭连接。
func (c *Conn) keepAlive() {
	// 创建一个定时器，用于检测连接是否闲置过久。
	idleTimer := time.NewTimer(c.maxConnectionIdle)
	// 确保在方法退出时停止定时器，以防止资源泄露。
	defer func() {
		idleTimer.Stop()
	}()

	for {
		select {
		case <-idleTimer.C:
			// 检查连接的闲置时间。
			c.idleMu.Lock()
			idle := c.idle
			if idle.IsZero() {
				// 如果连接从未闲置（即闲置时间戳为零），重置定时器并继续监控。
				c.idleMu.Unlock()
				idleTimer.Reset(c.maxConnectionIdle)
				continue
			}
			// 计算自连接最后一次使用以来的时间，并解锁。
			val := c.maxConnectionIdle - time.Since(idle)
			c.idleMu.Unlock()
			if val <= 0 {
				// 如果连接闲置时间超过了最大闲置时间，关闭连接。
				c.s.Close(c)
				return
			}
			// 重置定时器到连接即将过期的时间。
			idleTimer.Reset(val)
		case <-c.done:
			// 如果连接的 done 通道接收到信号，说明连接应被关闭，退出循环。
			return
		}
	}
}

func (c *Conn) ReadMessage() (messageType int, p []byte, err error) {
	messageType, p, err = c.Conn.ReadMessage()
	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	c.idle = time.Time{}
	return
}

func (c *Conn) WriteMessage(messageType int, data []byte) error {
	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	err := c.Conn.WriteMessage(messageType, data)
	c.idle = time.Now()
	return err
}

func (c *Conn) Close() error {
	select {
	case <-c.done:
	default:
		close(c.done)
	}
	return c.Conn.Close()
}
