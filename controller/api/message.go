package api

import (
	"chat/errcodes"
	"chat/global"
	"chat/logic"
	"chat/middlewares"
	"chat/model"
	"chat/model/request"
	"github.com/XYYSWK/Lutils/pkg/app"
	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
	"time"
)

type message struct {
}

func (message) GetMsgsByRelationIDAndTime(ctx *gin.Context) {
	//1.获取参数和参数校验
	//NewResponse就是包装一下ctx
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetMsgsByRelationIDAndTime)
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	//获取Token信息
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	//2.业务处理
	limit, offset := global.Page.GetPageSizeAndOffset(ctx.Request)
	//fmt.Println("params: ", params)
	//fmt.Println("params.LastTime: ", time.Unix(params.LastTime/1000, 0))
	result, err := logic.Logics.Message.GetMsgsByRelationIDAndTime(ctx, model.GetMsgsByRelationIDAndTime{
		AccountID:  content.ID,
		RelationID: params.RelationID,
		LastTime:   time.Unix(params.LastTime/1000, 0),
		Limit:      limit,
		Offset:     offset,
	})
	//3.返回响应
	reply.Reply(err, result)
}

func (message) GetPinMsgsByRelationID(ctx *gin.Context) {
	//1.获取参数和参数校验
	//NewResponse就是包装一下ctx
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetPinMsgsByRelationID)
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	//获取Token信息
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	//2.业务处理
	limit, offset := global.Page.GetPageSizeAndOffset(ctx.Request)

	result, err := logic.Logics.Message.GetPinMsgsByRelationID(ctx, content.ID, params.RelationID, limit, offset)

	//3.返回响应
	reply.Reply(err, result)
}

func (message) GetRlyMsgsInfoByMsgID(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetRlyMsgsInfoByMsgID)
	if err := ctx.ShouldBindQuery(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	limit, offset := global.Page.GetPageSizeAndOffset(ctx.Request)
	result, err := logic.Logics.Message.GetRlyMsgsInfoByMsgID(ctx, content.ID, params.RelationID, params.MsgID, limit, offset)
	reply.ReplyList(err, result.Total, result.List)
}

func (message) GetMsgsByContent(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetMsgsByContent)
	if err := ctx.ShouldBindQuery(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}

	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	limit, offset := global.Page.GetPageSizeAndOffset(ctx.Request)
	result, err := logic.Logics.Message.GetMsgsByContent(ctx, content.ID, params.RelationID, params.Content, limit, offset)
	reply.ReplyList(err, result.Total, result.List)
}

func (message) UpdateMsgPin(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamUpdateMsgPin)
	if err := ctx.ShouldBindQuery(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}

	err := logic.Logics.Message.UpdateMsgPin(ctx, content.ID, params)
	reply.Reply(err)
}
func (message) UpdateMsgTop(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamUpdateMsgTop)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Message.UpdateMsgTop(ctx, content.ID, params)
	reply.Reply(err)
}

func (message) RevokeMsg(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamRevokeMsg)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Message.RevokeMsg(ctx, content.ID, params.ID)
	reply.Reply(err)
}

func (message) DeleteMsg(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	//消息ID
	params := new(request.ParamDeleteMsg)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Message.DeleteMsg(ctx, content.ID, params.ID)
	reply.Reply(err)
}

func (message) GetTopMsgByRelationID(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetTopMsgByRelationID)
	if err := ctx.ShouldBindQuery(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.Message.GetTopMsgByRelationID(ctx, content.ID, params.RelationID)
	reply.Reply(err, result)
}
