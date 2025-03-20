package websocket

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type AckType int

const (
	NoAck AckType = iota
	OnlyAck
	RigorAck
)

func (t AckType) ToString() string {
	switch t {
	case OnlyAck:
		return "OnlyAck"
	case RigorAck:
		return "RigorA"
	default:
		return "NoAck"
	}
}

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
	// 并发控制
	messageMu sync.Mutex
	// 读取消息
	readMessage []*Message
	// 读取消息序列化
	readMessageSeq map[string]*Message
	// 用于在readAck 与 handlerWrite 之间传递消息
	message chan *Message
}

func NewConn(s *Server, w http.ResponseWriter, r *http.Request) *Conn {
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Errorf("Upgrade ws err %v", err)
		return nil
	}
	conn := &Conn{
		Conn: c,
		s:    s,
		idle: time.Now(),
		// readMessage 是一个用于存储已读取消息的切片，初始化时容量为2，意味着预期少量消息即可被读取。
		readMessage: make([]*Message, 0, 2),
		// readMessageSeq 是一个映射，用于根据消息序列号快速检索已读取的消息，预设容量为2，优化初始查找性能。
		readMessageSeq: make(map[string]*Message, 2),
		// message 是一个带缓冲的通道，用于在goroutine之间传递Message对象，缓冲区大小为1，确保一定程度的消息流通不会阻塞发送者。
		message: make(chan *Message, 1),
		// maxConnectionIdle 指定了最大空闲连接数，被设置为默认值defaultMaxConnectionIdle，用于控制空闲连接的数量，避免资源浪费。
		maxConnectionIdle: defaultMaxConnectionIdle,
		done:              make(chan struct{}),
	}
	// 启动keepAlive协程来保持连接活跃
	go conn.keepAlive()
	return conn
}
func (c *Conn) appendMsgMq(msg *Message) {
	c.messageMu.Lock()
	defer c.messageMu.Unlock()
	if m, ok := c.readMessageSeq[msg.Id]; ok {
		if len(c.readMessage) == 0 {
			// 如果切片为空，则直接返回，因为消息序列号已经存在但切片为空，无法找到对应的消息。
			return
		}
		// msg.Ack > m.Ack
		if msg.AckSeq <= m.AckSeq {
			// 如果消息序列号已经存在，但消息的AckSeq小于等于已存在的消息的AckSeq，则返回，
			// 因为新消息的AckSeq应该大于已存在的消息的AckSeq。
			return
		}
		// 更新映射中的消息
		c.readMessageSeq[msg.Id] = msg
		return
	}
	if msg.FrameType == FrameAck {
		// 如果消息类型为FrameAck，则从切片和映射中删除该消息。
		return
	}
	// 如果消息序列号不存在，则将消息添加到切片和映射中。
	c.readMessage = append(c.readMessage, msg)
	c.readMessageSeq[msg.Id] = msg
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
