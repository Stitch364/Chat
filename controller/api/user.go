package api

import (
	"chat/errcodes"
	"chat/logic"
	"chat/middlewares"
	"chat/model"
	"chat/model/request"
	"github.com/XYYSWK/Lutils/pkg/app"
	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
)

type user struct{}

func (user) Register(ctx *gin.Context) {
	//1.获取参数和参数校验
	//NewResponse就是包装一下ctx
	reply := app.NewResponse(ctx)
	params := new(request.ParamRegister)
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	//2.业务处理
	result, err := logic.Logics.User.Register(ctx, params.Email, params.Password, params.Code)
	//3.返回响应
	reply.Reply(err, result)
}

func (user) Login(ctx *gin.Context) {
	//1.获取参数和参数校验
	//NewResponse就是包装一下ctx
	reply := app.NewResponse(ctx)
	params := new(request.ParamLogin)
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}

	//2.业务处理
	result, err := logic.Logics.User.Login(ctx, params.Email, params.Password)
	//3.返回响应
	reply.Reply(err, result)
}
func (user) UpdateUserPassword(ctx *gin.Context) {
	//1.获取参数和参数校验
	//NewResponse就是包装一下ctx
	reply := app.NewResponse(ctx)
	params := new(request.ParamUpdateUserPassword)
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	//从context中拿 userID
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.UserToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	//2.业务处理
	err := logic.Logics.User.UpdateUserPassword(ctx, content.ID, params.Code, params.NewPassword)
	//3.返回响应
	reply.Reply(err)
}

func (user) UpdateUserEmail(ctx *gin.Context) {
	//1.获取参数和参数校验
	//NewResponse就是包装一下ctx
	reply := app.NewResponse(ctx)
	params := new(request.ParamUpdateUserEmail)
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	//从context中拿 userID
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.UserToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	//2.业务处理
	err := logic.Logics.User.UpdateUserEmail(ctx, content.ID, params.Code, params.Email)
	//3.返回响应
	reply.Reply(err)

}

func (user) Logout(ctx *gin.Context) {
	//1.获取参数和参数校验
	//NewResponse就是包装一下ctx
	reply := app.NewResponse(ctx)

	//2.业务处理
	if err := logic.Logics.User.Logout(ctx); err != nil {
		reply.Reply(err)
		return
	}
	//3.返回响应
	reply.Reply(nil, gin.H{
		"msg": "登出成功",
	})
}
