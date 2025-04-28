package mysql

import (
	db "chat/dao/mysql/sqlc"
	"chat/dao/mysql/tx"
	"chat/model/config"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type DB interface {
	db.Querier //查询接口
	tx.TXer    //事务接口
}

func Init(cfg config.MysqlConfig) DB { //使用全局结构体获取配置信息
	// 也可以使用MustConnect,连接不成功就panic,不返回错误
	Db, err := sql.Open(cfg.DirverName, cfg.DataSourceName)
	if err != nil {
		//直接用zap.L()
		//日志记录错误
		//zap.L().Error("connect DB failed", zap.Error(err))
		panic(err)
	}
	//设置最大连接数
	Db.SetMaxOpenConns(cfg.MaxOpenConns)
	//设置最大空闲连接数
	Db.SetMaxIdleConns(cfg.MaxIdleConns)
	// 设置连接的最大复用时间，超过该时间后连接将被关闭并重新创建
	Db.SetConnMaxLifetime(time.Minute * 10)
	return &tx.MySQLDB{Queries: db.New(Db), Db: Db}
}
