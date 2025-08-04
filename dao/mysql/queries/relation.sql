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

-- name: GetAllRelationIDs :many
select id
from relations;

-- name: GetAccountIDsByRelationID :many
select distinct account_id
from settings
where relation_id = ?;



-- name: CreateGroupRelation :exec
insert into relations (relation,name,description,avatar)
    values ('group', ?,?,?);

-- name: GetGroupRelationsId :one
select LAST_INSERT_ID()
from relations;



-- name: DeleteRelation :exec
delete
from relations
where id = ?;

-- name: UpdateGroupRelation :exec
update relations
set name=?,
       description=?,
       avatar=?
where relation = 'group'
  and id = ?;

-- name: GetGroupRelationByID :one
select id,
       account1_id,
       account2_id,
       name,
        description,
        avatar,
        created_at
from relations
where relation = 'group'
  and id = ?;

-- name: ExistsFriendRelation :one
select exists(select 1
              from relations
              where relation = 'friend'
                and account1_id = ?
               and account2_id = ?);

-- name: GetFriendRelationByID :one
select account1_id,
       account2_id,
       created_at
from relations
where relation = 'friend'
  and id = ?;

-- name: GetAllGroupRelation :many
select id
from relations
where relation = 'group';

-- name: GetAllRelationOnRelation :many
select *
from relations;


-- name: GetRelationIDByAccountID :one
select id
from relations
where account1_id = ?
  and account2_id = ?;

-- name: GetGroupList :many
select s.relation_id, s.nick_name, s.is_not_disturb, s.is_pin, s.pin_time, s.is_show, s.last_show, s.is_self,s.is_leader,
       r.id as relation_id,
       r.name as group_name,
       r.description,
       r.avatar as group_avatar,
       count(*) over () as total
from (select relation_id,
    nick_name,
    is_not_disturb,
    is_pin,
    pin_time,
    is_show,
    last_show,
    is_self,
    is_leader
    from settings,
    relations
    where settings.account_id = ?
    and settings.relation_id = relations.id
    and relations.relation = 'group') as s,
    relations r
where r.id = (select s.relation_id from settings where relation_id = s.relation_id and (settings.account_id = ?))
order by s.last_show;

-- name: GetGroupSettingsByName :many
select s.relation_id, s.nick_name, s.is_not_disturb, s.is_pin, s.pin_time, s.is_show, s.last_show, s.is_self,
       r.id as relation_id,
       name as group_name,
       avatar as group_avatar,
       description,
       count(*) over () as total
from (select relation_id,
    nick_name,
    is_not_disturb,
    is_pin,
    pin_time,
    is_show,
    last_show,
    is_self
    from settings,
    relations
    where settings.account_id = ?
    and settings.relation_id = relations.id
    and relations.relation = 'group') as s,
    relations r
where r.id = (select s.relation_id from settings where relation_id = s.relation_id and (settings.account_id = ?))
  and ((name like CONCAT('%' ,?, '%')))
order by name
limit ? offset ?;

-- name: GetGroupMembersByID :many
select a.id, a.name, a.avatar, s.nick_name, s.is_leader
from accounts a
         left join settings s on a.id = s.account_id
where s.relation_id = ?
limit ? offset ?;

-- name: GetGroupAvatarByID :one
select avatar
from relations
where relations.id = ?
limit 1;
