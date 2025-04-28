-- name: CreateAccount :exec
-- {id:int64, user_id:int64, name:string, avatar:string, gender:string, signature:string}
insert into accounts (id, user_id, name, avatar, gender, signature)
values (?, ?, ?, ?, ?, ?);

-- name: DeleteAccount :exec
-- {id:int64}
delete
from accounts
where id = ?;

-- name: DeleteAccountsByUserID :exec
-- {user_id:int64}
delete
from accounts
where user_id = ?;

-- name: UpdateAccount :exec
-- {id:int64, name:string, gender:string, signature:string}
update accounts
set name        = ?,
    gender      = ?,
    signature   = ?
where id = ?;

-- name: UpdateAccountAvatar :exec
-- {id:int64, avatar:string}
update accounts
set avatar = ?
where id = ?;


-- name: GetAccountsByUserID :many
-- {user_id:int64}
select id, name, avatar, gender
from accounts
where user_id = ?;

-- name: CountAccountsByUserID :one
-- {user_id:int64}
select count(id)
from accounts
where user_id = ?;

-- name: ExistsAccountByID :one
-- {id:int64}
select exists(select 1 from accounts where id = ?);

-- name: ExistsAccountByNameAndUserID :one
-- {user_id:int64, name:string}
select exists(
    select 1
    from accounts
    where user_id = ?
      and name = ?
);

-- name: GetAccountByID :one
-- {accounts.id:int64, r.account2_id:int64, r.account1_id:int64}
select a.id, a.user_id, a.name, a.avatar, a.gender, a.signature, a.created_at, r.id as relation_id
from (select id, user_id, name, avatar, gender, signature, created_at from accounts where accounts.id = ?) a
         left join relations r on
    r.relation = 'friend' and
    r.account1_id = a.id and r.account2_id = ? or
    r.account1_id = ? and r.account2_id = a.id
limit 1;

-- name: GetAccountsByName :many
-- {name:string, user_id:int64, page:int64, page_size:int64}
select a.*, r.id as relation_id, count(*) over () as total
from (select id, name, avatar, gender from accounts where accounts.name like CONCAT('%', ?, '%')) as a
         left join relations r on (r.relation = 'friend' and
                                   ((r.account1_id = a.id and r.account2_id = ?) or
                                    (r.account1_id = ? and r.account2_id = a.id)))
limit ? offset ?;