-- name: CreateMessage :exec
-- {notify_type:string,msg_type:string}
insert into messages
(notify_type, msg_type, msg_content, msg_extend, file_id, account_id, rly_msg_id, relation_id)
values
(?,?,?,JSON_ARRAY(),?,?,?,?);

-- name: CreateMessageReturn :one
SELECT
    id, msg_content, COALESCE(msg_extend,'{}'), file_id, create_at
FROM messages
WHERE id = LAST_INSERT_ID();

-- name: GetMessageInfoTx :one
SELECT m.id, msg_content, msg_extend,file_id, create_at,m.account_id,a.name,a.avatar,s.nick_name
FROM messages m
join accounts a on a.id = m.account_id
join settings s on s.account_id = a.id and s.relation_id = m.relation_id
WHERE m.id = LAST_INSERT_ID();



-- name: GetMessageByID :one
select id, notify_type, msg_type, msg_content, coalesce(msg_extend,'[]'), file_id, account_id,
       rly_msg_id, relation_id, create_at, is_revoke, is_top, is_pin, pin_time, read_ids, is_delete
from messages
where id = ?
limit 1;

-- name: GetMessageAndNameByID :one
select m.id, notify_type, msg_type, msg_content, coalesce(msg_extend,'[]'), file_id, m.account_id,a.name,s.nick_name,a.avatar,
       rly_msg_id, m.relation_id, create_at, is_revoke, is_top, m.is_pin, m.pin_time, read_ids, is_delete
from messages m
join accounts a on a.id = m.account_id
join settings s on s.account_id  = m.account_id and s.relation_id = m.relation_id
where m.id = ?
limit 1;

-- name: GetAccountInfoByID :one
select accounts.name,accounts.avatar,settings.nick_name
from accounts
join settings on accounts.id = settings.account_id  and relation_id = ?
where account_id = ?;

-- name: GetMsgsByRelationIDAndTime :many
select m1.id,
       m1.notify_type,
       m1.msg_type,
       m1.msg_content,
       coalesce(m1.msg_extend,'[]'),
       m1.file_id,
       m1.account_id,
       a.name,
       a.avatar,
       m1.rly_msg_id,
       m1.relation_id,
       m1.create_at,
       m1.is_revoke,
       m1.is_top,
       m1.is_pin,
       m1.pin_time,
       m1.read_ids,
       m1.is_delete,
       count(*) over () as total,
       (select count(id) from messages where rly_msg_id = m1.id and messages.relation_id = ?) as reply_count,
       s.nick_name
from messages m1
join accounts a on a.id = m1.account_id
JOIN   settings  s
       ON s.relation_id = m1.relation_id
           AND s.account_id = CASE
              WHEN (SELECT r.relation
                    FROM relations r
                    WHERE r.id = ?) = 'friend'
                  THEN  -- 好友关系：取对方账号的 setting
                  (SELECT CASE
                              WHEN m1.account_id = r1.account1_id THEN r1.account2_id
                              ELSE r1.account1_id
                              END
                   FROM relations r1
                   WHERE r1.id = m1.relation_id)
              ELSE  -- 非好友关系：取自己账号的 setting
                  m1.account_id
                  END
where m1.relation_id = ?
  and m1.create_at < ?
order by m1.create_at desc
limit ? offset ?;

-- name: GetNickNameByAccountIDAndRelation :one
select nick_name
from settings
where account_id = ?
  and relation_id = ?;

-- name: OfferMsgsByAccountIDAndTime :many
select m1.id,
       m1.notify_type,
       m1.msg_type,
       m1.msg_content,
       coalesce(m1.msg_extend,'[]'),
       m1.file_id,
       m1.account_id,
       m1.rly_msg_id,
       m1.relation_id,
       m1.create_at,
       m1.is_revoke,
       m1.is_top,
       m1.is_pin,
       m1.pin_time,
       m1.read_ids,
       m1.is_delete,
       count(*) over () as total,
       (select count(id) from messages where rly_msg_id = m1.id and messages.relation_id = m1.relation_id) as reply_count
from messages m1
         join settings s on m1.relation_id = s.relation_id and s.account_id = ?
where m1.create_at > ?
limit ? offset ?;

-- name: GetPinMsgsByRelationID :many
select m1.id,
       m1.notify_type,
       m1.msg_type,
       m1.msg_content,
       coalesce(m1.msg_extend,'[]'),
       m1.file_id,
       m1.account_id,
       m1.relation_id,
       m1.create_at,
       m1.is_revoke,
       m1.is_top,
       m1.is_pin,
       m1.pin_time,
       m1.read_ids,
       m1.is_delete,
       (select count(id) from messages where rly_msg_id = m1.id and messages.relation_id = ?) as reply_count,
       count(*) over () as total
from messages m1
where m1.relation_id = ? and m1.is_pin = true
order by m1.pin_time desc
limit ? offset ?;

-- name: GetRlyMsgsInfoByMsgID :many
select m1.id,
       m1.notify_type,
       m1.msg_type,
       m1.msg_content,
       coalesce(m1.msg_extend,'[]'),
       m1.file_id,
       m1.account_id,
       a.name,
       a.avatar,
       s.nick_name,
       m1.relation_id,
       m1.create_at,
       m1.is_revoke,
       m1.is_top,
       m1.is_pin,
       m1.pin_time,
       m1.read_ids,
       m1.is_delete,
       (select count(id) from messages where rly_msg_id = m1.id and messages.relation_id = ?) as reply_count,
       count(*) over () as total
from messages m1
join settings s on m1.relation_id = s.relation_id and s.account_id = m1.account_id
join accounts a on a.id = m1.account_id
where m1.relation_id = ? and m1.rly_msg_id = ?
order by m1.create_at
limit ? offset ?;

-- name: GetMsgsByContent :many
select m1.id,
       m1.notify_type,
       m1.msg_type,
       m1.msg_content,
       coalesce(m1.msg_extend,'[]'),
       m1.file_id,
       m1.account_id,
       a.name,
       a.avatar,
       s.nick_name,
       m1.relation_id,
       m1.create_at,
       m1.is_delete,
       count(*) over () as total
from messages m1
         join settings s on m1.relation_id = s.relation_id and s.account_id = ?
         join accounts a on a.id = m1.account_id
where (not is_revoke)
    and m1.msg_content like concat('%', ?, '%')
order by m1.create_at desc
    limit ? offset ?;

-- name: GetMsgsByContentAndRelation :many
select m1.id,
       m1.notify_type,
       m1.msg_type,
       m1.msg_content,
       coalesce(m1.msg_extend,'[]'),
       m1.file_id,
       m1.account_id,
       a.name,
       a.avatar,
       s.nick_name,
       m1.relation_id,
       m1.create_at,
       m1.is_delete,
       count(*) over () as total
from messages m1
         join settings s on m1.relation_id = ? and m1.relation_id = s.relation_id and s.account_id = ?
         join accounts a on a.id = m1.account_id
where (not is_revoke)
  and m1.msg_content like concat('%', ?, '%')
order by m1.create_at desc
limit ? offset ?;

-- name: GetTopMsgByRelationID :one
select m1.id,
       m1.notify_type,
       m1.msg_type,
       m1.msg_content,
       coalesce(m1.msg_extend,'[]'),
       m1.file_id,
       m1.account_id,
       m1.relation_id,
       m1.create_at,
       m1.is_revoke,
       m1.is_top,
       m1.is_pin,
       m1.pin_time,
       m1.read_ids,
       (select count(id) from messages where rly_msg_id = m1.id and messages.relation_id = ?) as reply_count,
       count(*) over () as total
from messages m1
where m1.relation_id = ? and m1.is_top = true
limit 1;

-- name: UpdateMsgPin :exec
update messages
set is_pin = ?
where id = ?;

-- name: UpdateMsgTop :exec
update messages
set is_top = ?
where id = ?;

-- name: UpdateMsgRevoke :exec
update messages
set is_revoke = ?
where id = ?;

-- name: UpdateMsgDelete :exec
update messages
set is_delete = ?
where id = ?;

-- name: GetAccountIDsByMsgID :one
select account1_id, account2_id
from relations r
where r.id = (
    select relation_id
    from messages m
    where m.id = ?
)
limit 1;

-- name: GetMsgDeleteById :one
select  is_delete
from messages
where id = ?;
