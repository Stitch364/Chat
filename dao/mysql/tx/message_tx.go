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
