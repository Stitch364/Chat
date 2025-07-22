package request

import "time"

type ParamCreatApplication struct {
	AccountID      int64  `json:"account_id" binding:"required,gte=1"` //被申请者的ID
	ApplicationMsg string `json:"application_msg" binding:"lte=200"`   //申请信息
}

type ParamDeleteApplication struct {
	AccountID int64     `json:"account_id" binding:"required,gte=1"` //被申请者的ID
	CreatAt   time.Time `json:"create_at" binding:"required"`        //申请时间
}

type ParamRefuseApplication struct {
	AccountID int64     `json:"account_id" binding:"required,gte=1"` //申请者的ID
	RefuseMsg string    `json:"refuse_msg" binding:"lte=200"`        //拒绝信息
	CreatAt   time.Time `json:"create_at" binding:"required"`        //申请时间
}
type ParamAcceptApplication struct {
	AccountID int64     `json:"account_id" binding:"required,gte=1"` //申请者的ID
	CreatAt   time.Time `json:"create_at" binding:"required"`        //申请时间
}
