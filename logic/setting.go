package logic

import (
	"chat/dao"
	db "chat/dao/mysql/sqlc"
	"chat/global"
	"context"
	"github.com/XYYSWK/Lutils/pkg/app/errcode"
)

type setting struct{}

// ExistsSetting 是否存在 account 和 relation 关系的联系
// 参数：accountID，relationDI
// 成功：是否存在，nil
// 失败：打印错误日志 errcode.ErrServer
func ExistsSetting(ctx context.Context, accountID, relationID int64) (bool, errcode.Err) {
	ok, err := dao.Database.DB.ExistsSetting(ctx, &db.ExistsSettingParams{
		AccountID:  accountID,
		RelationID: relationID,
	})
	if err != nil {
		global.Logger.Error(err.Error())
		return false, errcode.ErrServer
	}
	return ok, nil
}
