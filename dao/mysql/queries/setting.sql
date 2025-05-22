-- name: CreateSetting :exec
insert into settings (account_id, relation_id, nick_name,  is_leader, is_self)
values (?,?,'',?,?);

-- name: GetRelationIDsByAccountIDFromSettings :many
select relation_id
from settings
where account_id = ?;

-- name: DeleteSettingsByAccountID :exec
delete
from settings
where account_id = ?;

-- name: ExistsGroupLeaderByAccountIDWithLock :one
select exists(select 1 from settings where account_id = ? and is_leader = true) for update;

-- name: ExistsSetting :one
select exists(select 1 from settings where account_id = ? and relation_id = ?)