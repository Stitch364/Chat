package mysql

import (
	"fmt"
	"web_app/setting"

	"go.uber.org/zap"

	"github.com/spf13/viper"

	"github.com/jmoiron/sqlx"
)

var Db *sqlx.DB

func Init(cfg *setting.MySQLConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		//使用全局结构体获取配置信息
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		//使用viper获取配置信息
		//viper.GetString("mysql.user"),
		//viper.GetString("mysql.password"),
		//viper.GetString("mysql.host"),
		//viper.GetInt("mysql.port"),
		//viper.GetString("mysql.dbname"),
	)
	// 也可以使用MustConnect,连接不成功就panic,不返回错误
	Db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		//直接用zap.L()
		zap.L().Error("connect DB failed", zap.Error(err))
		return
	}
	//设置最大连接数
	Db.SetMaxOpenConns(viper.GetInt("mysql.max_open_conns"))
	//设置最大空闲连接数
	Db.SetMaxIdleConns(viper.GetInt("mysql.max_idle_conns"))
	return
}

// Close 定义一个Close方法，这样当Db时db（首字母小写），即没对外暴露也可以关闭
func Close() {
	_ = Db.Close()
}
