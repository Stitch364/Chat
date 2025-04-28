-- name: CreateApplication :exec
-- {account1_id:int64,account2_id:int64,apply_msg:string}
insert into applications (account1_id, account2_id, apply_msg, refuse_msg)
values ( ?, ?, ?, '');

-- name: ExistsApplicationByIDWithLock :one
-- {account1_id:int64,account2_id:int64}
select exists(
    select 1
    from applications
    where (account1_id = ? and account2_id = ?)
       or (account1_id = ? and account2_id = ?)
        for update );


-- name: DeleteApplication :exec
-- {account1_id:int64,account2_id:int64}
delete
from applications
where account1_id = ?
  and account2_id = ?;

-- name: GetApplicationByID :one
-- {account1_id:int64,account2_id:int64}
select *
from applications
where account1_id = ? and account2_id = ?
limit 1;

-- name: UpdateApplication :exec
-- {status:string,refuse_msg:string,account1_id:int64,account2_id:int64}
update applications
set status = ?,
    refuse_msg = ?
where account1_id = ?
  and account2_id = ?;

-- name: GetApplications :many
-- {account1_id:int64,account2_id:int64,limit:int32,offset:int32,total:int64}
select app.*,
       a1.name as account1_name,
       a1.avatar as account1_avatar,
       a2.name as account2_name,
       a2.avatar as account2_avatar
from accounts a1,
     accounts a2,
     (select *, count(*) over () as total
      from applications
      where account1_id = ?
         or account2_id = ?
      order by create_at desc
      limit ? offset ?) as app
where a1.id = app.account1_id
  and a2.id = app.account2_id;