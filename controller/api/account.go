package api

import (
	"chat/errcodes"
	"chat/global"
	"chat/logic"
	"chat/middlewares"
	"chat/model"
	"chat/model/request"
	"fmt"
	"github.com/XYYSWK/Lutils/pkg/app"
	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
)

type account struct {
}

func (account) CreateAccount(ctx *gin.Context) {
	//1.获取参数和参数校验
	//NewResponse就是包装一下ctx
	reply := app.NewResponse(ctx)
	params := new(request.ParamCreateAccount)
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	//获取Token信息
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok && content.TokenType != model.UserToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	//2.业务处理
	result, err := logic.Logics.Account.CreateAccount(ctx, content.ID, params.Name, global.PublicSetting.Rules.DefaultAvatarURL, params.Gender, params.Signature)
	//3.返回响应
	reply.Reply(err, result)
}

// GetAccountToken 获取账号Token
func (account) GetAccountToken(ctx *gin.Context) {
	//1.获取参数和参数校验
	//NewResponse就是包装一下ctx
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetAccountToken)
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	//获取UserToken信息
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok && content.TokenType != model.UserToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	//2.业务处理
	result, err := logic.Logics.Account.GetAccountToken(ctx, content.ID, params.AccountID)

	//3.返回响应
	reply.Reply(err, result)
}

func (account) DeleteAccount(ctx *gin.Context) {
	//1.获取参数和参数校验
	//NewResponse就是包装一下ctx
	reply := app.NewResponse(ctx)
	params := new(request.ParamDeleteAccount)
	if err := ctx.ShouldBind(params); err != nil {
		fmt.Println(err)
		fmt.Println(params.AccountID)
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	//fmt.Println("1 on ctrl")
	//获取UserToken信息
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok && content.TokenType != model.UserToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	//fmt.Println("2 on ctrl")
	//2.业务处理
	err := logic.Logics.Account.DeleteAccount(ctx, content.ID, params.AccountID)

	//3.返回响应
	reply.Reply(err)
}

func (account) GetAccountsByUserID(ctx *gin.Context) {
	//1.获取参数和参数校验
	//NewResponse就是包装一下ctx
	reply := app.NewResponse(ctx)
	//直接从登陆的用户中获取UserID
	//获取UserToken信息
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok && content.TokenType != model.UserToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	//global.Logger.Info("000")

	//2.业务处理
	result, err := logic.Logics.Account.GetAccountsByUserID(ctx, content.ID)

	//3.返回响应
	reply.Reply(err, result)
}

func (account) GetAccountByID(ctx *gin.Context) {
	//1.获取参数和参数校验
	//NewResponse就是包装一下ctx
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetAccountByID)
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	//获取UserToken信息
	//这里的Token是账号的Token不是用户的Token
	//所以Token ID是账号的ID
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	//2.业务处理
	result, err := logic.Logics.Account.GetAccountByID(ctx, params.AccountID, content.ID)

	//3.返回响应
	reply.Reply(err, result)
}

func (account) UpdateAccount(ctx *gin.Context) {
	//1.获取参数和参数校验
	//NewResponse就是包装一下ctx
	reply := app.NewResponse(ctx)
	params := new(request.ParamUpdateAccount)
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	//获取UserToken信息
	//这里的Token是账号的Token不是用户的Token
	//所以Token ID是账号的ID
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	//2.业务处理
	err := logic.Logics.Account.UpdateAccount(ctx, content.ID, params.Name, params.Gender, params.Signature)

	//3.返回响应
	reply.Reply(err)
}

func (account) GetAccountsByName(ctx *gin.Context) {
	//1.获取参数和参数校验
	//NewResponse就是包装一下ctx
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetAccountsByName)
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	//获取UserToken信息
	//这里的Token是账号的Token不是用户的Token
	//所以Token ID是账号的ID
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	limit, offset := global.Page.GetPageSizeAndOffset(ctx.Request)

	//2.业务处理
	result, err := logic.Logics.Account.GetAccountsByName(ctx, content.ID, params.Name, limit, offset)

	//3.返回响应
	reply.Reply(err, result)
}
