package core

import "encoding/json"

// h 全局对象
var Hub = hub{
	connections: make(map[*connection]bool),
	broadcast:   make(chan []byte),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
}

// 处理ws
func (h *hub) run() {
	// 监听数据通道，在后端处理通道数据
	for {
		select {
		case conn := <-h.register:
			// 标识注册
			h.connections[conn] = true
			// 组装data数据
			conn.data.IP = conn.ws.RemoteAddr().String()
			// 更新类型
			conn.data.Type = "handshake"
			// 更新用户列表
			conn.data.UserList = userList
			bin, _ := json.Marshal(conn.data)
			// 将数据放入通道
			conn.send <- bin
		case conn := <-h.unregister:
			if _, ok := h.connections[conn]; ok {
				delete(h.connections, conn)
				close(conn.send)
			}
		case data := <-h.broadcast:
			// 处理数据
			// 将数据广播到所有用户
			for conn := range h.connections {
				select {
				case conn.send <- data:
				default:
					delete(h.connections, conn)
					close(conn.send)
				}
			}
		}
	}
}
