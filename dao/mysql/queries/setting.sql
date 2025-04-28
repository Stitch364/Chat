-- name: CreateSetting :exec
insert into settings (account_id, relation_id, nick_name,  is_leader, is_self)
values (?,?,'',?,?);