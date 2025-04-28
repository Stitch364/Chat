-- name: GetUserByEmail :one
-- {email:string}
select id, email, password, created_at
from users
where email = ?
limit 1;

-- name: GetUserByID :one
-- {id:int64}
select id, email, password, created_at
from users
where id = ?
limit 1;

-- name: ExistEmail :one
-- {email:string}
select exists(select 1 from users where email = ?);

-- name: CreateUser :exec
-- {email:string,password:string}
insert into users(email,password)
values (?,?);

-- name: DeleteUser :exec
-- {id:int64}
delete
from users
where id = ?;

-- name: ExistsUserByID :one
-- {id:int64}
select exists(select 1 from users where id = ?);

-- name: GetAccountIDsByUserID :many
-- {user_id:int64}
select id
from accounts
where user_id = ?;

-- name: GetAllEmail :many
select email
from users;

-- name: UpdateUser :exec
update users
set email = ?,
    password = ?
where id = ?;