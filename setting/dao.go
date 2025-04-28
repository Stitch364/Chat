package setting

import (
	"chat/dao"
	"chat/dao/mysql"
	"chat/dao/redis"
	"chat/global"
)

type database struct {
}

// Init 数据库初始化
func (database) Init() {
	//mysql 初始化
	dao.Database.DB = mysql.Init(global.PrivateSetting.Mysql)
	//redis 初始化
	dao.Database.Redis = redis.Init(
		global.PrivateSetting.Redis.Address,
		global.PrivateSetting.Redis.Password,
		global.PrivateSetting.Redis.PoolSize,
		global.PrivateSetting.Redis.DB,
	)
}
