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
)

type application struct {
}

func (application) CreateApplication(ctx *gin.Context) {
	//1.获取参数和参数校验
	//NewResponse就是包装一下ctx
	reply := app.NewResponse(ctx)
	params := new(request.ParamCreatApplication)
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	//获取UserToken信息
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	//2.业务处理
	err := logic.Logics.Application.CreateApplication(ctx, content.ID, params.AccountID, params.ApplicationMsg)

	//3.返回响应
	reply.Reply(err)
}

func (application) DeleteApplication(ctx *gin.Context) {
	//1.获取参数和参数校验
	//NewResponse就是包装一下ctx
	reply := app.NewResponse(ctx)
	params := new(request.ParamDeleteApplication)
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
	err := logic.Logics.Application.DeleteApplication(ctx, content.ID, params.AccountID, params.CreatAt)

	//3.返回响应
	reply.Reply(err)
}

func (application) AcceptApplication(ctx *gin.Context) {
	//1.获取参数和参数校验
	//NewResponse就是包装一下ctx
	reply := app.NewResponse(ctx)
	params := new(request.ParamAcceptApplication)
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
	err := logic.Logics.Application.AcceptApplication(ctx, params.AccountID, content.ID, params.CreatAt)

	//3.返回响应
	reply.Reply(err)
}

func (application) RefuseApplication(ctx *gin.Context) {
	//1.获取参数和参数校验
	//NewResponse就是包装一下ctx
	reply := app.NewResponse(ctx)
	params := new(request.ParamRefuseApplication)
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
	// 被申请人拒绝申请人发的申请
	err := logic.Logics.Application.RefuseApplication(ctx, params.AccountID, content.ID, params.RefuseMsg, params.CreatAt)

	//3.返回响应
	reply.Reply(err)
}

func (application) ListApplications(ctx *gin.Context) {
	//1.获取参数和参数校验
	//这个不需要传参数，只要token即可
	reply := app.NewResponse(ctx)
	//获取账号Token信息
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	//2.业务处理
	limit, offset := global.Page.GetPageSizeAndOffset(ctx.Request)

	result, err := logic.Logics.Application.ListApplications(ctx, content.ID, limit, offset)

	//3.返回响应
	reply.Reply(err, result)
}
