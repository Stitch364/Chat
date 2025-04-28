package operate

import (
	"context"
	"github.com/XYYSWK/Lutils/pkg/utils"
)

const keyGroup = "KeyGroup"

// AddRelationAccount 向群聊名单中添加人员(两个人的好友关系，相当于一个特殊的群聊)
func (r *RDB) AddRelationAccount(ctx context.Context, relationID int64, accountIDs ...int64) error {
	if len(accountIDs) == 0 {
		return nil
	}
	data := make([]interface{}, len(accountIDs))
	for i, v := range accountIDs {
		data[i] = utils.IDToString(v)
	}
	return r.rdb.SAdd(ctx, utils.LinkStr(keyGroup, utils.IDToString(relationID)), data...).Err()
}
