package main


// h 全局对象
var h = hub{
		connections: make(map[*connection]bool),
		broadcast :make(chan []byte),
		register :make(chan *connection),
		unregister:make(chan *connection),
	}
}