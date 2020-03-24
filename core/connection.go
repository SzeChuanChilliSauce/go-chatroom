package core

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	for msg := range conn.send {
		conn.ws.WriteMessage(websocket.TextMessage, msg)
	}

	conn.ws.Close()
}

// 用户列表
var userList = []string{}

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
			// 弹出窗口，输入用户名
			conn.data.User = conn.data.Content
			conn.data.From = conn.data.User
			// 登录后，将用户加入到列表
			userList = append(userList, conn.data.User)
			// 每个用户都加载所有登录的用户列表
			conn.data.UserList = userList
			// 数据序列化
			bin, _ := json.Marshal(conn.data)
			h.broadcast <- bin
		case "user":
			conn.data.Type = "user"
			bin, _ := json.Marshal(conn.data)
			h.broadcast <- bin
		case "logout":
			conn.data.Type = "logout"
			// 用户列表删除
			userList = remove(userList, conn.data.User)
			conn.data.UserList = userList
			conn.data.Content = conn.data.User
			// 序列化
			bin, _ := json.Marshal(conn.data)
			h.broadcast <- bin
			h.unregister <- conn
		default:
			fmt.Println("other...")
		}
	}
}

// 从列表中删除元素
func remove(slice []string, user string) []string {
	l := len(slice)
	if l == 0 {
		return slice
	}

	var res = []string{}

	// 是第一个元素
	if slice[0] == user {
		return slice[1:]
	}

	// 是最后一个元素
	if slice[l-1] == user {
		return slice[:l-1]
	}

	for i := range slice {
		// 是中间的元素
		if slice[i] == user {
			res = append(res, slice[:i]...)
			res = append(res, slice[i+1:]...)
			return res
		}
	}

	// 没有要删除的元素
	return slice
}

// 定义升级器，将http请求升级为ws请求
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 跨域访问
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ws的回调函数
func wsHandler(w http.ResponseWriter, r *http.Request) {
	// 升级
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// 初始化连接对象
	conn := &connection{
		send: make(chan []byte, 128),
		ws:   ws,
		data: &Data{},
	}

	// 在ws中注册
	h.register <- conn

	go conn.write()
	conn.read()

	defer func() {
		conn.data.Type = "logout"
		userList = remove(userList, conn.data.User)
		conn.data.UserList = userList
		conn.data.Content = conn.data.User
		bin, _ := json.Marshal(conn.data)
		h.broadcast <- bin
		h.unregister <- conn
	}()
}
