package api

import (
	"chat/global"
	"chat/logic"
	"chat/model/request"
	"github.com/XYYSWK/Lutils/pkg/app"
	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type email struct{}

//先验证邮箱是否存在，若存在才能发送验证码

func (email) ExistEmail(ctx *gin.Context) {
	//包装ctx
	reply := app.NewResponse(ctx)
	//创建request结构体
	params := &request.ParamExistEmail{}
	//参数绑定，错误处理
	if err := ctx.ShouldBindQuery(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	//几乎是固定的格式了
	//调用logic层的方法处理业务
	result, err := logic.Logics.Email.ExistEmail(ctx, params.Email)
	//返回结果
	reply.Reply(err, result)
}

// SendCode 发送验证码
func (email) SendCode(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := &request.ParamSendEmail{}
	if err := ctx.ShouldBind(params); err != nil {
		global.Logger.Error("err", zap.Error(err))
		//返回错误的详细信息
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	err := logic.Logics.Email.SendCode(params.Email)
	reply.Reply(err)
}
