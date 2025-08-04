package tx

import (
	db "chat/dao/mysql/sqlc"
	"chat/global"
	"context"
	"database/sql"
	"errors"
)

func (store *MySQLDB) UploadGroupAvatarWithTx(ctx context.Context, arg db.CreateFileParams) error {
	return store.execTx(ctx, func(queries *db.Queries) error {
		var err error
		_, err = queries.GetGroupAvatar(ctx, arg.RelationID)
		if err != nil && arg.Url != global.PublicSetting.Rules.DefaultAvatarURL {
			// 如果没有设置过群头像 或者 是默认头像 并且 新的头像不是默认头像
			if errors.Is(sql.ErrNoRows, err) {
				err = queries.CreateFile(ctx, &db.CreateFileParams{
					FileName:   arg.FileName,
					FileType:   "image",
					FileSize:   arg.FileSize,
					FileKey:    arg.FileKey,
					Url:        arg.Url,
					RelationID: arg.RelationID,
					AccountID:  sql.NullInt64{},
				})
			} else {
				return err
			}
		} else {
			// 在 file 表中覆盖之前的群头像
			err = queries.UpdateGroupAvatar(ctx, &db.UpdateGroupAvatarParams{
				Url:        arg.Url,
				RelationID: arg.RelationID,
			})
		}
		data, err := queries.GetGroupRelationByID(ctx, arg.RelationID.Int64)
		if err != nil {
			return err
		}
		// 更新 relation 表中的头像数据
		return queries.UpdateGroupRelation(ctx, &db.UpdateGroupRelationParams{
			Name: data.Name,

			Description: data.Description,
			Avatar: sql.NullString{
				String: arg.Url,
				Valid:  true,
			},
			ID: arg.RelationID.Int64,
		})
	})
}
