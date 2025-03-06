package redis

import (
	"fmt"
	"web_app/setting"

	"github.com/go-redis/redis"
)

// 声明一个全局的rdb变量
var rdb *redis.Client

// Init 初始化连接
// 普通连接
func Init(cfg *setting.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host, cfg.Port,
			//viper.GetString("redis.host"),
			//viper.GetString("redis.port"),
		),
		Password: cfg.Password, //viper.GetString("redis.password"), // 密码
		DB:       cfg.Database, //viper.GetInt("redis.db"),          // 数据库
		PoolSize: cfg.PollSize, //viper.GetInt("redis.pool_size"),   // 连接池大小
	})

	//接收Ping的结果
	_, err = rdb.Ping().Result()
	return err
}

func Close() {
	_ = rdb.Close()
}
