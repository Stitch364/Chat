package reply

import "time"

type ParamApplicationInfo struct {
	AccountID1 int64     `json:"account_id_1,omitempty"` //申请者ID
	AccountID2 int64     `json:"account_id_2,omitempty"` //被申请者的ID
	ApplyMsg   string    `json:"apply_msg,omitempty"`    //申请消息
	Refuse     string    `json:"refuse,omitempty"`       //拒绝信息
	Status     string    `json:"status,omitempty"`       //申请状态[已申请，已拒绝，已同意]
	CreateAt   time.Time `json:"create_at,omitempty"`    //创建时间
	UpdateAt   time.Time `json:"update_at,omitempty"`    //更新时间
	Name       string    `json:"name,omitempty"`         //对方名字
	Avatar     string    `json:"avatar,omitempty"`       //对方头像
}

// ParamListApplication 所有申请信息
type ParamListApplication struct {
	List  []*ParamApplicationInfo `json:"list,omitempty"`
	Total int64                   `json:"total,omitempty"` //总数
}
