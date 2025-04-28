package dao

import (
	my "chat/dao/mysql"
	"chat/dao/redis/operate"
)

type datebase struct {
	DB    my.DB //查询接口和事务接口
	Redis *operate.RDB
}

var Database = new(datebase)
