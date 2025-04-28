package tx

import (
	db "chat/dao/mysql/sqlc"
	"chat/dao/redis/operate"
	"chat/errcodes"
	"chat/model"
	"chat/pkg/tool"
	"context"
	"database/sql"
	"encoding/json"
)

// CreateApplicationTx 用事务先判断是否存在申请，不存在则创建
// 用户1给用户2发送申请，2在没回应的情况下，是不能给用户1发送申请的,因为这里查询申请存在时查的是双向的申请
// 开事务是为了防止查询时，更改了数据库中申请的信息，有可能导致用户2已经同意了申请，用户1再发送申请
func (store *MySQLDB) CreateApplicationTx(ctx context.Context, params *db.CreateApplicationParams) error {
	return store.execTx(ctx, func(q *db.Queries) error {
		//查看申请是否存在
		ok, err := q.ExistsApplicationByIDWithLock(ctx, &db.ExistsApplicationByIDWithLockParams{
			Account1ID:   params.Account1ID,
			Account1ID_2: params.Account1ID,
		})
		if err != nil {
			return err
		}
		if ok {
			//存在申请
			return errcodes.ApplicationExists
		}
		return q.CreateApplication(ctx, params)
	})
}

func (store *MySQLDB) AcceptApplicationTx(ctx context.Context, rdb *operate.RDB, account1, account2 *db.GetAccountByIDRow) (*db.Message, error) {
	var result *db.Message
	err := store.execTx(ctx, func(q *db.Queries) error {
		var err error
		//修改申请状态
		err = tool.DoThat(err, func() error {
			return q.UpdateApplication(ctx, &db.UpdateApplicationParams{
				Account1ID: account1.ID,
				Account2ID: account2.ID,
				Status:     db.ApplicationsStatusValue1, //已同意
			})
		})
		id1, id2 := account1.ID, account2.ID
		if id1 > id2 {
			id1, id2 = id2, id1
		}
		//创建好友关系
		var relationID int64
		err = tool.DoThat(err, func() error {
			err = q.CreateFriendRelation(ctx, &db.CreateFriendRelationParams{
				Account1ID: sql.NullInt64{Int64: id1, Valid: true},
				Account2ID: sql.NullInt64{Int64: id2, Valid: true},
			})
			//查询刚创建的好友关系ID
			err = tool.DoThat(err, func() error {
				relationID, err = q.GetFriendRelationIdByID1AndID1(ctx, &db.GetFriendRelationIdByID1AndID1Params{
					Account1ID: sql.NullInt64{Int64: id1, Valid: true},
					Account2ID: sql.NullInt64{Int64: id2, Valid: true},
				})
				return err
			})
			return err
		})
		//创建好友关系设置(双方都要有)
		err = tool.DoThat(err, func() error {
			return q.CreateSetting(ctx, &db.CreateSettingParams{
				AccountID:  account1.ID,
				RelationID: relationID,
				IsLeader:   false,
				IsSelf:     false,
			})
		})
		err = tool.DoThat(err, func() error {
			return q.CreateSetting(ctx, &db.CreateSettingParams{
				AccountID:  account2.ID,
				RelationID: relationID,
				IsLeader:   false,
				IsSelf:     false,
			})
		})
		//发送通知信息
		err = tool.DoThat(err, func() error {
			arg := &db.CreateMessageParams{
				NotifyType: db.MessagesNotifyTypeCommon,
				MsgType:    db.MessagesMsgType(model.MsgTypeText),
				MsgContent: "我们已经是好友了，现在可以开始聊天啦！",
				MsgExtend:  json.RawMessage(``),
				AccountID:  sql.NullInt64{Int64: account1.ID, Valid: true},
				RelationID: relationID,
			}
			//创建消息
			err = q.CreateMessage(ctx, arg)
			var MsgInfo *db.GetMessageInfoTxRow
			err = tool.DoThat(err, func() error {
				//获取刚创建的消息的信息
				MsgInfo, err = q.GetMessageInfoTx(ctx)
				return err
			})

			result = &db.Message{
				ID:         MsgInfo.ID,
				NotifyType: arg.NotifyType,
				MsgType:    arg.MsgType,
				MsgContent: arg.MsgContent,
				RelationID: relationID,
				CreateAt:   MsgInfo.CreateAt,
			}
			return err
		})
		//保存好友关系到redis
		err = tool.DoThat(err, func() error {
			return rdb.AddRelationAccount(ctx, relationID, account1.ID, account2.ID)
		})
		return err
	})
	return result, err
}
