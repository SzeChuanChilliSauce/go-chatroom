package core

// Data ...
type Data struct {
	IP string `json:"ip"`

	// 标识信息类型
	// 	login - 登录信息
	// 	handshake - 握手信息，刚打开网页的状态
	// 	system - 系统信息
	// 	logout - 退出信息
	// 	user - 普通信息
	Type string `json:"type"`

	// 消息由谁发出
	From string `json:"from"`

	// 内容
	Content string `json:"content"`

	// 用户名
	User string `json:"user"`

	// 用户列表
	UserList []string `json:"user_list"`
}
