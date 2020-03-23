package main

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

// ws连接器 数据 管道

// connection ...
type connection struct {
	// ws连接器
	ws *websocket.Conn

	// 管道
	send chan []byte

	// 数据
	data *Data
}

// hub 处理ws中的各种逻辑
type hub struct {
	// 已注册的连接器
	connections map[*connection]bool

	// 广播消息
	broadcast chan []byte

	// 注册
	register chan *connection

	// 注销
	unregister chan *connection
}

// write 写数据
func (conn *connection) write() {
	// 从管道取数据
	for msg := range c.send {
		conn.ws.WriteMessage(websocket.TextMessage, msg)
	}

	conn.ws.Close()
}

// read 读数据
func (conn *connection) read() {
	for {
		// 读数据
		_, msg, err := conn.ws.ReadMessage()
		if err != nil {
			// 移除
			h.unregister <- conn
			break
		}

		// 反序列化
		err = json.Unmarshal(msg, conn.data)
		if err != nil {
			break
		}

		switch conn.data.Type {
		case "login":
			conn.data.User = conn.data.Content
			conn.data.From = conn.data.User

		}
	}
}
