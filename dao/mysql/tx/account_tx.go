package tx

import (
	db "chat/dao/mysql/sqlc"
	"chat/dao/redis/operate"
	"chat/pkg/tool"
	"context"
	"database/sql"
	"errors"
)

var (
	ErrAccountOverNum     = errors.New("账户数量超过限制")
	ErrAccountNameExists  = errors.New("账户名已存在")
	ErrAccountGroupLeader = errors.New("账户是群主")
)

// CreateAccountWithTx 检查账户数量，账户名后创建和自己的关系
func (store *MySQLDB) CreateAccountWithTx(ctx context.Context, rdb *operate.RDB, maxAccountNum int64, arg *db.CreateAccountParams) error {
	return store.execTx(ctx, func(queries *db.Queries) error {
		var err error
		var accountNum int64
		//检查数量
		err = tool.DoThat(err, func() error {
			accountNum, err = queries.CountAccountsByUserID(ctx, arg.UserID)
			return err
		})
		if accountNum >= maxAccountNum {
			return ErrAccountOverNum
		}
		//检查账户名
		var exist bool
		err = tool.DoThat(err, func() error {
			exist, err = queries.ExistsAccountByNameAndUserID(ctx, &db.ExistsAccountByNameAndUserIDParams{
				Name:   arg.Name,
				UserID: arg.UserID,
			})
			return err
		})
		if exist {
			return ErrAccountNameExists
		}
		// 数量未超限制并且用户名合法
		// 创建账户
		err = tool.DoThat(err, func() error {
			err = queries.CreateAccount(ctx, arg)
			return err
		})

		//建立与自己的关系（就是自己和自己的好友关系）
		var relationID int64
		err = tool.DoThat(err, func() error {
			err = queries.CreateFriendRelation(ctx, &db.CreateFriendRelationParams{
				Account1ID: sql.NullInt64{Int64: arg.ID, Valid: true},
				Account2ID: sql.NullInt64{Int64: arg.ID, Valid: true},
			})
			err := tool.DoThat(err, func() error {
				relationID, err = queries.GetFriendRelationIdByID1AndID1(ctx, &db.GetFriendRelationIdByID1AndID1Params{
					Account1ID: sql.NullInt64{Int64: arg.ID, Valid: true},
					Account2ID: sql.NullInt64{Int64: arg.ID, Valid: true},
				})
				//err是查询的err
				return err
			})
			return err
		})
		err = tool.DoThat(err, func() error {
			return queries.CreateSetting(ctx, &db.CreateSettingParams{
				AccountID:  arg.ID,
				RelationID: relationID,
				IsSelf:     true,
			})
		})
		// 添加自己一个人的关系到 redis
		err = tool.DoThat(err, func() error {
			return rdb.AddRelationAccount(ctx, relationID, arg.ID)
		})
		return err
	})
}

func (store *MySQLDB) DeleteAccountWithTx(ctx context.Context, rdb *operate.RDB, accountID int64) error {
	return store.execTx(ctx, func(queries *db.Queries) error {
		//var err error
		//判断用户是否是群主

		//删除好友

		//删除群

		//删除账户

		//从redis中删除关系

		//从redis中删除账户所在的群中的该账户
		return nil
	})
}
