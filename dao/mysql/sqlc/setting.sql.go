// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: setting.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createSetting = `-- name: CreateSetting :exec
insert into settings (account_id, relation_id, nick_name,  is_leader, is_self)
values (?,?,'',?,?)
`

type CreateSettingParams struct {
	AccountID  int64
	RelationID int64
	IsLeader   bool
	IsSelf     bool
}

func (q *Queries) CreateSetting(ctx context.Context, arg *CreateSettingParams) error {
	_, err := q.exec(ctx, q.createSettingStmt, createSetting,
		arg.AccountID,
		arg.RelationID,
		arg.IsLeader,
		arg.IsSelf,
	)
	return err
}

const deleteSetting = `-- name: DeleteSetting :exec
delete
from settings
where account_id = ?
  and relation_id = ?
`

type DeleteSettingParams struct {
	AccountID  int64
	RelationID int64
}

func (q *Queries) DeleteSetting(ctx context.Context, arg *DeleteSettingParams) error {
	_, err := q.exec(ctx, q.deleteSettingStmt, deleteSetting, arg.AccountID, arg.RelationID)
	return err
}

const deleteSettingsByAccountID = `-- name: DeleteSettingsByAccountID :exec
delete
from settings
where account_id = ?
`

func (q *Queries) DeleteSettingsByAccountID(ctx context.Context, accountID int64) error {
	_, err := q.exec(ctx, q.deleteSettingsByAccountIDStmt, deleteSettingsByAccountID, accountID)
	return err
}

const existsFriendSetting = `-- name: ExistsFriendSetting :one
select exists(select 1
              from settings s,
                   relations r
              where r.relation = 'friend'
                and (account1_id = ? and account2_id = ?) or
                     (account1_id = ? and account2_id = ?)
                and s.account_id = ?
)
`

type ExistsFriendSettingParams struct {
	Account1ID   sql.NullInt64
	Account2ID   sql.NullInt64
	Account1ID_2 sql.NullInt64
	Account2ID_2 sql.NullInt64
	AccountID    int64
}

func (q *Queries) ExistsFriendSetting(ctx context.Context, arg *ExistsFriendSettingParams) (bool, error) {
	row := q.queryRow(ctx, q.existsFriendSettingStmt, existsFriendSetting,
		arg.Account1ID,
		arg.Account2ID,
		arg.Account1ID_2,
		arg.Account2ID_2,
		arg.AccountID,
	)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const existsGroupLeaderByAccountIDWithLock = `-- name: ExistsGroupLeaderByAccountIDWithLock :one
select exists(select 1 from settings where account_id = ? and is_leader = true) for update
`

func (q *Queries) ExistsGroupLeaderByAccountIDWithLock(ctx context.Context, accountID int64) (bool, error) {
	row := q.queryRow(ctx, q.existsGroupLeaderByAccountIDWithLockStmt, existsGroupLeaderByAccountIDWithLock, accountID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const existsIsLeader = `-- name: ExistsIsLeader :one
select exists(select 1 from settings where relation_id = ? and account_id = ? and is_leader is true)
`

type ExistsIsLeaderParams struct {
	RelationID int64
	AccountID  int64
}

func (q *Queries) ExistsIsLeader(ctx context.Context, arg *ExistsIsLeaderParams) (bool, error) {
	row := q.queryRow(ctx, q.existsIsLeaderStmt, existsIsLeader, arg.RelationID, arg.AccountID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const existsSetting = `-- name: ExistsSetting :one
select exists(select 1 from settings where account_id = ? and relation_id = ?)
`

type ExistsSettingParams struct {
	AccountID  int64
	RelationID int64
}

func (q *Queries) ExistsSetting(ctx context.Context, arg *ExistsSettingParams) (bool, error) {
	row := q.queryRow(ctx, q.existsSettingStmt, existsSetting, arg.AccountID, arg.RelationID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const getFriendPinSettingsOrderByPinTime = `-- name: GetFriendPinSettingsOrderByPinTime :many
SELECT s.relation_id, s.nick_name, s.pin_time,
       a.id AS account_id,
       a.name AS account_name,
       a.avatar AS account_avatar
FROM (
         SELECT settings.relation_id,
                settings.nick_name,
                settings.pin_time
         FROM settings, relations
         WHERE settings.account_id = ?
           AND settings.is_pin = true
           AND settings.relation_id = relations.id
           AND relations.relation = 'friend'
     ) AS s,
     accounts a
WHERE a.id = (
    -- 在子查询中显式指定表别名解决歧义
    SELECT sub_settings.account_id
    FROM settings AS sub_settings  -- 使用别名区分
    WHERE sub_settings.relation_id = s.relation_id
      AND (sub_settings.account_id != ? OR sub_settings.is_self = true)
    LIMIT 1  -- 确保子查询返回单值
)
ORDER BY s.pin_time
`

type GetFriendPinSettingsOrderByPinTimeParams struct {
	AccountID   int64
	AccountID_2 int64
}

type GetFriendPinSettingsOrderByPinTimeRow struct {
	RelationID    int64
	NickName      string
	PinTime       time.Time
	AccountID     int64
	AccountName   string
	AccountAvatar string
}

func (q *Queries) GetFriendPinSettingsOrderByPinTime(ctx context.Context, arg *GetFriendPinSettingsOrderByPinTimeParams) ([]*GetFriendPinSettingsOrderByPinTimeRow, error) {
	rows, err := q.query(ctx, q.getFriendPinSettingsOrderByPinTimeStmt, getFriendPinSettingsOrderByPinTime, arg.AccountID, arg.AccountID_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetFriendPinSettingsOrderByPinTimeRow{}
	for rows.Next() {
		var i GetFriendPinSettingsOrderByPinTimeRow
		if err := rows.Scan(
			&i.RelationID,
			&i.NickName,
			&i.PinTime,
			&i.AccountID,
			&i.AccountName,
			&i.AccountAvatar,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFriendSettingsByName = `-- name: GetFriendSettingsByName :many
SELECT s.relation_id, s.nick_name, s.is_not_disturb, s.is_pin, s.pin_time, s.is_show, s.last_show, s.is_self,
       a.id AS account_id,
       a.name AS account_name,
       a.avatar AS account_avatar,
       COUNT(*) OVER () AS total
FROM (
         SELECT settings.relation_id,
                settings.nick_name,
                settings.is_not_disturb,
                settings.is_pin,
                settings.pin_time,
                settings.is_show,
                settings.last_show,
                settings.is_self
         FROM settings, relations
         WHERE settings.account_id = ?
           AND settings.relation_id = relations.id
           AND relations.relation = 'friend'
     ) AS s,
     accounts a
WHERE a.id in (
    SELECT sub_s.account_id
    FROM settings AS sub_s
    WHERE sub_s.relation_id = s.relation_id
      AND (sub_s.account_id != ? OR s.is_self = true)
)
  AND ((a.name LIKE CONCAT('%', ?, '%'))
    OR (s.nick_name LIKE CONCAT('%', ?, '%')))
ORDER BY a.name
LIMIT ? OFFSET ?
`

type GetFriendSettingsByNameParams struct {
	AccountID   int64
	AccountID_2 int64
	CONCAT      interface{}
	CONCAT_2    interface{}
	Limit       int32
	Offset      int32
}

type GetFriendSettingsByNameRow struct {
	RelationID    int64
	NickName      string
	IsNotDisturb  bool
	IsPin         bool
	PinTime       time.Time
	IsShow        bool
	LastShow      time.Time
	IsSelf        bool
	AccountID     int64
	AccountName   string
	AccountAvatar string
	Total         interface{}
}

func (q *Queries) GetFriendSettingsByName(ctx context.Context, arg *GetFriendSettingsByNameParams) ([]*GetFriendSettingsByNameRow, error) {
	rows, err := q.query(ctx, q.getFriendSettingsByNameStmt, getFriendSettingsByName,
		arg.AccountID,
		arg.AccountID_2,
		arg.CONCAT,
		arg.CONCAT_2,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetFriendSettingsByNameRow{}
	for rows.Next() {
		var i GetFriendSettingsByNameRow
		if err := rows.Scan(
			&i.RelationID,
			&i.NickName,
			&i.IsNotDisturb,
			&i.IsPin,
			&i.PinTime,
			&i.IsShow,
			&i.LastShow,
			&i.IsSelf,
			&i.AccountID,
			&i.AccountName,
			&i.AccountAvatar,
			&i.Total,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFriendSettingsOrderByName = `-- name: GetFriendSettingsOrderByName :many
select s.relation_id, s.nick_name, s.is_not_disturb, s.is_pin, s.pin_time, s.is_show, s.last_show, s.is_self,
       a.id as account_id,
       a.name as account_name,
       a.avatar as account_avatar
from (select settings.relation_id,
             settings.nick_name,
             settings.is_not_disturb,
             settings.is_pin,
             settings.pin_time,
             settings.is_show,
             settings.last_show,
             settings.is_self
      from settings,
           relations
      where settings.account_id = ?
        and settings.relation_id = relations.id
        and relation = 'friend') as s,
     accounts a
where a.id in (select account_id from settings where relation_id = s.relation_id) and a.id != ?
order by a.name
`

type GetFriendSettingsOrderByNameParams struct {
	AccountID int64
	ID        int64
}

type GetFriendSettingsOrderByNameRow struct {
	RelationID    int64
	NickName      string
	IsNotDisturb  bool
	IsPin         bool
	PinTime       time.Time
	IsShow        bool
	LastShow      time.Time
	IsSelf        bool
	AccountID     int64
	AccountName   string
	AccountAvatar string
}

func (q *Queries) GetFriendSettingsOrderByName(ctx context.Context, arg *GetFriendSettingsOrderByNameParams) ([]*GetFriendSettingsOrderByNameRow, error) {
	rows, err := q.query(ctx, q.getFriendSettingsOrderByNameStmt, getFriendSettingsOrderByName, arg.AccountID, arg.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetFriendSettingsOrderByNameRow{}
	for rows.Next() {
		var i GetFriendSettingsOrderByNameRow
		if err := rows.Scan(
			&i.RelationID,
			&i.NickName,
			&i.IsNotDisturb,
			&i.IsPin,
			&i.PinTime,
			&i.IsShow,
			&i.LastShow,
			&i.IsSelf,
			&i.AccountID,
			&i.AccountName,
			&i.AccountAvatar,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFriendShowSettingsOrderByShowTime = `-- name: GetFriendShowSettingsOrderByShowTime :many
SELECT s.relation_id, s.nick_name, s.is_not_disturb, s.is_pin, s.pin_time, s.is_show, s.last_show, s.is_self,
       a.id AS account_id,
       a.name AS account_name,
       a.avatar AS account_avatar
FROM (
         SELECT settings.relation_id,
                settings.nick_name,
                settings.is_not_disturb,
                settings.is_pin,
                settings.pin_time,
                settings.is_show,
                settings.last_show,
                settings.is_self
         FROM settings
                  JOIN relations ON settings.relation_id = relations.id
         WHERE settings.account_id = ?  -- 当前用户ID
           AND settings.is_show = true
           AND relations.relation = 'friend'
     ) AS s
         JOIN accounts a ON a.id = (
    SELECT sub_s.account_id
    FROM settings AS sub_s  -- 使用别名解决歧义
    WHERE sub_s.relation_id = s.relation_id
      AND (sub_s.account_id != ? OR sub_s.is_self = true)  -- 明确使用别名
    LIMIT 1  -- 确保返回单值
)
ORDER BY s.last_show DESC
`

type GetFriendShowSettingsOrderByShowTimeParams struct {
	AccountID   int64
	AccountID_2 int64
}

type GetFriendShowSettingsOrderByShowTimeRow struct {
	RelationID    int64
	NickName      string
	IsNotDisturb  bool
	IsPin         bool
	PinTime       time.Time
	IsShow        bool
	LastShow      time.Time
	IsSelf        bool
	AccountID     int64
	AccountName   string
	AccountAvatar string
}

func (q *Queries) GetFriendShowSettingsOrderByShowTime(ctx context.Context, arg *GetFriendShowSettingsOrderByShowTimeParams) ([]*GetFriendShowSettingsOrderByShowTimeRow, error) {
	rows, err := q.query(ctx, q.getFriendShowSettingsOrderByShowTimeStmt, getFriendShowSettingsOrderByShowTime, arg.AccountID, arg.AccountID_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetFriendShowSettingsOrderByShowTimeRow{}
	for rows.Next() {
		var i GetFriendShowSettingsOrderByShowTimeRow
		if err := rows.Scan(
			&i.RelationID,
			&i.NickName,
			&i.IsNotDisturb,
			&i.IsPin,
			&i.PinTime,
			&i.IsShow,
			&i.LastShow,
			&i.IsSelf,
			&i.AccountID,
			&i.AccountName,
			&i.AccountAvatar,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getGroupPinSettingsOrderByPinTime = `-- name: GetGroupPinSettingsOrderByPinTime :many
SELECT s.relation_id,
       s.nick_name,
       s.pin_time,
       r.id,
       r.name,
       r.description,
       r.avatar
FROM (
         SELECT settings.relation_id,
                settings.nick_name,
                settings.pin_time
         FROM settings
                  JOIN relations ON settings.relation_id = relations.id
         WHERE settings.account_id = ?
           AND settings.is_pin = true
           AND relations.relation = 'group'  -- 明确指定表别名
     ) AS s
         JOIN relations r ON r.id = s.relation_id
WHERE EXISTS (
    SELECT 1
    FROM settings sub_s  -- 使用别名区分子查询
    WHERE sub_s.relation_id = s.relation_id
      AND sub_s.account_id = ?
)
ORDER BY s.pin_time
`

type GetGroupPinSettingsOrderByPinTimeParams struct {
	AccountID   int64
	AccountID_2 int64
}

type GetGroupPinSettingsOrderByPinTimeRow struct {
	RelationID  int64
	NickName    string
	PinTime     time.Time
	ID          int64
	Name        sql.NullString
	Description sql.NullString
	Avatar      sql.NullString
}

func (q *Queries) GetGroupPinSettingsOrderByPinTime(ctx context.Context, arg *GetGroupPinSettingsOrderByPinTimeParams) ([]*GetGroupPinSettingsOrderByPinTimeRow, error) {
	rows, err := q.query(ctx, q.getGroupPinSettingsOrderByPinTimeStmt, getGroupPinSettingsOrderByPinTime, arg.AccountID, arg.AccountID_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetGroupPinSettingsOrderByPinTimeRow{}
	for rows.Next() {
		var i GetGroupPinSettingsOrderByPinTimeRow
		if err := rows.Scan(
			&i.RelationID,
			&i.NickName,
			&i.PinTime,
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Avatar,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getGroupShowSettingsOrderByShowTime = `-- name: GetGroupShowSettingsOrderByShowTime :many
SELECT s.relation_id, s.nick_name, s.is_not_disturb, s.is_pin, s.pin_time, s.is_show, s.last_show, s.is_self,
       r.id,
       r.name,
       r.description,
       r.avatar
FROM (
         SELECT settings.relation_id,
                settings.nick_name,
                settings.is_not_disturb,
                settings.is_pin,
                settings.pin_time,
                settings.is_show,
                settings.last_show,
                settings.is_self
         FROM settings
                  JOIN relations ON settings.relation_id = relations.id
         WHERE settings.account_id = ?
           AND settings.is_show = true
           AND relations.relation = 'group'
     ) AS s
         JOIN relations r ON r.id = s.relation_id
WHERE EXISTS (
    SELECT 1
    FROM settings sub_s  -- 使用别名解决歧义
    WHERE sub_s.relation_id = s.relation_id
      AND sub_s.account_id = ?
)
ORDER BY s.last_show DESC
`

type GetGroupShowSettingsOrderByShowTimeParams struct {
	AccountID   int64
	AccountID_2 int64
}

type GetGroupShowSettingsOrderByShowTimeRow struct {
	RelationID   int64
	NickName     string
	IsNotDisturb bool
	IsPin        bool
	PinTime      time.Time
	IsShow       bool
	LastShow     time.Time
	IsSelf       bool
	ID           int64
	Name         sql.NullString
	Description  sql.NullString
	Avatar       sql.NullString
}

func (q *Queries) GetGroupShowSettingsOrderByShowTime(ctx context.Context, arg *GetGroupShowSettingsOrderByShowTimeParams) ([]*GetGroupShowSettingsOrderByShowTimeRow, error) {
	rows, err := q.query(ctx, q.getGroupShowSettingsOrderByShowTimeStmt, getGroupShowSettingsOrderByShowTime, arg.AccountID, arg.AccountID_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetGroupShowSettingsOrderByShowTimeRow{}
	for rows.Next() {
		var i GetGroupShowSettingsOrderByShowTimeRow
		if err := rows.Scan(
			&i.RelationID,
			&i.NickName,
			&i.IsNotDisturb,
			&i.IsPin,
			&i.PinTime,
			&i.IsShow,
			&i.LastShow,
			&i.IsSelf,
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Avatar,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRelationIDsByAccountIDFromSettings = `-- name: GetRelationIDsByAccountIDFromSettings :many
select relation_id
from settings
where account_id = ?
`

func (q *Queries) GetRelationIDsByAccountIDFromSettings(ctx context.Context, accountID int64) ([]int64, error) {
	rows, err := q.query(ctx, q.getRelationIDsByAccountIDFromSettingsStmt, getRelationIDsByAccountIDFromSettings, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []int64{}
	for rows.Next() {
		var relation_id int64
		if err := rows.Scan(&relation_id); err != nil {
			return nil, err
		}
		items = append(items, relation_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSettingByID = `-- name: GetSettingByID :one
select account_id, relation_id, nick_name, is_not_disturb, is_pin, pin_time, is_show, last_show, is_leader, is_self
from settings
where account_id = ?
  and relation_id = ?
`

type GetSettingByIDParams struct {
	AccountID  int64
	RelationID int64
}

func (q *Queries) GetSettingByID(ctx context.Context, arg *GetSettingByIDParams) (*Setting, error) {
	row := q.queryRow(ctx, q.getSettingByIDStmt, getSettingByID, arg.AccountID, arg.RelationID)
	var i Setting
	err := row.Scan(
		&i.AccountID,
		&i.RelationID,
		&i.NickName,
		&i.IsNotDisturb,
		&i.IsPin,
		&i.PinTime,
		&i.IsShow,
		&i.LastShow,
		&i.IsLeader,
		&i.IsSelf,
	)
	return &i, err
}

const transferIsLeaderFalse = `-- name: TransferIsLeaderFalse :exec
update settings
set is_leader = false
where relation_id = ?
  and account_id = ?
`

type TransferIsLeaderFalseParams struct {
	RelationID int64
	AccountID  int64
}

func (q *Queries) TransferIsLeaderFalse(ctx context.Context, arg *TransferIsLeaderFalseParams) error {
	_, err := q.exec(ctx, q.transferIsLeaderFalseStmt, transferIsLeaderFalse, arg.RelationID, arg.AccountID)
	return err
}

const transferIsLeaderTrue = `-- name: TransferIsLeaderTrue :exec
update settings
set is_leader = true
where relation_id = ?
  and account_id = ?
`

type TransferIsLeaderTrueParams struct {
	RelationID int64
	AccountID  int64
}

func (q *Queries) TransferIsLeaderTrue(ctx context.Context, arg *TransferIsLeaderTrueParams) error {
	_, err := q.exec(ctx, q.transferIsLeaderTrueStmt, transferIsLeaderTrue, arg.RelationID, arg.AccountID)
	return err
}

const updateSettingDisturb = `-- name: UpdateSettingDisturb :exec
update settings
set is_not_disturb = ?
where account_id = ?
  and relation_id = ?
`

type UpdateSettingDisturbParams struct {
	IsNotDisturb bool
	AccountID    int64
	RelationID   int64
}

func (q *Queries) UpdateSettingDisturb(ctx context.Context, arg *UpdateSettingDisturbParams) error {
	_, err := q.exec(ctx, q.updateSettingDisturbStmt, updateSettingDisturb, arg.IsNotDisturb, arg.AccountID, arg.RelationID)
	return err
}

const updateSettingLeader = `-- name: UpdateSettingLeader :exec
update settings
set is_leader = ?
where account_id = ?
  and relation_id = ?
`

type UpdateSettingLeaderParams struct {
	IsLeader   bool
	AccountID  int64
	RelationID int64
}

func (q *Queries) UpdateSettingLeader(ctx context.Context, arg *UpdateSettingLeaderParams) error {
	_, err := q.exec(ctx, q.updateSettingLeaderStmt, updateSettingLeader, arg.IsLeader, arg.AccountID, arg.RelationID)
	return err
}

const updateSettingNickName = `-- name: UpdateSettingNickName :exec
update settings
set nick_name = ?
where account_id = ?
  and relation_id = ?
`

type UpdateSettingNickNameParams struct {
	NickName   string
	AccountID  int64
	RelationID int64
}

func (q *Queries) UpdateSettingNickName(ctx context.Context, arg *UpdateSettingNickNameParams) error {
	_, err := q.exec(ctx, q.updateSettingNickNameStmt, updateSettingNickName, arg.NickName, arg.AccountID, arg.RelationID)
	return err
}

const updateSettingPin = `-- name: UpdateSettingPin :exec
update settings
set is_pin = ?
where account_id = ?
  and relation_id = ?
`

type UpdateSettingPinParams struct {
	IsPin      bool
	AccountID  int64
	RelationID int64
}

func (q *Queries) UpdateSettingPin(ctx context.Context, arg *UpdateSettingPinParams) error {
	_, err := q.exec(ctx, q.updateSettingPinStmt, updateSettingPin, arg.IsPin, arg.AccountID, arg.RelationID)
	return err
}

const updateSettingShow = `-- name: UpdateSettingShow :exec
update settings
set is_show = ?
where account_id = ?
  and relation_id = ?
`

type UpdateSettingShowParams struct {
	IsShow     bool
	AccountID  int64
	RelationID int64
}

func (q *Queries) UpdateSettingShow(ctx context.Context, arg *UpdateSettingShowParams) error {
	_, err := q.exec(ctx, q.updateSettingShowStmt, updateSettingShow, arg.IsShow, arg.AccountID, arg.RelationID)
	return err
}
