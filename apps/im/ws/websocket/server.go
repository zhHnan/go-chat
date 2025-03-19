package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type Server struct {
	// 加锁
	sync.RWMutex
	routes         map[string]HandlerFunc
	addr           string
	patten         string
	authentication Authentication
	// 容易产生并发
	connToUser map[*Conn]string
	userToConn map[string]*Conn
	upgrader   websocket.Upgrader
	logx.Logger
}

func (s *Server) AddRoute(rs []Route) {
	for _, r := range rs {
		s.routes[r.Method] = r.Handler
	}
}
func NewServer(addr string, opts ...ServerOptions) *Server {
	opt := newServerOptions(opts...)
	return &Server{
		routes:         make(map[string]HandlerFunc),
		addr:           addr,
		patten:         opt.patten,
		authentication: opt.Authentication,
		upgrader:       websocket.Upgrader{},
		Logger:         logx.WithContext(context.Background()),
		connToUser:     make(map[*Conn]string),
		userToConn:     make(map[string]*Conn),
	}
}
func (s *Server) addConn(conn *Conn, req *http.Request) {
	uid := s.authentication.UserId(req)
	s.Infof("用户尝试连接: %s", uid)
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()
	if c := s.userToConn[uid]; c != nil {
		// 关闭之前的连接
		s.Infof("关闭用户之前的连接: %s", uid)
		c.Close()
	}
	s.connToUser[conn] = uid
	s.userToConn[uid] = conn
	s.Infof("用户连接成功: %s, 当前连接数量: %d", uid, len(s.connToUser))
}

func (s *Server) GetConn(uid string) *Conn {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()
	conn := s.userToConn[uid]
	if conn == nil {
		s.Infof("获取连接失败，用户 %s 不存在", uid)
		// 打印当前所有在线用户，帮助调试
		s.Infof("当前在线用户: %v", s.GetUsers())
	}
	return conn
}
func (s *Server) ServerWs(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			s.Errorf("Server handler ws recover err %v", r)
		}
	}()
	// 记录连接尝试
	s.Infof("收到WebSocket连接请求: %s", r.URL.String())

	// 升级协议，获取连接对象
	conn := NewConn(s, w, r)
	if conn == nil {
		s.Errorf("创建连接对象失败")
		return
	}
	//conn, err := s.upgrader.Upgrade(w, r, nil)
	//if err != nil {
	//	s.Errorf("Upgrade ws err %v", err)
	//	return
	//}
	if !s.authentication.Auth(w, r) {
		//conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("forbidden visit")))
		s.Errorf("认证失败，拒绝连接")
		s.Send(&Message{FrameType: FrameData, Data: fmt.Sprint("forbidden visit")}, conn)
		conn.Close()
		return
	}

	s.Infof("认证成功，添加连接")
	// 记录连接
	s.addConn(conn, r)
	// 根据连接对象创建客户端对象
	// 使用协程 处理连接
	go s.handlerConn(conn)
}

// 根据连接对象执行任务处理
func (s *Server) handlerConn(conn *Conn) {
	s.Infof("开始处理连接...")
	// 此方法用于检索与特定连接相关联的所有用户ID
	uids := s.GetUsers(conn)
	s.Infof("连接关联的用户ID: %v", uids)
	// 设置连接的用户ID为用户ID列表中的第一个ID
	// 这里假设uids数组非空，且conn对象具有uid属性
	if len(uids) == 0 {
		s.Errorf("未找到连接对应的用户ID，关闭连接")
		s.Close(conn)
		return
	}

	conn.Uid = uids[0]
	s.Infof("设置连接用户ID: %s", conn.Uid)

	for {
		s.Infof("等待消息...")
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("websocket connection read message err %v", err)
			s.Close(conn)
			return
		}
		// 解析消息
		var message Message
		if err = json.Unmarshal(msg, &message); err != nil {
			s.Errorf("json unmarshal err %v, msg: %v", err, string(msg))
			s.Close(conn)
			return
		}
		// 依据消息进行处理
		switch message.FrameType {
		case FramePing:
			s.Send(&Message{FrameType: FramePing}, conn)
		case FrameData:
			// 根据请求的方法分发路由
			if handler, ok := s.routes[message.Method]; ok {
				// 执行方法
				handler(s, conn, &message)
			} else {
				s.Send(&Message{FrameType: FrameData, Data: fmt.Sprintf("method %s not found", message.Method)}, conn)
				//conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("method %s not found", message.Method)))
			}
		}
	}
}

func (s *Server) GetUsers(conns ...*Conn) []string {

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	var res []string
	if len(conns) == 0 {
		// 获取全部
		res = make([]string, 0, len(s.connToUser))
		for _, uid := range s.connToUser {
			res = append(res, uid)
		}
	} else {
		// 获取部分
		res = make([]string, 0, len(conns))
		for _, conn := range conns {
			res = append(res, s.connToUser[conn])
		}
	}

	return res
}

func (s *Server) SendByUserId(msg interface{}, sendIds ...string) error {
	if len(sendIds) == 0 {
		return nil
	}
	return s.Send(msg, s.GetConns(sendIds...)...)
}
func (s *Server) GetConns(uids ...string) []*Conn {
	if len(uids) == 0 {
		return nil
	}

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	res := make([]*Conn, 0, len(uids))
	for _, uid := range uids {
		res = append(res, s.userToConn[uid])
	}
	return res
}
func (s *Server) Send(msg interface{}, conns ...*Conn) error {
	if len(conns) == 0 {
		return nil
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return err
		}
	}
	return nil
}
func (s *Server) Start() {
	http.HandleFunc(s.patten, s.ServerWs)
	fmt.Printf("starting server at %s", s.addr)
	s.Info(http.ListenAndServe(s.addr, nil))
}

func (s *Server) Stop() {

	fmt.Println("stop server")
}

func (s *Server) Close(conn *Conn) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	uid := s.connToUser[conn]
	if uid == "" {
		// 已经被关闭
		return
	}

	delete(s.connToUser, conn)
	delete(s.userToConn, uid)
	conn.Close()
}
