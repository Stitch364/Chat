package logic

import (
	"chat/dao"
	db "chat/dao/mysql/sqlc"
	"chat/errcodes"
	"chat/global"
	"chat/middlewares"
	"chat/model/reply"
	"chat/task"
	"database/sql"
	"errors"
	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
)

type application struct {
}

func (application) CreateApplication(ctx *gin.Context, accountID1, accountID2 int64, msg string) errcode.Err {
	//判断两个ID是否一样，不能给自己发好友申请
	if accountID1 == accountID2 {
		return errcodes.ApplicationNotValid
	}
	//判断是否是好友？（已经是好友了怎么发申请？）

	//创建申请
	err := dao.Database.DB.CreateApplicationTx(ctx, &db.CreateApplicationParams{
		Account1ID: accountID1,
		Account2ID: accountID2,
		ApplyMsg:   msg,
	})
	switch {
	case errors.Is(err, errcodes.ApplicationExists):
		return errcodes.ApplicationExists
	case errors.Is(err, nil):
		//global.Worker.SendTask(task.Application(accountID2))
		return nil
	default:
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	//将申请消息发送到被申请方
}

func (application) DeleteApplication(ctx *gin.Context, accountID1, accountID2 int64) errcode.Err {
	apply, err := getApplication(ctx, accountID1, accountID2)
	if err != nil {
		return errcode.ErrServer
	}
	//判断删申请的人是不是发送申请的人
	//只能删自己发出去的
	if apply.Account1ID != accountID1 {
		return errcodes.AuthPermissionsInsufficient
	}
	if err := dao.Database.DB.DeleteApplication(ctx, &db.DeleteApplicationParams{
		Account1ID: accountID1,
		Account2ID: accountID2,
	}); err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	return nil
}

func getApplication(ctx *gin.Context, accountID1, accountID2 int64) (*db.Application, errcode.Err) {
	apply, err := dao.Database.DB.GetApplicationByID(ctx, &db.GetApplicationByIDParams{
		Account1ID: accountID1,
		Account2ID: accountID2,
	})
	switch {
	case errors.Is(err, nil):
		return apply, nil
	case errors.Is(err, sql.ErrNoRows):
		return nil, errcodes.ApplicationNotExists
	default:
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
}

// AcceptApplication 同意申请
func (application) AcceptApplication(ctx *gin.Context, accountID1, accountID2 int64) errcode.Err {
	//申请中只有双方ID，应展示头像，名字
	apply, myerr := getApplication(ctx, accountID1, accountID2)
	if myerr != nil {
		return myerr
	}
	if apply.Status == db.ApplicationsStatusValue2 {
		return errcodes.ApplicationRepeatOpt
	}
	//获取两人的账号信息
	accountInfo1, myerr := getAccountInfoByID(ctx, accountID1, accountID1)
	if myerr != nil {
		return myerr
	}
	accountInfo2, myerr := getAccountInfoByID(ctx, accountID2, accountID2)
	if myerr != nil {
		return myerr
	}
	//同意后需要创建两人好友关系，创建两人互相的设置，推送消息
	msgInfo, err := dao.Database.DB.AcceptApplicationTx(ctx, dao.Database.Redis, accountInfo1, accountInfo2)
	//_, err := dao.Database.DB.AcceptApplicationTx(ctx, dao.Database.Redis, accountInfo1, accountInfo2)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	//推送消息
	global.Worker.SendTask(task.PublishMsg(reply.ParamMsgInfoWithRly{
		ParamMsgInfo: reply.ParamMsgInfo{
			ID:         msgInfo.ID,
			NotifyType: string(msgInfo.NotifyType),
			MsgType:    string(msgInfo.MsgType),
			MsgContent: msgInfo.MsgContent,
			RelationID: msgInfo.RelationID,
			CreateAt:   msgInfo.CreateAt,
		},
		RlyMsg: nil,
	}))
	return nil
}

// RefuseApplication 拒绝申请
// 拒绝申请只需要修改申请状态即可
func (application) RefuseApplication(ctx *gin.Context, accountID1, accountID2 int64, refuseMsg string) errcode.Err {
	//申请中只有双方ID，应展示头像，名字
	apply, myerr := getApplication(ctx, accountID1, accountID2)
	if myerr != nil {
		return myerr
	}
	if apply.Status == db.ApplicationsStatusValue1 {
		return errcodes.ApplicationRepeatOpt
	}
	if err := dao.Database.DB.UpdateApplication(ctx, &db.UpdateApplicationParams{
		Status:     db.ApplicationsStatusValue2,
		RefuseMsg:  refuseMsg,
		Account1ID: accountID1,
		Account2ID: accountID2,
	}); err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	return nil
}

func (application) ListApplications(ctx *gin.Context, accountID int64, limit, offset int32) (reply.ParamListApplication, errcode.Err) {
	//直接查询
	list, err := dao.Database.DB.GetApplications(ctx, &db.GetApplicationsParams{
		Account1ID: accountID,
		Account2ID: accountID,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return reply.ParamListApplication{}, errcode.ErrServer
	}
	if len(list) == 0 {
		//没有申请
		return reply.ParamListApplication{}, nil
	}
	//格式转换
	data := make([]*reply.ParamApplicationInfo, len(list))
	for i, v := range list {
		name, avatar := v.Account1Name, v.Account1Avatar
		if v.Account1ID == accountID {
			//显示的是哪个账号的申请
			name, avatar = v.Account2Name, v.Account2Avatar
		}
		if v.Account2ID == accountID && v.Status == db.ApplicationsStatusValue0 {
			v.Status = db.ApplicationsStatusValue3
		}
		data[i] = &reply.ParamApplicationInfo{
			AccountID1: v.Account1ID,
			AccountID2: v.Account2ID,
			ApplyMsg:   v.ApplyMsg,
			Avatar:     avatar,
			CreateAt:   v.CreateAt,
			Name:       name,
			Refuse:     v.RefuseMsg,
			Status:     string(v.Status),
			UpdateAt:   v.UpdateAt,
		}
	}
	//interface{} 转 int64
	total, ok := list[0].Total.(int64)
	if !ok {
		return reply.ParamListApplication{}, errcode.ErrServer
	}

	return reply.ParamListApplication{
		List:  data,
		Total: TotalToPageTotal(total, limit),
	}, nil

}
