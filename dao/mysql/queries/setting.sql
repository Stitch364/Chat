-- name: CreateSetting :exec
insert into settings (account_id, relation_id, nick_name,  is_leader, is_self)
values (?,?,'',?,?);

-- name: DeleteSetting :exec
delete
from settings
where account_id = ?
  and relation_id = ?;

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
select exists(select 1 from settings where account_id = ? and relation_id = ?);

-- name: GetFriendSettingsOrderByName :many
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
order by a.name;

-- name: UpdateSettingNickName :exec
update settings
set nick_name = ?
where account_id = ?
  and relation_id = ?;

-- name: UpdateSettingDisturb :exec
update settings
set is_not_disturb = ?
where account_id = ?
  and relation_id = ?;

-- name: UpdateSettingPin :exec
update settings
set is_pin = ?
where account_id = ?
  and relation_id = ?;

-- name: UpdateSettingLeader :exec
update settings
set is_leader = ?
where account_id = ?
  and relation_id = ?;

-- name: UpdateSettingShow :exec
update settings
set is_show = ?
where account_id = ?
  and relation_id = ?;

-- name: GetSettingByID :one
select *
from settings
where account_id = ?
  and relation_id = ?;

-- name: GetFriendPinSettingsOrderByPinTime :many
SELECT s.*,
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
ORDER BY s.pin_time;



-- name: GetGroupPinSettingsOrderByPinTime :many
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
ORDER BY s.pin_time;

-- name: GetFriendShowSettingsOrderByShowTime :many
SELECT s.*,
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
ORDER BY s.last_show DESC;

-- name: GetGroupShowSettingsOrderByShowTime :many
SELECT s.*,
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
ORDER BY s.last_show DESC;


-- name: GetFriendSettingsByName :many
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
LIMIT ? OFFSET ?;

-- name: TransferIsLeaderFalse :exec
update settings
set is_leader = false
where relation_id = ?
  and account_id = ?;

-- name: TransferIsLeaderTrue :exec
update settings
set is_leader = true
where relation_id = ?
  and account_id = ?;

-- name: ExistsIsLeader :one
select exists(select 1 from settings where relation_id = ? and account_id = ? and is_leader is true);

-- name: ExistsFriendSetting :one
select exists(select 1
              from settings s,
                   relations r
              where r.relation = 'friend'
                and (account1_id = ? and account2_id = ?) or
                     (account1_id = ? and account2_id = ?)
                and s.account_id = ?
)