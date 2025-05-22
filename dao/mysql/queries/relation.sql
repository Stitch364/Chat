-- name: CreateFriendRelation :exec
insert into relations (relation,account1_id,account2_id)
values ('friend',?,?);

-- name: GetFriendRelationIdByID1AndID1 :one
select id from relations where account1_id = ? and account2_id = ? and relation = 'friend';

-- name: GetFriendRelationIDsByAccountID :many
select  id
from relations
where relation = 'friend'
  and (account1_id = ? or account2_id = ?);


-- name: DeleteFriendRelationsByAccountID :exec
delete
from relations
where relation = 'friend'
  and (account1_id = ? or account2_id = ?);