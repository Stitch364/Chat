-- name: CreateMessage :exec
-- {notify_type:string,msg_type:string}
insert into messages
(notify_type, msg_type, msg_content, msg_extend, file_id, account_id, rly_msg_id, relation_id)
values
(?,?,?,?,?,?,?,?);

-- name: GetMessageInfoTx :one
SELECT id, msg_content, msg_extend,file_id, create_at
FROM messages
WHERE id = LAST_INSERT_ID();



-- name: GetMessageByID :one
select id, notify_type, msg_type, msg_content, msg_extend, file_id, account_id,
       rly_msg_id, relation_id, create_at, is_revoke, is_top, is_pin, pin_time, read_ids
from messages
where id = ?
limit 1;

-- name: GetMsgsByRelationIDAndTime :many
select m1.id,
       m1.notify_type,
       m1.msg_type,
       m1.msg_content,
       m1.msg_extend,
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
       count(*) over () as total,
       (select count(id) from messages where rly_msg_id = m1.id and messages.relation_id = ?) as reply_count
from messages m1
where m1.relation_id = ?
  and m1.create_at < ?
order by m1.create_at desc
limit ? offset ?;

-- name: OfferMsgsByAccountIDAndTime :many
select m1.id,
       m1.notify_type,
       m1.msg_type,
       m1.msg_content,
       m1.msg_extend,
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
       m1.msg_extend,
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
where m1.relation_id = ? and m1.is_pin = true
order by m1.pin_time desc
limit ? offset ?;

-- name: GetRlyMsgsInfoByMsgID :many
select m1.id,
       m1.notify_type,
       m1.msg_type,
       m1.msg_content,
       m1.msg_extend,
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
where m1.relation_id = ? and m1.rly_msg_id = ?
order by m1.create_at
limit ? offset ?;

-- name: GetMsgsByContent :many
select m1.id,
       m1.notify_type,
       m1.msg_type,
       m1.msg_content,
       m1.msg_extend,
       m1.file_id,
       m1.account_id,
       m1.relation_id,
       m1.create_at,
       count(*) over () as total
from messages m1
         join settings s on m1.relation_id = s.relation_id and s.account_id = ?
where (not is_revoke)
    and (m1.msg_content like concat('%', ?, '%') or m1.msg_extend like concat('%', ?, '%'))
order by m1.create_at desc
    limit ? offset ?;

-- name: GetMsgsByContentAndRelation :many
select m1.id,
       m1.notify_type,
       m1.msg_type,
       m1.msg_content,
       m1.msg_extend,
       m1.file_id,
       m1.account_id,
       m1.relation_id,
       m1.create_at,
       count(*) over () as total
from messages m1
         join settings s on m1.relation_id = ? and m1.relation_id = s.relation_id and s.account_id = ?
where (not is_revoke)
  and (m1.msg_content like concat('%', ?, '%') or m1.msg_extend like concat('%', ?, '%'))
order by m1.create_at desc
limit ? offset ?;