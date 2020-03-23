package main

// h 全局对象
var h = hub{
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

		}
	}
}
