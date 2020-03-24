package core

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

// 用户列表
var UserList = []string{}

// 连接
type Conn struct {
	// websocket连接
	WS *websocket.Conn
	// 发送管道
	SendCh chan []byte
	// 数据
	Data *Message
}

// write 写数据
func (conn *Conn) Write() {
	// 从管道取数据
	for msg := range conn.SendCh {
		_ = conn.WS.WriteMessage(websocket.TextMessage, msg)
	}

	_ = conn.WS.Close()
}

// read 读数据
func (conn *Conn) Read() {
	for {
		// 读数据
		_, msg, err := conn.WS.ReadMessage()
		if err != nil {
			// 移除
			ChatHub.CancelCh <- conn
			break
		}

		err = json.Unmarshal(msg, conn.Data)
		if err != nil {
			break
		}

		switch conn.Data.Type {
		case Login:
			// 弹出窗口，输入用户名
			conn.Data.User = conn.Data.Content
			conn.Data.From = conn.Data.User
			// 登录后，将用户加入到列表
			UserList = append(UserList, conn.Data.User)
			// 每个用户都加载所有登录的用户列表
			conn.Data.UserList = UserList
			// 数据序列化
			bin, _ := json.Marshal(conn.Data)
			ChatHub.BroadcastCh <- bin
		case User:
			conn.Data.Type = User
			bin, _ := json.Marshal(conn.Data)
			ChatHub.BroadcastCh <- bin
		case Logout:
			conn.Data.Type = Logout
			// 用户列表删除
			UserList = Remove(UserList, conn.Data.User)
			conn.Data.UserList = UserList
			conn.Data.Content = conn.Data.User
			// 序列化
			bin, _ := json.Marshal(conn.Data)
			ChatHub.BroadcastCh <- bin
			ChatHub.CancelCh <- conn
		default:
			fmt.Println("other...")
		}
	}
}
