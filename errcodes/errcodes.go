package errcodes

import "github.com/XYYSWK/Lutils/pkg/app/errcode"

const SevenDay = int(7 * 24 * 60)

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
	DeleteNotLegal              = errcode.NewErr(3007, "不能删除自己")
	CoolingOffPeriod            = errcode.NewErr(3008, "7天冷却期内不可发送申请")
	RelationNotExists           = errcode.NewErr(4002, "关系不存在")
	MsgNotExists                = errcode.NewErr(5001, "消息不存在")
	MsgAlreadyRevoke            = errcode.NewErr(5002, "消息已经撤销")
	RlyMsgHasRevoked            = errcode.NewErr(5003, "回复消息已经撤销")
	RlyMsgNotOneRelation        = errcode.NewErr(5004, "回复的消息和发送的消息并非在一个群")
	MsgRevokeTimeOut            = errcode.NewErr(5006, "消息撤销时间已过")
	NotLeader                   = errcode.NewErr(7001, "非群主")
	IsLeader                    = errcode.NewErr(7002, "群主不可退群")
	NotGroupMember              = errcode.NewErr(7003, "非该群成员")
	FileNotExist                = errcode.NewErr(8002, "文件不存在")
	FileTooBig                  = errcode.NewErr(8004, "文件过大")
	FileDeleteFailed            = errcode.NewErr(8003, "文件删除失败")
	FileIsEmpty                 = errcode.NewErr(8005, "文件为空")
	//FailedStore                 = errcode.NewErr(8001, "存储文件失败")
	//NotifyNotExist              = errcode.NewErr(9001, "通知不存在")
)
