// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: relation.sql

package db

import (
	"context"
	"database/sql"
)

const createFriendRelation = `-- name: CreateFriendRelation :exec
insert into relations (relation,account1_id,account2_id)
values ('friend',?,?)
`

type CreateFriendRelationParams struct {
	Account1ID sql.NullInt64
	Account2ID sql.NullInt64
}

func (q *Queries) CreateFriendRelation(ctx context.Context, arg *CreateFriendRelationParams) error {
	_, err := q.exec(ctx, q.createFriendRelationStmt, createFriendRelation, arg.Account1ID, arg.Account2ID)
	return err
}

const deleteFriendRelationsByAccountID = `-- name: DeleteFriendRelationsByAccountID :exec
delete
from relations
where relation = 'friend'
  and (account1_id = ? or account2_id = ?)
`

type DeleteFriendRelationsByAccountIDParams struct {
	Account1ID sql.NullInt64
	Account2ID sql.NullInt64
}

func (q *Queries) DeleteFriendRelationsByAccountID(ctx context.Context, arg *DeleteFriendRelationsByAccountIDParams) error {
	_, err := q.exec(ctx, q.deleteFriendRelationsByAccountIDStmt, deleteFriendRelationsByAccountID, arg.Account1ID, arg.Account2ID)
	return err
}

const getAccountIDsByRelationID = `-- name: GetAccountIDsByRelationID :many
select distinct account_id
from settings
where relation_id = ?
`

func (q *Queries) GetAccountIDsByRelationID(ctx context.Context, relationID int64) ([]int64, error) {
	rows, err := q.query(ctx, q.getAccountIDsByRelationIDStmt, getAccountIDsByRelationID, relationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []int64{}
	for rows.Next() {
		var account_id int64
		if err := rows.Scan(&account_id); err != nil {
			return nil, err
		}
		items = append(items, account_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllRelationIDs = `-- name: GetAllRelationIDs :many
select id
from relations
`

func (q *Queries) GetAllRelationIDs(ctx context.Context) ([]int64, error) {
	rows, err := q.query(ctx, q.getAllRelationIDsStmt, getAllRelationIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []int64{}
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFriendRelationIDsByAccountID = `-- name: GetFriendRelationIDsByAccountID :many
select  id
from relations
where relation = 'friend'
  and (account1_id = ? or account2_id = ?)
`

type GetFriendRelationIDsByAccountIDParams struct {
	Account1ID sql.NullInt64
	Account2ID sql.NullInt64
}

func (q *Queries) GetFriendRelationIDsByAccountID(ctx context.Context, arg *GetFriendRelationIDsByAccountIDParams) ([]int64, error) {
	rows, err := q.query(ctx, q.getFriendRelationIDsByAccountIDStmt, getFriendRelationIDsByAccountID, arg.Account1ID, arg.Account2ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []int64{}
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFriendRelationIdByID1AndID1 = `-- name: GetFriendRelationIdByID1AndID1 :one
select id from relations where account1_id = ? and account2_id = ? and relation = 'friend'
`

type GetFriendRelationIdByID1AndID1Params struct {
	Account1ID sql.NullInt64
	Account2ID sql.NullInt64
}

func (q *Queries) GetFriendRelationIdByID1AndID1(ctx context.Context, arg *GetFriendRelationIdByID1AndID1Params) (int64, error) {
	row := q.queryRow(ctx, q.getFriendRelationIdByID1AndID1Stmt, getFriendRelationIdByID1AndID1, arg.Account1ID, arg.Account2ID)
	var id int64
	err := row.Scan(&id)
	return id, err
}
