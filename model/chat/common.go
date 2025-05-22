package chat

/*
chat 中事件的关键字，以 Server 开头的事件为服务端发送的事件，以 client 开头的事件为客户端发送的事件
*/ // 客户端推送的事件

const (
	ClientSendMsg = "send_msg" // 发送消息
	ClientReadMsg = "read_msg" // 已读消息
	ClientTest    = "test"     // 测试
	ClientAuth    = "auth"     // 认证
)
