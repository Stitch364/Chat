package logic

import (
	"chat/dao"
	"chat/errcodes"
	"chat/global"
	"chat/middlewares"
	"chat/model/reply"
	"chat/pkg/emailMark"
	"errors"
	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	"github.com/XYYSWK/Lutils/pkg/utils"
	"github.com/gin-gonic/gin"
)

type email struct {
}

func (email) ExistEmail(ctx *gin.Context, emailStr string) (*reply.ParamExistEmail, errcode.Err) {
	// 先在redis中找
	ok, err := dao.Database.Redis.ExistEmail(ctx, emailStr)
	if ok {
		// ok == true 邮箱存在，在redis中找到的
		return &reply.ParamExistEmail{Exist: ok}, nil
	}
	// 记录错误日志
	if err != nil {
		global.Logger.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
	}
	// redis中没找到再去mysql中找
	ok, err = dao.Database.DB.ExistEmail(ctx, emailStr)
	if err != nil {
		global.Logger.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	return &reply.ParamExistEmail{Exist: ok}, nil
}

// CheckEmailNotExist 判断邮箱是否存在
// 先查缓存，缓存中没有再去查数据库，如果存在，将邮箱写入缓存中，返回邮箱已注册的错误
// 根据ExistEmail方法返回情况，做不同的操作
func CheckEmailNotExist(ctx *gin.Context, emailStr string) errcode.Err {
	// 调用上面的ExistEmail
	result, err := email{}.ExistEmail(ctx, emailStr)
	if err != nil {
		return err
	}
	if result.Exist {
		// 说明邮箱已经注册
		return errcodes.EmailExists
	}
	// result.Exist == false 时表明redis和mysql中都没有
	// 未找到邮箱信息（包括 缓存中没有数据库中有，缓存中和数据库中都没有）
	// 查数据库
	exist, myErr := dao.Database.DB.ExistEmail(ctx, emailStr)
	if myErr != nil {
		global.Logger.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	// 数据库有的话写进缓存
	if exist {
		//添加进缓存
		myErr = dao.Database.Redis.AddEmails(ctx, emailStr)
		if myErr != nil {
			global.Logger.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
			return errcode.ErrServer
		}
		return errcodes.EmailExists
	}
	// mysql数据库中没有就是真正的没有
	return nil
}

// SendCode 发送验证码(邮件)
func (email) SendCode(emailStr string) errcode.Err {
	// 判断发送邮件的频率
	if global.EmailMark.CheckUserExist(emailStr) {
		return errcodes.EmailSendMany
	}
	// 异步发送邮件(使用工作池)
	global.Worker.SendTask(func() {
		//utils.RandomString生成一个长度为n的随机字符串
		//生成验证码
		code := utils.RandomString(global.PublicSetting.Rules.CodeLength)
		if err := global.EmailMark.SendMark(emailStr, code); err != nil && !errors.Is(err, emailMark.ErrSendTooMany) {
			global.Logger.Error(err.Error())
		}
	})
	return nil
}
