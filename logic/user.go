package logic

import (
	"chat/dao"
	db "chat/dao/mysql/sqlc"
	"chat/errcodes"
	"chat/global"
	"chat/middlewares"
	"chat/model"
	"chat/model/reply"
	"database/sql"
	"errors"
	"fmt"
	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	"github.com/XYYSWK/Lutils/pkg/password"
	"github.com/gin-gonic/gin"
)

type user struct{}

func (user) Register(ctx *gin.Context, emailStr, pwd, code string) (*reply.ParamRegister, errcode.Err) {
	//判断邮箱是否注册过
	if err := CheckEmailNotExist(ctx, emailStr); err != nil {
		fmt.Println("err:", err.Error())
		return nil, err
	}

	//校验验证码
	if !global.EmailMark.CheckCode(emailStr, code) {
		return nil, errcodes.EmailCodeNotValid
	}
	//密码加密
	hashPassword, err := password.HashPassword(pwd)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	//将 user 写入数据并返回 UserInfo
	//只返回一个err
	err = dao.Database.DB.CreateUser(ctx, &db.CreateUserParams{
		Email:    emailStr,
		Password: hashPassword,
	})

	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	//重新通过邮箱获取用户信息
	userInfo, err := dao.Database.DB.GetUserByEmail(ctx, emailStr)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	//添加邮箱进 redis
	err = dao.Database.Redis.AddEmails(ctx, emailStr)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	//创建 Token

	//accessToken
	accessToken, accessPayload, err := newUserToken(model.UserToken, userInfo.ID, global.PrivateSetting.Token.AccessTokenExpire)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	//refreshToken
	refreshToken, _, err := newUserToken(model.UserToken, userInfo.ID, global.PrivateSetting.Token.RefreshTokenExpire)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	//将token保存进redis
	if err = dao.Database.Redis.SaveUserToken(ctx, userInfo.ID, []string{accessToken, refreshToken}); err != nil {
		return nil, errcode.ErrServer.WithDetails(err.Error())
	}
	//返回数据以及错误信息
	return &reply.ParamRegister{
		ParamUserInfo: reply.ParamUserInfo{
			ID:       userInfo.ID,
			Email:    userInfo.Email,
			CreateAt: userInfo.CreatedAt,
		},
		Token: reply.ParamToken{
			AccessToken:   accessToken,
			AccessPayload: accessPayload,
			RefreshToken:  refreshToken,
		},
	}, nil
}

func (user) Login(ctx *gin.Context, emailStr, pwd string) (*reply.ParamLogin, errcode.Err) {
	//校验账号是否存在
	//通过用户邮箱获取用户信息
	userInfo, err := dao.Database.DB.GetUserByEmail(ctx, emailStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errcodes.UserNotFound
		}
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	//校验密码
	if err := password.CheckPassword(pwd, userInfo.Password); err != nil {
		return nil, errcodes.PasswordNotValid
	}
	//创建新的Token
	//accessToken
	accessToken, accessPayload, err := newUserToken(model.UserToken, userInfo.ID, global.PrivateSetting.Token.AccessTokenExpire)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	//refreshToken
	refreshToken, _, err := newUserToken(model.UserToken, userInfo.ID, global.PrivateSetting.Token.RefreshTokenExpire)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	//将token保存进redis
	if err = dao.Database.Redis.SaveUserToken(ctx, userInfo.ID, []string{accessToken, refreshToken}); err != nil {
		return nil, errcode.ErrServer.WithDetails(err.Error())
	}
	return &reply.ParamLogin{
		ParamUserInfo: reply.ParamUserInfo{
			ID:       userInfo.ID,
			Email:    userInfo.Email,
			CreateAt: userInfo.CreatedAt,
		},
		Token: reply.ParamToken{
			AccessToken:   accessToken,
			AccessPayload: accessPayload,
			RefreshToken:  refreshToken,
		},
	}, nil
}

// UpdateUserPassword 更新用户密码
func (user) UpdateUserPassword(ctx *gin.Context, userID int64, code, newPwd string) errcode.Err {
	userInfo, myerr := getUserInfoByID(ctx, userID)
	if myerr != nil {
		global.Logger.Error(myerr.Error(), middlewares.ErrLogMsg(ctx)...)
		return myerr
	}
	//校验验证码
	if !global.EmailMark.CheckCode(userInfo.Email, code) {
		return errcodes.EmailCodeNotValid
	}
	//密码加密
	hashPassword, err := password.HashPassword(newPwd)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	//保存密码
	if err := dao.Database.DB.UpdateUser(ctx, &db.UpdateUserParams{
		Email:    userInfo.Email,
		ID:       userID,
		Password: hashPassword,
	}); err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}

	//清除Token（改密码后要成新登录）
	if err := dao.Database.Redis.DeleteAllTokenByUser(ctx, userID); err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer.WithDetails(err.Error())
	}
	return nil
}
func (user) UpdateUserEmail(ctx *gin.Context, userID int64, code, newEmail string) errcode.Err {
	userInfo, myerr := getUserInfoByID(ctx, userID)
	if myerr != nil {
		global.Logger.Error(myerr.Error(), middlewares.ErrLogMsg(ctx)...)
		return myerr
	}
	// 邮箱不能和之前的邮箱重复
	if userInfo.Email == newEmail {
		return errcodes.EmailSame
	}
	// 判断邮箱是否已经注册过
	if err := CheckEmailNotExist(ctx, newEmail); err != nil {
		return err
	}

	//校验验证码
	if !global.EmailMark.CheckCode(newEmail, code) {
		return errcodes.EmailCodeNotValid
	}
	//保存邮箱
	if err := dao.Database.DB.UpdateUser(ctx, &db.UpdateUserParams{
		Email:    newEmail,
		ID:       userID,
		Password: userInfo.Password,
	}); err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	//更新 redis 中的邮箱
	if err := dao.Database.Redis.UpdateEmail(ctx, userInfo.Email, newEmail); err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	/*	// 推送更改邮箱通知
		accessToken, _ := middlewares.GetToken(ctx.Request.Header)
		global.Worker.SendTask(task.UpdateEmail(accessToken, userID, emailStr))*/
	return nil
}

// Logout 退出登录
func (user) Logout(ctx *gin.Context) errcode.Err {
	Token, payload, err := GetTokenAndPayload(ctx)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcodes.AuthenticationFailed
	}
	content := &model.Content{}
	_ = content.Unmarshal(payload.Content)
	// 判断用户在redis中是否存在
	if ok := dao.Database.Redis.CheckUserTokenValid(ctx, content.ID, Token); !ok {
		return errcodes.UserNotFound
	}
	//先将token从redis中清除
	if err := dao.Database.Redis.DeleteAllTokenByUser(ctx, content.ID); err != nil {
		return errcode.ErrServer.WithDetails(err.Error())
	}
	return nil
}

// 通过ID获取用户信息
func getUserInfoByID(ctx *gin.Context, userID int64) (*db.User, errcode.Err) {
	userInfo, err := dao.Database.DB.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errcodes.UserNotFound
		}
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	return userInfo, nil
}

// 通过邮箱获取用户信息
func getUserInfoByEmail(ctx *gin.Context, emailStr string) (*db.User, errcode.Err) {
	userInfo, err := dao.Database.DB.GetUserByEmail(ctx, emailStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errcodes.UserNotFound
		}
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	return userInfo, nil
}
