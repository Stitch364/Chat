package chat

import (
	"chat/dao"
	db "chat/dao/mysql/sqlc"
	"chat/errcodes"
	"chat/global"
	"chat/logic"
	"chat/model"
	"chat/model/chat/client"
	"chat/model/reply"
	"chat/task"
	"context"
	"database/sql"
	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	"time"
)

type message struct {
}

// SendMsg 发送消息
func (message) SendMsg(ctx context.Context, params *model.HandleSendMsg) (*client.HandleSendMsgRly, errcode.Err) {
	//判断权限
	ok, myErr := logic.ExistsSetting(ctx, params.AccountID, params.RelationID)
	if myErr != nil {
		return nil, myErr
	}
	if !ok {
		return nil, errcodes.AuthPermissionsInsufficient
	}

	var rlyMsgID int64
	var rlyMsg *reply.ParamRlyMsg
	//判断是否是回复的消息，如果是并回复
	if params.RlyMsgID > 0 {
		//是回复的消息
		//获取被回复的消息信息
		rlyInfo, myerr := logic.GetMsgInfoAndNameByID(ctx, params.RlyMsgID)
		if myerr != nil {
			return nil, myerr
		}
		//不能回复别的群的消息
		if rlyInfo.RelationID != params.RelationID {
			return nil, errcodes.RlyMsgNotOneRelation
		}
		//不能回复已撤回的消息
		if rlyInfo.IsRevoke {
			return nil, errcodes.RlyMsgHasRevoked
		}
		//回复的消息ID
		rlyMsgID = params.RlyMsgID
		//消息扩展信息
		rlyMsgExtend, err := model.JsonToExtend(rlyInfo.MsgExtend)
		if err != nil {
			return nil, errcode.ErrServer
		}

		//是文件类消息
		fileInfo := &db.File{
			ID:         0,
			FileName:   "",
			FileType:   "",
			FileSize:   0,
			FileKey:    "",
			Url:        "",
			RelationID: sql.NullInt64{},
			AccountID:  sql.NullInt64{},
			CreateAt:   time.Time{},
		}
		var myErr1 error
		if rlyInfo.MsgType == db.MessagesMsgTypeFile {
			fileInfo, myErr1 = dao.Database.DB.GetFileDetailsByID(ctx, rlyInfo.FileID.Int64)
			if myErr1 != nil {
				//查询出错了
				global.Logger.Error(myErr1.Error())
				return nil, errcode.ErrServer
			}
		}
		//被回复的消息信息
		rlyMsg = &reply.ParamRlyMsg{
			MsgID:         rlyInfo.ID, //被回复的消息ID
			MsgContent:    rlyInfo.MsgContent,
			MsgExtend:     rlyMsgExtend,
			MsgType:       string(rlyInfo.MsgType),
			IsRevoked:     rlyInfo.IsRevoke,
			AccountID:     rlyInfo.AccountID.Int64,
			AccountName:   rlyInfo.Name,
			NickName:      rlyInfo.NickName,
			AccountAvatar: rlyInfo.Avatar,
			FileID:        fileInfo.ID,
			FileName:      fileInfo.FileName,
			FileType:      string(fileInfo.FileType),
			FileSize:      fileInfo.FileSize,
		}
	}
	//msgExtend, err := model.ExtendToJson(params.MsgExtend)
	//if err != nil {
	//	global.Logger.Error(err.Error())
	//	return nil, errcode.ErrServer
	//}
	//将消息存储到数据库
	result, err := dao.Database.DB.CreateMessageTx(ctx, &db.CreateMessageParams{
		NotifyType: db.MessagesNotifyTypeCommon,           //通知类型
		MsgType:    db.MessagesMsgType(model.MsgTypeText), //消息类型（文本消息）
		MsgContent: params.MsgContent,                     //消息内容
		//MsgExtend:  msgExtend,                                         //扩展信息
		AccountID:  sql.NullInt64{Int64: params.AccountID, Valid: true}, //发送账号
		RlyMsgID:   sql.NullInt64{Int64: rlyMsgID, Valid: rlyMsgID > 0}, //回复的消息ID
		RelationID: params.RelationID,                                   //关系ID
	})
	if err != nil {
		global.Logger.Error(err.Error())
		return nil, errcode.ErrServer
	}
	//推送消息
	global.Worker.SendTask(task.PublishMsg(reply.ParamMsgInfoWithRly{
		ParamMsgInfo: reply.ParamMsgInfo{
			ID:            result.ID,
			NotifyType:    string(db.MessagesNotifyTypeCommon), //common
			MsgType:       string(model.MsgTypeText),           //text
			MsgContent:    result.MsgContent,                   //消息内容
			MsgExtend:     params.MsgExtend,                    //扩展信息
			AccountID:     params.AccountID,                    //发送的账号ID
			AccountName:   result.Name,                         //账号名称
			AccountAvatar: result.Avatar,                       //账号头像
			NickName:      result.NickName,                     //昵称
			RelationID:    params.RelationID,                   //关系ID
			CreateAt:      result.CreateAt,                     //消息创建时间
		},
		RlyMsg: rlyMsg, //被回复的消息信息
	}))
	//返回消息ID和创建时间
	return &client.HandleSendMsgRly{
		MsgID:    result.ID,
		CreateAt: result.CreateAt,
	}, nil
}

//func (message) ReadMsg(ctx context.Context, params *model.HandleReadMsg) errcode.Err {
//	//判断权限
//	ok, myErr := logic.ExistsSetting(ctx, params.ReaderID, params.RelationID)
//	if myErr != nil {
//		return myErr
//	}
//	if !ok {
//		return errcodes.AuthPermissionsInsufficient
//	}
//	//没有做已读消息的功能
//
//	return nil
//}
