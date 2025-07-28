package chat

import (
	"chat/global"
	"chat/model/chat/client"
	"chat/model/common"
	"chat/task"
	"fmt"
	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	socketio "github.com/googollee/go-socket.io"
	"log"
	"time"
)

type handle struct {
}

const AuthLimitTimeout = 10 * time.Second

//func (handle) OnConnect(s socketio.Conn) error {
//	global.Logger.Info(fmt.Sprintf("连接建立 | ID=%s | IP=%s", s.ID(), s.RemoteAddr())) // s.RemoteAddr()获取客户端的 IP 地址和端口号信息。
//
//	return nil
//}

func (handle) OnConnect(s socketio.Conn) error {
	global.Logger.Info(fmt.Sprintf("连接建立 | ID=%s | IP=%s", s.ID(), s.RemoteAddr())) // s.RemoteAddr()获取客户端的 IP 地址和端口号信息。
	time.AfterFunc(AuthLimitTimeout, func() {
		if !global.ChatMap.HasSID(s.ID()) {
			global.Logger.Info(fmt.Sprintln("onConnect auth failed:", s.RemoteAddr().String(), s.ID()))
			_ = s.Close()
		}
	})
	return nil
}

func (handle) OnError(s socketio.Conn, err error) {
	//log.Println("OnError on error:", err)
	if s == nil {
		return
	}
	global.ChatMap.Leave(s)
	//log.Println("OnError disconnected: ", s.RemoteAddr().String(), s.ID())
	fmt.Printf("\033[32m[Error ConnMap] Error: %d sid: %s\033[0m\n", err.Error(), s.ID())
	//global.Logger.Error(fmt.Sprintf("连接错误 | ID=%s | 错误=%s", s.ID(), err.Error()))
	_ = s.Close()
}

var consumerMap = make(map[int64]bool)

func (handle) Auth(s socketio.Conn, accessToken string) string {
	//fmt.Println("Auth", accessToken)
	token, myErr := MustAccount(accessToken)
	if myErr != nil {
		return common.NewState(myErr).MustJson()
	}
	s.SetContext(token)

	global.ChatMap.Link(s, token.Content.ID)
	global.Worker.SendTask(task.AccountLogin(accessToken, s.RemoteAddr().String(), token.Content.ID))
	//go consumer.StartConsumer(token.Content.ID)
	//
	//if !consumerMap[token.Content.ID] {
	//	go consumer.StartConsumer(token.Content.ID)
	//	consumerMap[token.Content.ID] = true
	//}

	return common.NewState(nil).MustJson()
}

func (handle) Test(s socketio.Conn, msg string) string {
	_, ok := CheckAuth(s)
	if !ok {
		return ""
	}
	param := new(client.TestParams)
	log.Println(msg)
	if err := common.Decode(msg, param); err != nil {
		return common.NewState(errcode.ErrParamsNotValid.WithDetails(err.Error())).MustJson()
	}
	result := common.NewState(nil, client.TestRly{
		Name:    param.Name,
		Age:     param.Age,
		Address: s.RemoteAddr().String(),
		ID:      s.ID(),
	}).MustJson()
	s.Emit("test", "test")
	return result
}

func (handle) OnDisconnect(s socketio.Conn, _ string) {
	global.ChatMap.Leave(s)
	//fmt.Printf("\033[32m[Error ConnMap] sid: %s\033[0m\n", s.ID())
	global.Logger.Warn(fmt.Sprintf("连接断开 | ID=%s ", s.ID()))
}
