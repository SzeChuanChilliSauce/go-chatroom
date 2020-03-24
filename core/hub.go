package core

import "encoding/json"

// 全局控制器对象
var ChatHub = Hub{
	Conns:       make(map[*Conn]bool),
	BroadcastCh: make(chan []byte),
	RegisterCh:  make(chan *Conn),
	CancelCh:    make(chan *Conn),
}

// Hub 消息控制中心,处理websocket中各种消息
type Hub struct {
	// 已注册的连接
	Conns map[*Conn]bool
	// 广播消息管道
	BroadcastCh chan []byte
	// 注册管道
	RegisterCh chan *Conn
	// 注销管道
	CancelCh chan *Conn
}

// 处理ws
func (h *Hub) Run() {
	// 监听数据通道，在后端处理通道数据
	for {
		select {
		case conn := <-h.RegisterCh:
			// 标识注册
			h.Conns[conn] = true
			// 组装data数据
			conn.Data.IP = conn.WS.RemoteAddr().String()
			// 更新类型
			conn.Data.Type = Handshake
			// 更新用户列表
			conn.Data.UserList = UserList
			bin, _ := json.Marshal(conn.Data)
			// 将数据放入通道
			conn.SendCh <- bin
		case conn := <-h.CancelCh:
			if _, ok := h.Conns[conn]; ok {
				delete(h.Conns, conn)
				close(conn.SendCh)
			}
		case data := <-h.BroadcastCh:
			// 处理数据
			// 将数据广播到所有用户
			for conn := range h.Conns {
				select {
				case conn.SendCh <- data:
				default:
					delete(h.Conns, conn)
					close(conn.SendCh)
				}
			}
		}
	}
}
