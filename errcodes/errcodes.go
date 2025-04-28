package errcodes

import "github.com/XYYSWK/Lutils/pkg/app/errcode"

var (
	UserNotFound                = errcode.NewErr(2001, "用户不存在")
	PasswordNotValid            = errcode.NewErr(2002, "密码错误")
	EmailSendMany               = errcode.NewErr(2003, "邮件发送频繁，请稍后再试")
	EmailCodeNotValid           = errcode.NewErr(2004, "邮箱验证码校验失败")
	EmailSame                   = errcode.NewErr(2005, "邮箱相同")
	EmailExists                 = errcode.NewErr(2006, "邮箱已经注册")
	AuthNotExist                = errcode.NewErr(2007, "身份不存在") // 没找到header上的token
	AuthenticationFailed        = errcode.NewErr(2008, "身份验证失败")
	AuthPermissionsInsufficient = errcode.NewErr(2009, "权限不足") // 对不属于自己的资源进行操作
	AccountNotFound             = errcode.NewErr(2010, "账号不存在")
	AccountNameExists           = errcode.NewErr(2011, "账号名已经存在")
	AccountNumExcessive         = errcode.NewErr(2012, "账号数量超过限制")
	AuthOverTime                = errcode.NewErr(2013, "身份过期")
	ApplicationExists           = errcode.NewErr(3001, "申请已经存在")
	ApplicationNotExists        = errcode.NewErr(3002, "申请不存在")
	ApplicationNotValid         = errcode.NewErr(3003, "申请不合法")
	ApplicationRepeatOpt        = errcode.NewErr(3004, "重复操作申请")
	UserHasAccount              = errcode.NewErr(3005, "用户存在账户")     // 删除用户时用户仍存在账户
	AccountGroupLeader          = errcode.NewErr(3006, "账号作为群主不能删除") // 删除账户时，账户存在群主身份
)
