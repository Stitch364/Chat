package tx

import (
	db "chat/dao/mysql/sqlc"
	"chat/dao/redis/operate"
	"chat/global"
	"context"
	"database/sql"
)

// DeleteRelationWithTx 从数据库中删除关系并删除 redis 中的关系
func (store *MySQLDB) DeleteRelationWithTx(ctx context.Context, rdb *operate.RDB, relationID int64) error {
	return store.execTx(ctx, func(queries *db.Queries) error {
		err := queries.DeleteRelation(ctx, relationID)
		if err != nil {
			return err
		}
		return rdb.DeleteRelations(ctx, relationID)
	})
}

func (store *MySQLDB) CreateGroupRelationWithTx(ctx context.Context, accountID int64, name string, description string) (err error, relationId int64) {
	var relationIdTemp int64
	return store.execTx(ctx, func(queries *db.Queries) error {
		err := queries.CreateGroupRelation(ctx, &db.CreateGroupRelationParams{
			Name: sql.NullString{
				String: name,
				Valid:  true,
			},
			Description: sql.NullString{
				String: description,
				Valid:  true,
			},
			Avatar: sql.NullString{
				String: global.PublicSetting.Rules.DefaultAvatarURL,
				Valid:  true,
			},
		})
		if err != nil {
			return err
		}
		relationIdTemp, err = queries.GetGroupRelationsId(ctx)
		return err
	}), relationIdTemp
}
