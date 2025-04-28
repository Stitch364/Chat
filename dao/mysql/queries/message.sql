-- name: CreateMessage :exec
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

-- name: UpdateMsgReads :exec
update messages
set read_ids = json_array_append(read_ids, '$', @accountID)
where relation_id = ?;

