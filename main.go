package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"go-chatroom/core"
	"net/http"

	"github.com/gorilla/mux"
)

// 定义升级器，将http请求升级为ws请求
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 跨域访问
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ws的回调函数
func WsHandler(w http.ResponseWriter, r *http.Request) {
	// 升级
	ws, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// 初始化连接对象
	conn := &core.Conn{
		WS:     ws,
		SendCh: make(chan []byte, 128),
		Data:   &core.Message{},
	}

	// 在ws中注册
	core.ChatHub.RegisterCh <- conn

	go conn.Write()
	conn.Read()

	defer func() {
		// 准备数据
		conn.Data.Type = core.Logout
		core.UserList = core.Remove(core.UserList, conn.Data.User)
		conn.Data.UserList = core.UserList
		conn.Data.Content = conn.Data.User
		bin, _ := json.Marshal(conn.Data)
		// 发送数据
		core.ChatHub.BroadcastCh <- bin
		core.ChatHub.CancelCh <- conn
	}()
}

func main() {
	// 创建路由
	r := mux.NewRouter()
	// 指定ws回调函数
	r.HandleFunc("/ws", WsHandler)

	go core.ChatHub.Run()

	// 开启服务端监听
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println("err:", err)
	}
}
