package model

import (
	"encoding/json"
	"time"
)

type MsgType string

const (
	MsgTypeText MsgType = "text"
	MsgTypeFile MsgType = "file"
)

type Remind struct {
	Idx       int64 `json:"idx,omitempty" binding:"required,gte=1" validate:"required,gte=1"`        // 第几个 @
	AccountID int64 `json:"account_id,omitempty" binding:"required,gte=1" validate:"required,gte=1"` // 被 @ 的账号 ID
}

type MsgExtend struct {
	Remind []Remind `json:"remind"` // @ 的描述信息
}

type GetMsgsByRelationIDAndTime struct {
	AccountID  int64
	RelationID int64
	LastTime   time.Time
	Limit      int32
	Offset     int32
}

type OfferMsgsByAccountIDAndTime struct {
	AccountID int64
	LastTime  time.Time
	Limit     int32
	Offset    int32
}

// ExtendToJson 将 MsgExtend 转化为 json格式的字符串，可以是 nil
// 参数：消息扩展信息
// 返回：json格式的 string 对象
func ExtendToJson(extend *MsgExtend) (json.RawMessage, error) {
	if extend == nil {
		return nil, nil // 如果 extend 是 nil，返回 nil
	}
	data, err := json.Marshal(extend)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// JsonToExtend 将 Json 转化为 MsgExtend
// 参数：Json 对象（如果存的 json 为 nil 或未定义则返回 nil）
// 返回：解析后的消息扩展信息（可能为 nil）
func JsonToExtend(data json.RawMessage) (*MsgExtend, error) {
	if data == nil {
		return nil, nil
	}
	extend := &MsgExtend{}
	err := json.Unmarshal(data, extend)
	if err != nil {
		return nil, err
	}
	return extend, nil
}
