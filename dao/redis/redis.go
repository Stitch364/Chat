package redis

import (
	"chat/dao/redis/operate"
	"context"
	"github.com/go-redis/redis/v8"
)

// Init 初始化连接
// 普通连接
//func Init(cfg *setting.RedisConfig) (err error) {
//	rdb = redis.NewClient(&redis.Options{
//		Addr: fmt.Sprintf("%s:%d",
//			cfg.Host, cfg.Port,
//			//viper.GetString("redis.host"),
//			//viper.GetString("redis.port"),
//		),
//		Password: cfg.Password, //viper.GetString("redis.password"), // 密码
//		DB:       cfg.Database, //viper.GetInt("redis.db"),          // 数据库
//		PoolSize: cfg.PollSize, //viper.GetInt("redis.pool_size"),   // 连接池大小
//	})
//
//	//接收Ping的结果
//	_, err = rdb.Ping().Result()
//	return err
//}

func Init(Addr, Password string, PoolSize, DB int) *operate.RDB {
	rdb := redis.NewClient(&redis.Options{
		Addr:     Addr,     // host:port
		Password: Password, // 密码
		PoolSize: PoolSize, // 连接池
		DB:       DB,       // 默认连接数据库（0-15）
	})
	//context.Background()提供上下文环境
	//Result() 获取命令执行的结果和错误信息
	_, err := rdb.Ping(context.Background()).Result() //测试连接
	if err != nil {
		panic(err)
	}
	//operate.New(rdb)就是将rdb包装成RDB类型
	return operate.New(rdb)
}
