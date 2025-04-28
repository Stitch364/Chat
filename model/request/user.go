package request

//用户请求相关的请求参数结构体

// ParamRegister 注册请求参数
type ParamRegister struct {
	Email    string `json:"email" binding:"required"`                 //邮箱
	Password string `json:"password" binding:"required,gte=6,lte=50"` //密码
	Code     string `json:"code" binding:"required,gte=6,lte=6"`      //验证码
}

// ParamLogin 登录请求参数
type ParamLogin struct {
	Email    string `json:"email" binding:"required,email,lte=50"`    // 邮箱
	Password string `json:"password" binding:"required,gte=6,lte=50"` //验证码
}

// ParamUpdateUserPassword 更新密码请求参数
type ParamUpdateUserPassword struct {
	Code        string `json:"code" binding:"required,gte=6,lte=50"`        // 验证码
	NewPassword string `json:"newPassword" binding:"required,gte=6,lte=50"` // 新密码
}

// ParamUpdateUserEmail 更新邮箱请求参数
type ParamUpdateUserEmail struct {
	Email string `json:"email" binding:"required,email,lte=50"` // 邮箱
	Code  string `json:"code" binding:"required,gte=6,lte=6"`   // 验证码
}
