package tx

import (
	db "chat/dao/mysql/sqlc"
	"chat/pkg/tool"
	"context"
)

func (store *MySQLDB) CreateMessageTx(ctx context.Context, params *db.CreateMessageParams) (*db.GetMessageInfoTxRow, error) {
	var msg *db.GetMessageInfoTxRow
	//用于实现创建消息并返回新建消息并返回id等信息
	err := store.execTx(ctx, func(q *db.Queries) error {
		err := q.CreateMessage(ctx, params)
		err = tool.DoThat(err, func() error {
			//获取刚创建的消息的信息
			msg, err = q.GetMessageInfoTx(ctx)
			return err
		})
		return err
	})
	return msg, err
}

// RevokeMsgWithTx 撤回消息，如果消息 pin 或者置顶，则全部取消
func (store *MySQLDB) RevokeMsgWithTx(ctx context.Context, msgID int64, isPin, isTop bool) error {
	return store.execTx(ctx, func(queries *db.Queries) error {
		var err error
		err = tool.DoThat(err, func() error {
			return queries.UpdateMsgRevoke(ctx, &db.UpdateMsgRevokeParams{
				ID:       msgID,
				IsRevoke: true,
			})
		})
		if isPin {
			err = tool.DoThat(err, func() error {
				return queries.UpdateMsgPin(ctx, &db.UpdateMsgPinParams{
					ID:    msgID,
					IsPin: false,
				})
			})
		}
		if isTop {
			err = tool.DoThat(err, func() error {
				return queries.UpdateMsgTop(ctx, &db.UpdateMsgTopParams{
					ID:    msgID,
					IsTop: false,
				})
			})
		}
		return err
	})
}

func (store *MySQLDB) DeleteMsgWithTx(ctx context.Context, msgID int64, isDel int32) error {
	return store.execTx(ctx, func(queries *db.Queries) error {
		var err error
		err = tool.DoThat(err, func() error {
			IsDelete, err := queries.GetMsgDeleteById(ctx, msgID)
			if IsDelete == 0 {
				err = tool.DoThat(err, func() error {
					return queries.UpdateMsgDelete(ctx, &db.UpdateMsgDeleteParams{
						ID:       msgID,
						IsDelete: isDel,
					})
				})
				return err
			} else {
				err = tool.DoThat(err, func() error {
					return queries.UpdateMsgDelete(ctx, &db.UpdateMsgDeleteParams{
						ID:       msgID,
						IsDelete: 2,
					})
				})
				return err
			}
		})
		return err
	})
}
