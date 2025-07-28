package logic

import (
	"chat/dao"
	db "chat/dao/mysql/sqlc"
	"chat/errcodes"
	"chat/global"
	"chat/middlewares"
	"chat/model"
	"chat/model/chat/server"
	"chat/model/format"
	"chat/model/reply"
	"chat/model/request"
	"chat/task"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type message struct {
}

func (message) GetMsgsByRelationIDAndTime(ctx *gin.Context, params model.GetMsgsByRelationIDAndTime) (*reply.ParamGetMsgsRelationIDAndTime, errcode.Err) {
	//通过时间和ids查询消息
	// 权限验证
	ok, myErr := ExistsSetting(ctx, params.AccountID, params.RelationID)
	if myErr != nil {
		return nil, myErr
	}
	if !ok {
		return nil, errcodes.AuthPermissionsInsufficient
	}
	data, err := dao.Database.DB.GetMsgsByRelationIDAndTime(ctx, &db.GetMsgsByRelationIDAndTimeParams{
		RelationID:   params.RelationID,
		RelationID_2: params.RelationID,
		CreateAt:     params.LastTime,
		Limit:        params.Limit,
		Offset:       params.Offset,
	})
	if err != nil {
		//查询出错了
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	//没有消息
	if len(data) == 0 {
		return &reply.ParamGetMsgsRelationIDAndTime{List: []*reply.ParamMsgInfoWithRly{}}, nil
	}
	//查询用户备注
	//NickName, err := dao.Database.DB.GetNickNameByAccountIDAndRelation(ctx, &db.GetNickNameByAccountIDAndRelationParams{
	//	AccountID:  params.AccountID,
	//	RelationID: params.RelationID,
	//})

	result := make([]*reply.ParamMsgInfoWithRly, 0, len(data))
	var DelMsgSum int64 = 0
	for _, v := range data {
		var content string
		var extend *model.MsgExtend
		if !v.IsRevoke { // 该消息没有被撤回
			//添加消息内容以及扩展信息
			content = v.MsgContent
			extend, err = model.JsonToExtend(v.MsgExtend)
			if err != nil {
				global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
				continue
			}
		}
		//var readIDs []int64
		//if params.AccountID == v.AccountID.Int64 {
		//	readIDs = v.ReadIds
		//}
		var rlyMsg *reply.ParamRlyMsg
		if v.RlyMsgID.Valid { // 该 ID 有意义
			rlyMsgInfo, myErr := GetMsgInfoByID(ctx, v.RlyMsgID.Int64)
			if myErr != nil {
				continue
			}
			var rlyContent string
			var rlyExtend *model.MsgExtend
			if !rlyMsgInfo.IsRevoke { // 回复消息没有撤回
				rlyContent = rlyMsgInfo.MsgContent
				rlyExtend, err = model.JsonToExtend(rlyMsgInfo.MsgExtend)
				if err != nil {
					global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
					continue
				}
			}
			rlyMsg = &reply.ParamRlyMsg{
				MsgID:      v.RlyMsgID.Int64,
				MsgType:    string(rlyMsgInfo.MsgType),
				MsgContent: rlyContent,
				MsgExtend:  rlyExtend,
				IsRevoked:  rlyMsgInfo.IsRevoke,
			}
		}
		//不显示自己删除的消息
		if params.AccountID == v.AccountID.Int64 && v.IsDelete == 1 || v.IsDelete == 2 {
			DelMsgSum++
			continue
		}
		if params.AccountID != v.AccountID.Int64 && v.IsDelete == -1 || v.IsDelete == 2 {
			DelMsgSum++
			continue
		}
		result = append(result, &reply.ParamMsgInfoWithRly{
			ParamMsgInfo: reply.ParamMsgInfo{
				ID:         v.ID,
				NotifyType: string(v.NotifyType),
				MsgType:    string(v.MsgType),
				MsgContent: content,
				MsgExtend:  extend,
				FileID:     v.FileID.Int64,
				AccountID:  v.AccountID.Int64,
				RelationID: v.RelationID,
				CreateAt:   v.CreateAt,
				IsRevoke:   v.IsRevoke,
				IsTop:      v.IsTop,
				IsPin:      v.IsPin,
				PinTime:    v.PinTime,
				//ReadIds:    readIDs,
				ReplyCount: v.ReplyCount,
			},
			RlyMsg:   rlyMsg,
			NickName: v.NickName,
		})
	}
	return &reply.ParamGetMsgsRelationIDAndTime{List: result, Total: TotalToPageTotal(data[0].Total.(int64), params.Limit) - DelMsgSum}, nil
}

func GetMsgInfoByID(ctx context.Context, msgID int64) (*db.Message, errcode.Err) {
	result, err := dao.Database.DB.GetMessageByID(ctx, msgID)
	if err != nil {
		//数据库中不存在该消息
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errcodes.MsgNotExists
		}
		global.Logger.Error(err.Error())
		return nil, errcode.ErrServer
	}
	return result, nil
}

func (message) GetPinMsgsByRelationID(ctx *gin.Context, accountID, relationID int64, limit, offset int32) (*reply.ParamGetPinMsgsByRelationID, errcode.Err) {
	//  权限验证
	ok, err := ExistsSetting(ctx, accountID, relationID)
	if err != nil {
		return &reply.ParamGetPinMsgsByRelationID{Total: 0}, err
	}
	if !ok {
		return &reply.ParamGetPinMsgsByRelationID{Total: 0}, errcodes.AuthPermissionsInsufficient
	}
	//查数据库
	data, myErr := dao.Database.DB.GetPinMsgsByRelationID(ctx, &db.GetPinMsgsByRelationIDParams{
		RelationID:   relationID,
		RelationID_2: relationID,
		Limit:        limit,
		Offset:       offset,
	})
	if myErr != nil {
		global.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
		return &reply.ParamGetPinMsgsByRelationID{Total: 0}, errcode.ErrServer
	}
	//没数据
	if len(data) == 0 {
		return &reply.ParamGetPinMsgsByRelationID{List: []*reply.ParamMsgInfo{}}, nil
	}
	result := make([]*reply.ParamMsgInfo, 0, len(data))
	//格式转换
	for _, v := range data {
		var content string
		var extend *model.MsgExtend
		if !v.IsRevoke {
			content = v.MsgContent
			extend, myErr = model.JsonToExtend(v.MsgExtend)
			if myErr != nil {
				global.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
				return &reply.ParamGetPinMsgsByRelationID{Total: 0}, errcode.ErrServer
			}
		}

		result = append(result, &reply.ParamMsgInfo{
			ID:         v.ID,
			NotifyType: string(v.NotifyType),
			MsgType:    string(v.MsgType),
			MsgContent: content,
			MsgExtend:  extend,
			FileID:     v.FileID.Int64,
			AccountID:  v.AccountID.Int64,
			RelationID: v.RelationID,
			CreateAt:   v.CreateAt,
			IsRevoke:   v.IsRevoke,
			IsTop:      v.IsTop,
			IsPin:      v.IsPin,
			PinTime:    v.PinTime,
			ReplyCount: v.ReplyCount,
		})
	}
	return &reply.ParamGetPinMsgsByRelationID{
		List:  result,
		Total: TotalToPageTotal(data[0].Total.(int64), limit),
	}, nil
}

func (message) GetRlyMsgsInfoByMsgID(ctx *gin.Context, accountID, relationID, msgID int64, limit, offset int32) (*reply.ParamGetRlyMsgsInfoByMsgID, errcode.Err) {
	//  权限验证
	ok, err := ExistsSetting(ctx, accountID, relationID)
	if err != nil {
		return &reply.ParamGetRlyMsgsInfoByMsgID{Total: 0}, err
	}
	if !ok {
		return &reply.ParamGetRlyMsgsInfoByMsgID{Total: 0}, errcodes.AuthPermissionsInsufficient
	}
	//查数据库
	data, myErr := dao.Database.DB.GetRlyMsgsInfoByMsgID(ctx, &db.GetRlyMsgsInfoByMsgIDParams{
		RelationID:   relationID,
		RelationID_2: relationID,
		RlyMsgID:     sql.NullInt64{Int64: msgID, Valid: true},
		Limit:        limit,
		Offset:       offset,
	})
	if myErr != nil {
		global.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
		return &reply.ParamGetRlyMsgsInfoByMsgID{Total: 0}, errcode.ErrServer
	}
	//没数据
	if len(data) == 0 {
		return &reply.ParamGetRlyMsgsInfoByMsgID{List: []*reply.ParamMsgInfo{}}, nil
	}
	result := make([]*reply.ParamMsgInfo, 0, len(data))
	//格式转换
	for _, v := range data {
		var content string
		var extend *model.MsgExtend
		if !v.IsRevoke {
			content = v.MsgContent
			extend, myErr = model.JsonToExtend(v.MsgExtend)
			if myErr != nil {
				global.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
				return &reply.ParamGetRlyMsgsInfoByMsgID{Total: 0}, errcode.ErrServer
			}
		}

		result = append(result, &reply.ParamMsgInfo{
			ID:         v.ID,
			NotifyType: string(v.NotifyType),
			MsgType:    string(v.MsgType),
			MsgContent: content,
			MsgExtend:  extend,
			FileID:     v.FileID.Int64,
			AccountID:  v.AccountID.Int64,
			RelationID: v.RelationID,
			CreateAt:   v.CreateAt,
			IsRevoke:   v.IsRevoke,
			IsTop:      v.IsTop,
			IsPin:      v.IsPin,
			PinTime:    v.PinTime,
			ReplyCount: v.ReplyCount,
		})
	}
	return &reply.ParamGetRlyMsgsInfoByMsgID{
		List:  result,
		Total: TotalToPageTotal(data[0].Total.(int64), limit),
	}, nil
}

// 从指定关系中模糊查找指定内容的信息
func getMsgsByContentAndRelation(ctx *gin.Context, params *db.GetMsgsByContentAndRelationParams) (*reply.ParamGetMsgsByContent, errcode.Err) {
	ok, err := ExistsSetting(ctx, params.AccountID, params.RelationID)
	if err != nil {
		return &reply.ParamGetMsgsByContent{}, err
	}
	if !ok {
		return &reply.ParamGetMsgsByContent{}, errcodes.AuthPermissionsInsufficient
	}
	data, myErr := dao.Database.DB.GetMsgsByContentAndRelation(ctx, params)
	if myErr != nil {
		global.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
		return &reply.ParamGetMsgsByContent{}, errcode.ErrServer
	}
	if len(data) == 0 {
		return &reply.ParamGetMsgsByContent{List: []*reply.ParamBriefMsgInfo{}}, nil
	}
	result := make([]*reply.ParamBriefMsgInfo, 0, len(data))
	var DelMsgSum int64 = 0
	for _, v := range data {
		var extend *model.MsgExtend
		extend, myErr = model.JsonToExtend(v.MsgExtend)
		if myErr != nil {
			global.Logger.Error(myErr.Error(), zap.Any("msgExtend", v.MsgExtend))
			continue
		}
		//不显示自己删除的消息
		if params.AccountID == v.AccountID.Int64 && v.IsDelete == 1 || v.IsDelete == 2 {
			DelMsgSum++
			continue
		}
		if params.AccountID != v.AccountID.Int64 && v.IsDelete == -1 || v.IsDelete == 2 {
			DelMsgSum++
			continue
		}
		result = append(result, &reply.ParamBriefMsgInfo{
			ID:         v.ID,
			NotifyType: string(v.NotifyType),
			MsgType:    string(v.MsgType),
			MsgContent: v.MsgContent,
			Extend:     extend,
			FileID:     v.FileID.Int64,
			AccountID:  v.AccountID.Int64,
			RelationID: v.RelationID,
			CreateAt:   v.CreateAt,
		})
	}
	return &reply.ParamGetMsgsByContent{
		List:  result,
		Total: TotalToPageTotal(data[0].Total.(int64)-DelMsgSum, params.Limit),
	}, nil
}

// GetMsgsByContent 从所有用户消息中查和从指定用户消息中查
func (message) GetMsgsByContent(ctx *gin.Context, accountID, relationID int64, content string, limit, offset int32) (*reply.ParamGetMsgsByContent, errcode.Err) {

	if relationID >= 0 {
		//模糊查找指定关系中的聊天信息
		return getMsgsByContentAndRelation(ctx, &db.GetMsgsByContentAndRelationParams{
			RelationID: relationID,
			AccountID:  accountID,
			CONCAT:     content,
			Limit:      limit,
			Offset:     offset,
		})
	}
	//模糊查找所有关系中的信息
	//查数据库
	data, myErr := dao.Database.DB.GetMsgsByContent(ctx, &db.GetMsgsByContentParams{
		AccountID: accountID,
		CONCAT:    content,
		Limit:     limit,
		Offset:    offset,
	})
	if myErr != nil {
		global.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
		return &reply.ParamGetMsgsByContent{Total: 0}, errcode.ErrServer
	}
	//没数据
	if len(data) == 0 {
		return &reply.ParamGetMsgsByContent{List: []*reply.ParamBriefMsgInfo{}}, nil
	}
	result := make([]*reply.ParamBriefMsgInfo, 0, len(data))
	//格式转换
	var DelMsgSum int64 = 0
	for _, v := range data {
		extend, myErr := model.JsonToExtend(v.MsgExtend)
		if myErr != nil {
			global.Logger.Error(myErr.Error(), zap.Any("msgExtend", v.MsgExtend))
			continue
		}
		//不显示自己删除的消息
		if accountID == v.AccountID.Int64 && v.IsDelete == 1 || v.IsDelete == 2 {
			DelMsgSum++
			continue
		}
		if accountID != v.AccountID.Int64 && v.IsDelete == -1 || v.IsDelete == 2 {
			DelMsgSum++
			continue
		}
		result = append(result, &reply.ParamBriefMsgInfo{
			ID:         v.ID,
			NotifyType: string(v.NotifyType),
			MsgType:    string(v.MsgType),
			MsgContent: v.MsgContent,
			Extend:     extend,
			FileID:     v.FileID.Int64,
			AccountID:  v.AccountID.Int64,
			RelationID: v.RelationID,
			CreateAt:   v.CreateAt,
		})
	}
	return &reply.ParamGetMsgsByContent{
		List:  result,
		Total: TotalToPageTotal(data[0].Total.(int64)-DelMsgSum, limit),
	}, nil
}

func (message) UpdateMsgPin(ctx *gin.Context, accountID int64, params *request.ParamUpdateMsgPin) errcode.Err {
	ok, err := ExistsSetting(ctx, accountID, params.RelationID)
	if err != nil {
		return err
	}
	if !ok {
		return errcodes.AuthPermissionsInsufficient
	}
	msgInfo, err := GetMsgInfoByID(ctx, params.ID)
	if err != nil {
		return err
	}
	if msgInfo.IsPin == params.IsPin {
		return nil
	}
	myErr := dao.Database.DB.UpdateMsgPin(ctx, &db.UpdateMsgPinParams{
		ID:    params.ID,
		IsPin: params.IsPin,
	})
	if myErr != nil {
		global.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	// 推送 pin 通知
	accessToken, _ := middlewares.GetToken(ctx.Request.Header)
	global.Worker.SendTask(task.UpdateMsgState(accessToken, params.RelationID, params.ID, server.MsgPin, params.IsPin))
	return nil
}

func (message) UpdateMsgTop(ctx *gin.Context, accountID int64, params *request.ParamUpdateMsgTop) errcode.Err {
	//权限验证
	ok, err := ExistsSetting(ctx, accountID, params.RelationID)
	if err != nil {
		return err
	}
	if !ok {
		return errcodes.AuthPermissionsInsufficient
	}
	msgInfo, err := GetMsgInfoByID(ctx, params.ID)
	if err != nil {
		return err
	}
	if msgInfo.IsTop == params.IsTop {
		return nil
	}
	myErr := dao.Database.DB.UpdateMsgTop(ctx, &db.UpdateMsgTopParams{
		ID:    params.ID,
		IsTop: params.IsTop,
	})
	if myErr != nil {
		global.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	// 推送 置顶 消息
	accessToken, _ := middlewares.GetToken(ctx.Request.Header)
	global.Worker.SendTask(task.UpdateMsgState(accessToken, params.RelationID, params.ID, server.MsgTop, params.IsTop))
	// 创建并推送 top 消息
	f := func() error {
		arg := &db.CreateMessageParams{
			NotifyType: db.MessagesNotifyTypeSystem,
			MsgType:    db.MessagesMsgType(model.MsgTypeText),
			MsgContent: fmt.Sprintf(format.TopMessage, accountID),
			//MsgExtend:  json.RawMessage{},
			RelationID: msgInfo.RelationID,
		}
		msgRly, err := dao.Database.DB.CreateMessageTx(ctx, arg)
		if err != nil {
			return err
		}
		global.Worker.SendTask(task.PublishMsg(reply.ParamMsgInfoWithRly{
			ParamMsgInfo: reply.ParamMsgInfo{
				ID:         msgRly.ID,
				NotifyType: string(arg.NotifyType),
				MsgType:    string(arg.MsgType),
				MsgContent: arg.MsgContent,
				RelationID: arg.RelationID,
				CreateAt:   msgRly.CreateAt,
			},
			RlyMsg: nil,
		}))
		return nil
	}
	if err := f(); err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		reTry("UpdateMsgTop", f)
	}
	return nil
}

func (message) RevokeMsg(ctx *gin.Context, accountID, msgID int64) errcode.Err {
	msgInfo, err := GetMsgInfoByID(ctx, msgID)
	if err != nil {
		return err
	}
	// 检查权限(是不是本人)
	if msgInfo.AccountID.Int64 != accountID {
		return errcodes.AuthPermissionsInsufficient
	}
	if msgInfo.IsRevoke {
		return errcodes.MsgAlreadyRevoke
	}
	if diffMinutes := int(time.Since(msgInfo.CreateAt).Minutes()); diffMinutes > 2 {
		return errcodes.MsgRevokeTimeOut
	}
	myErr := dao.Database.DB.RevokeMsgWithTx(ctx, msgID, msgInfo.IsPin, msgInfo.IsTop)
	if myErr != nil {
		global.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	accessToken, _ := middlewares.GetToken(ctx.Request.Header)
	global.Worker.SendTask(task.UpdateMsgState(accessToken, msgInfo.RelationID, msgID, server.MsgRevoke, true))
	if msgInfo.IsTop {
		// 推送 top 通知
		global.Worker.SendTask(task.UpdateMsgState(accessToken, msgInfo.RelationID, msgID, server.MsgTop, false))
		// 创建并推送 top 消息
		f := func() error {
			arg := &db.CreateMessageParams{
				NotifyType: db.MessagesNotifyTypeSystem,
				MsgType:    db.MessagesMsgType(model.MsgTypeText),
				MsgContent: fmt.Sprintf(format.UnTopMessage, accountID),
				//MsgExtend:  json.RawMessage{},
				RelationID: msgInfo.RelationID,
			}
			msgRly, err := dao.Database.DB.CreateMessageTx(ctx, arg)
			if err != nil {
				return err
			}
			global.Worker.SendTask(task.PublishMsg(reply.ParamMsgInfoWithRly{
				ParamMsgInfo: reply.ParamMsgInfo{
					ID:         msgRly.ID,
					NotifyType: string(arg.NotifyType),
					MsgType:    string(arg.MsgType),
					MsgContent: arg.MsgContent,
					RelationID: arg.RelationID,
					CreateAt:   msgRly.CreateAt,
				},
				RlyMsg: nil,
			}))
			return nil
		}
		if err := f(); err != nil {
			global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
			reTry("RevokeMsg", f)
		}
	}
	return nil
}

func (message) DeleteMsg(ctx *gin.Context, accountID, msgID int64) errcode.Err {
	var Isdel int32
	msgInfo, err := GetMsgInfoByID(ctx, msgID)
	if err != nil {
		return err
	}
	// 检查是删除自己发的消息还是别人发的消息
	if msgInfo.IsRevoke {
		return errcodes.MsgAlreadyRevoke
	}
	//发消息人不是自己
	//fmt.Println(msgInfo.AccountID.Int64)
	//fmt.Println(accountID)
	if msgInfo.AccountID.Int64 != accountID {
		Isdel = -1
	} else {
		//发消息人是自己
		Isdel = 1
	}

	myErr := dao.Database.DB.DeleteMsgWithTx(ctx, msgID, Isdel)
	if myErr != nil {
		global.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	return nil
}

func (message) GetTopMsgByRelationID(ctx *gin.Context, accountID, relationID int64) (*reply.ParamGetTopMsgByRelationID, errcode.Err) {
	ok, err := ExistsSetting(ctx, accountID, relationID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errcodes.AuthPermissionsInsufficient
	}
	data, myErr := dao.Database.DB.GetTopMsgByRelationID(ctx, &db.GetTopMsgByRelationIDParams{
		RelationID:   relationID,
		RelationID_2: relationID,
	})
	if myErr != nil {
		if errors.Is(myErr, sql.ErrNoRows) {
			return nil, nil
		}
		global.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	var content string
	var extend *model.MsgExtend
	if !data.IsRevoke {
		content = data.MsgContent
		extend, myErr = model.JsonToExtend(data.MsgExtend)
		if myErr != nil {
			global.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
			return nil, errcode.ErrServer
		}
	}

	return &reply.ParamGetTopMsgByRelationID{MsgInfo: reply.ParamMsgInfo{
		ID:         data.ID,
		NotifyType: string(data.NotifyType),
		MsgType:    string(data.MsgType),
		MsgContent: content,
		MsgExtend:  extend,
		FileID:     data.FileID.Int64,
		AccountID:  data.AccountID.Int64,
		RelationID: data.RelationID,
		CreateAt:   data.CreateAt,
		IsRevoke:   data.IsRevoke,
		IsTop:      data.IsTop,
		IsPin:      data.IsPin,
		PinTime:    data.PinTime,
		ReplyCount: data.ReplyCount,
	}}, nil
}

func TotalToPageTotal(total int64, limit int32) int64 {
	if total%int64(limit) == 0 {
		return total / int64(limit)
	}
	return total/int64(limit) + 1
}
