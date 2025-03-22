package msgTransfer

import (
	"github.com/zeromicro/go-zero/core/logx"
	"go-chat/apps/im/ws/ws"
	"go-chat/pkg/constants"
	"sync"
	"time"
)

type groupMsgRead struct {
	mu sync.Mutex
	// 记录消息
	push           *ws.Push
	conversationId string
	// 异步推送消息
	pushCh chan *ws.Push

	count int
	// 上次推送时间
	pushTime time.Time
	done     chan struct{}
}

func newGroupMsgRead(push *ws.Push, pushCh chan *ws.Push) *groupMsgRead {
	m := &groupMsgRead{
		conversationId: push.ConversationId,
		push:           push,
		pushCh:         pushCh,
		count:          1,
		done:           make(chan struct{}),
	}
	go m.transfer()
	return m
}

func (m *groupMsgRead) mergePush(push *ws.Push) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.push == nil {
		m.push = push
	}
	m.count++

	for msgId, read := range m.push.ReadRecords {
		m.push.ReadRecords[msgId] = read
	}
}

// transfer 消息推送
func (m *groupMsgRead) transfer() {
	// 超时发送
	// 超量发送

	timer := time.NewTimer(GroupMsgReadRecordDelayTime / 2)
	defer timer.Stop()
	for {
		select {
		case <-m.done:
			return
		case <-timer.C:
			m.mu.Lock()
			pushTime := m.pushTime
			val := GroupMsgReadRecordDelayTime*2 - time.Since(pushTime)
			push := m.push

			if val > 0 && m.count < GroupMsgReadRecordDelayCount || push == nil {
				if val > 0 {
					timer.Reset(val)
				}
				// 未达标
				m.mu.Unlock()
				continue
			}

			m.pushTime = time.Now()
			m.push = nil
			m.count = 0
			timer.Reset(GroupMsgReadRecordDelayTime / 2)
			m.mu.Unlock()
			// 推送消息
			logx.Infof("超过消息合并的条件，推送消息: %s", push)
			m.pushCh <- push
		default:
			m.mu.Lock()
			if m.count >= GroupMsgReadRecordDelayCount {
				push := m.push
				m.push = nil
				m.count = 0
				m.mu.Unlock()

				// 推送消息
				logx.Infof("超过消息合并的条件，推送消息: %s", push)
				m.pushCh <- push
				continue
			}
			if m.IsIdle() {
				m.mu.Unlock()
				// 释放资源
				m.pushCh <- &ws.Push{
					ConversationId: m.conversationId,
					ChatType:       constants.GroupChatType,
				}
				continue
			}
			m.mu.Unlock()

			temDelay := GroupMsgReadRecordDelayTime / 4
			if temDelay > time.Second {
				temDelay = time.Second
			}
			time.Sleep(temDelay)
		}
	}
}

// 判断消息是否活跃
func (m *groupMsgRead) IsIdle() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.isIdle()
}

func (m *groupMsgRead) isIdle() bool {
	pushTime := m.pushTime
	val := GroupMsgReadRecordDelayTime*2 - time.Since(pushTime)
	if val <= 0 && m.push == nil && m.count == 0 {
		return true
	}
	return false
}

// clear 清除消息
func (m *groupMsgRead) clear() {
	select {
	case <-m.done:
	default:
		close(m.done)
	}
	m.push = nil
}
