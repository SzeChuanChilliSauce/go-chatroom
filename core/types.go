package core

// 消息类型
// 	 login - 登录信息
// 	 handshake - 握手信息，刚打开网页的状态
//   system - 系统信息
// 	 logout - 退出信息
// 	 user - 普通信息
const (
	Login     = "login"
	Logout    = "logout"
	Handshake = "handshake"
	System    = "system"
	User      = "user"
)

// Message
type Message struct {
	IP string `json:"ip"`
	// 消息类型
	Type string `json:"type"`
	// 消息发出者
	From string `json:"from"`
	// 内容
	Content string `json:"content"`
	// 用户名
	User string `json:"user"`
	// 用户列表
	UserList []string `json:"user_list"`
}

// 从列表中删除元素
func Remove(slice []string, user string) []string {
	l := len(slice)
	if l == 0 {
		return slice
	}
	// 是第一个元素
	if slice[0] == user {
		return slice[1:]
	}
	// 是最后一个元素
	if slice[l-1] == user {
		return slice[:l-1]
	}

	var res = []string{}
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
