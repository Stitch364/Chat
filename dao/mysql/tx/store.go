package tx

import (
	db "chat/dao/mysql/sqlc"
	"chat/dao/redis/operate"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // 导入 MySQL 驱动
)

// TXer 用于执行事务相关操作
type TXer interface {
	// CreateAccountWithTx 创建账号并建立和自己的关系
	CreateAccountWithTx(ctx context.Context, rdb *operate.RDB, maxAccountNum int64, arg *db.CreateAccountParams) error
	DeleteAccountWithTx(ctx context.Context, rdb *operate.RDB, accountID int64) error
	CreateApplicationTx(ctx context.Context, params *db.CreateApplicationParams) error
	AcceptApplicationTx(ctx context.Context, rdb *operate.RDB, account1, account2 *db.GetAccountByIDRow) (*db.Message, error)
	CreateMessageTx(ctx context.Context, params *db.CreateMessageParams) (*db.GetMessageInfoTxRow, error)
	RevokeMsgWithTx(ctx context.Context, msgID int64, isPin, isTop bool) error
}

// MySQLDB 实现了 DB 接口
// 直接嵌入实现了DB接口的db.Queries，就相当于MySQLDB也实现了
type MySQLDB struct {
	*db.Queries
	Db *sql.DB
}

// 通过事务执行回调函数
func (store *MySQLDB) execTx(ctx context.Context, fn func(tx *db.Queries) error) error {
	// 开启一个数据库事务
	tx, err := store.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := store.WithTx(tx) // 使用开启的事务创建一个查询
	// 调用传入的回调函数执行数据库操作
	if err := fn(q); err != nil {
		// 如果回调函数执行失败，回滚事务
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err:%v, rb err:%v", err, rbErr)
		}
		return err
	}
	// 提交事务
	return tx.Commit()
}
