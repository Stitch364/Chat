package logic

import (
	"chat/dao"
	db "chat/dao/mysql/sqlc"
	"chat/dao/mysql/tx"
	"chat/errcodes"
	"chat/global"
	"chat/middlewares"
	"chat/model"
	"chat/model/common"
	"chat/model/reply"
	"chat/task"
	"database/sql"
	"errors"
	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
)

type account struct{}

func (account) CreateAccount(ctx *gin.Context, userID int64, name, avatar, gender, signature string) (*reply.ParamCreateAccount, errcode.Err) {
	arg := &db.CreateAccountParams{
		ID:        global.GenerateID.GetID(),
		UserID:    userID,
		Name:      name,
		Avatar:    avatar,
		Gender:    db.AccountsGender(gender),
		Signature: signature,
	}
	//创建账户和自己的关系
	err := dao.Database.DB.CreateAccountWithTx(ctx, dao.Database.Redis, global.PublicSetting.Rules.AccountNumMax, arg)
	//处理错误
	switch {
	case errors.Is(err, tx.ErrAccountOverNum):
		return nil, errcodes.AccountNumExcessive
	case errors.Is(err, tx.ErrAccountNameExists):
		return nil, errcodes.AccountNameExists
	case err == nil:
	default:
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	//创建token
	token, payload, err := newAccountToken(model.AccountToken, arg.ID)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	return &reply.ParamCreateAccount{
		ParamAccountInfo: reply.ParamAccountInfo{
			ID:     arg.ID,
			Name:   arg.Name,
			Avatar: arg.Avatar,
			Gender: string(arg.Gender),
		},
		ParamGetAccountToken: reply.ParamGetAccountToken{
			AccountToken: common.Token{
				Token:    token,
				ExpireAt: payload.ExpiredAt, //过期时间
			},
		},
	}, nil
}

// GetAccountToken 生成账号Token
func (account) GetAccountToken(ctx *gin.Context, userID, accountID int64) (*reply.ParamGetAccountToken, errcode.Err) {
	//用账号ID获取账号信息
	accountInfo, myerr := getAccountInfoByID(ctx, accountID, accountID)
	if myerr != nil {
		return nil, myerr
	}
	//账号所属用户ID与当前用户ID对比
	if accountInfo.UserID != userID {
		return nil, errcodes.AuthPermissionsInsufficient
	}
	//给账号创建Token
	token, payload, err := newAccountToken(model.AccountToken, accountID)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	return &reply.ParamGetAccountToken{AccountToken: common.Token{
		Token:    token,
		ExpireAt: payload.ExpiredAt,
	}}, nil
}

// 通过账号 ID 获取账号信息
func getAccountInfoByID(ctx *gin.Context, accountID, selfID int64) (*db.GetAccountByIDRow, errcode.Err) {
	accountInfo, err := dao.Database.DB.GetAccountByID(ctx, &db.GetAccountByIDParams{
		ID:         accountID,
		Account2ID: sql.NullInt64{Int64: selfID, Valid: true},
		Account1ID: sql.NullInt64{Int64: selfID, Valid: true},
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errcodes.AccountNotFound
		}
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	return accountInfo, nil
}

func (account) DeleteAccount(ctx *gin.Context, userID, accountID int64) errcode.Err {
	//用账号ID获取账号信息
	accountInfo, myerr := getAccountInfoByID(ctx, accountID, accountID)
	if myerr != nil {
		return myerr
	}
	//账号所属用户ID与当前用户ID对比
	if accountInfo.UserID != userID {
		return errcodes.AuthPermissionsInsufficient
	}

	//删除账号
	err := dao.Database.DB.DeleteAccountWithTx(ctx, dao.Database.Redis, accountID)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	return nil
}

func (account) GetAccountsByUserID(ctx *gin.Context, userID int64) (*reply.ParamGetAccountsByUserID, errcode.Err) {
	//global.Logger.Info("GetAccountsByUserID1")

	accounts, err := dao.Database.DB.GetAccountsByUserID(ctx, userID)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	//ParamAccountInfo内容和GetAccountsByUserIDRow内容一样，只是名字不一样
	result := make([]reply.ParamAccountInfo, len(accounts))
	for i, accountInfo := range accounts {
		result[i] = reply.ParamAccountInfo{
			ID:     accountInfo.ID,
			Name:   accountInfo.Name,
			Avatar: accountInfo.Avatar,
			Gender: string(accountInfo.Gender),
		}
	}

	return &reply.ParamGetAccountsByUserID{
		List:  result,
		Total: int64(len(result)),
	}, nil
}

// GetAccountByID 效果和getAccountInfoByID一样，不过getAccountInfoByID是为包内函数使用的
func (account) GetAccountByID(ctx *gin.Context, accountID, selfID int64) (*reply.ParamGetAccountByID, errcode.Err) {
	//用账号ID获取账号信息
	accountInfo, myerr := getAccountInfoByID(ctx, accountID, selfID)
	if myerr != nil {
		return nil, myerr
	}
	return &reply.ParamGetAccountByID{
		Info: reply.ParamAccountInfo{
			ID:     accountInfo.ID,
			Name:   accountInfo.Name,
			Avatar: accountInfo.Avatar,
			Gender: string(accountInfo.Gender),
		},
		Signature:  accountInfo.Signature,
		CreateAt:   accountInfo.CreatedAt,
		RelationID: accountInfo.RelationID.Int64,
	}, nil
}

func (account) UpdateAccount(ctx *gin.Context, accountID int64, name, gender, signature string) errcode.Err {
	err := dao.Database.DB.UpdateAccount(ctx, &db.UpdateAccountParams{
		Name:      name,
		Gender:    db.AccountsGender(gender),
		Signature: signature,
		ID:        accountID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	//获取Token
	accessToken, _ := middlewares.GetToken(ctx.Request.Header)
	//推送更新消息
	global.Worker.SendTask(task.UpdateAccount(accessToken, accountID, name, gender, signature))

	return nil
}

func (account) GetAccountsByName(ctx *gin.Context, accountID int64, name string, limit, offset int32) (*reply.ParamGetAccountsByName, errcode.Err) {
	var Total int64
	var pageTotal int64
	accounts, err := dao.Database.DB.GetAccountsByName(ctx, &db.GetAccountsByNameParams{
		CONCAT: name,
		Account2ID: sql.NullInt64{
			Int64: accountID,
			Valid: true,
		},
		Account1ID: sql.NullInt64{
			Int64: accountID,
			Valid: true,
		},
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	if len(accounts) == 0 {
		return &reply.ParamGetAccountsByName{List: []*reply.ParamFriendInfo{}}, nil
	}

	//类型转换
	accountInfos := make([]*reply.ParamFriendInfo, len(accounts))
	for i, account := range accounts {
		accountInfos[i] = &reply.ParamFriendInfo{
			ParamAccountInfo: reply.ParamAccountInfo{
				ID:     account.ID,
				Name:   account.Name,
				Avatar: account.Avatar,
				Gender: string(account.Gender),
			},
			RelationID: account.RelationID.Int64,
		}
	}

	Int64limit := int64(limit)
	Total = accounts[0].Total.(int64)
	if Total%Int64limit == 0 {
		pageTotal = Total / Int64limit
	} else {
		pageTotal = Total/Int64limit + 1
	}

	return &reply.ParamGetAccountsByName{
		List:  accountInfos,
		Total: pageTotal,
	}, nil
}
