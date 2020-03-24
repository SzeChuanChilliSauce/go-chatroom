package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// 创建路由
	r := mux.NewRouter()
	// 指定ws回调函数
	r.HandleFunc("/ws", wsHandler)

	go Hub.run()

	// 开启服务端监听
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println("err:", err)
	}
}
