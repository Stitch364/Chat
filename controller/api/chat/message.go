package chat

import (
	"chat/chat"
	"chat/global"
	"chat/model"
	"chat/model/chat/client"
	"chat/model/common"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	socketio "github.com/googollee/go-socket.io"
	"log"
)

// 用于处理客户端发送的 event
type message struct {
}

// SendMsg 发送消息
// 参数：client.HandleSendMsgParams
// 返回：client.HandleSendMsgRly
func (message) SendMsg(s socketio.Conn, msg map[string]interface{}) string {
	//fmt.Println("-----------------接收到的消息：", msg)
	token, ok := CheckAuth(s)
	if !ok {
		return ""
	}

	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Serialization failed: %v", err)
		return common.NewState(errcode.ErrParamsNotValid.WithDetails(err.Error())).MustJson()
	}
	//fmt.Printf("Received raw JSON: %s", string(jsonBytes))

	// 2. 反序列化到结构体
	params := new(client.HandleSendMsgParams)
	if err := json.Unmarshal(jsonBytes, &params); err != nil {
		log.Printf("Deserialization failed: %v", err)
		return common.NewState(errcode.ErrParamsNotValid.WithDetails(err.Error())).MustJson()
	}
	//base64解码
	decodedBytes, err := base64.StdEncoding.DecodeString(params.MsgContent)
	if err != nil {
		fmt.Println("Encoding error:", err)
		return common.NewState(errcode.ErrParamsNotValid.WithDetails(err.Error())).MustJson()
	}
	params.MsgContent = string(decodedBytes)
	//fmt.Println("解码的中文:-----------------------", params.MsgContent)
	//params := new(client.HandleSendMsgParams)
	//if err := common.Decode(msg, params); err != nil {
	//	return common.NewState(errcode.ErrParamsNotValid.WithDetails(err.Error())).MustJson()
	//}
	//fmt.Println("SendMsg Params:", params)
	ctx, cancel := global.DefaultContextWithTimeout()
	defer cancel()
	result, myErr := chat.Group.Message.SendMsg(ctx, &model.HandleSendMsg{
		AccessToken: token.AccessToken,
		RelationID:  params.RelationID,
		AccountID:   token.Content.ID,
		MsgContent:  params.MsgContent,
		MsgExtend:   params.MsgExtend,
		RlyMsgID:    params.RlyMsgID,
	})
	return common.NewState(myErr, result).MustJson()
}

//// ReadMsg 已读消息
//// 参数：client.HandleReadMsgParams
//// 返回：无
//func (message) ReadMsg(s socketio.Conn, msg string) string {
//	token, ok := CheckAuth(s)
//	if !ok {
//		return ""
//	}
//	params := new(client.HandleReadMsgParams)
//	if err := common.Decode(msg, params); err != nil {
//		return common.NewState(errcode.ErrParamsNotValid.WithDetails(err.Error())).MustJson()
//	}
//	ctx, cancel := global.DefaultContextWithTimeout()
//	defer cancel()
//	myErr := chat.Group.Message.ReadMsg(ctx, &model.HandleReadMsg{
//		AccessToken: token.AccessToken,
//		MsgIDs:      params.MsgIDs,
//		RelationID:  params.RelationID,
//		ReaderID:    token.Content.ID,
//	})
//	return common.NewState(myErr).MustJson()
//}
