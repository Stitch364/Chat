// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: application.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createApplication = `-- name: CreateApplication :exec
insert into applications (account1_id, account2_id, apply_msg, refuse_msg)
values ( ?, ?, ?, '')
`

type CreateApplicationParams struct {
	Account1ID int64
	Account2ID int64
	ApplyMsg   string
}

// {account1_id:int64,account2_id:int64,apply_msg:string}
func (q *Queries) CreateApplication(ctx context.Context, arg *CreateApplicationParams) error {
	_, err := q.exec(ctx, q.createApplicationStmt, createApplication, arg.Account1ID, arg.Account2ID, arg.ApplyMsg)
	return err
}

const deleteApplication = `-- name: DeleteApplication :exec
delete
from applications
where account1_id = ?
  and account2_id = ?
  and create_at = ?
`

type DeleteApplicationParams struct {
	Account1ID int64
	Account2ID int64
	CreateAt   time.Time
}

// {account1_id:int64,account2_id:int64}
func (q *Queries) DeleteApplication(ctx context.Context, arg *DeleteApplicationParams) error {
	_, err := q.exec(ctx, q.deleteApplicationStmt, deleteApplication, arg.Account1ID, arg.Account2ID, arg.CreateAt)
	return err
}

const existRelation = `-- name: ExistRelation :one
select exists(
    select 1
    from relations
    where (account1_id = ? and account2_id = ?)
       or (account1_id = ? and account2_id = ?)
        for update )
`

type ExistRelationParams struct {
	Account1ID   sql.NullInt64
	Account2ID   sql.NullInt64
	Account1ID_2 sql.NullInt64
	Account2ID_2 sql.NullInt64
}

func (q *Queries) ExistRelation(ctx context.Context, arg *ExistRelationParams) (bool, error) {
	row := q.queryRow(ctx, q.existRelationStmt, existRelation,
		arg.Account1ID,
		arg.Account2ID,
		arg.Account1ID_2,
		arg.Account2ID_2,
	)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const existsApplicationByIDWithLock = `-- name: ExistsApplicationByIDWithLock :one
select exists(
    select 1
    from applications
    where (account1_id = ? and account2_id = ?)
       or (account1_id = ? and account2_id = ?)
        for update )
`

type ExistsApplicationByIDWithLockParams struct {
	Account1ID   int64
	Account2ID   int64
	Account1ID_2 int64
	Account2ID_2 int64
}

// {account1_id:int64,account2_id:int64}
func (q *Queries) ExistsApplicationByIDWithLock(ctx context.Context, arg *ExistsApplicationByIDWithLockParams) (bool, error) {
	row := q.queryRow(ctx, q.existsApplicationByIDWithLockStmt, existsApplicationByIDWithLock,
		arg.Account1ID,
		arg.Account2ID,
		arg.Account1ID_2,
		arg.Account2ID_2,
	)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const getApplicationByID = `-- name: GetApplicationByID :one
select account1_id, account2_id, apply_msg, refuse_msg, status, create_at, update_at
from applications
where account1_id = ? and account2_id = ? and create_at = ?
limit 1
`

type GetApplicationByIDParams struct {
	Account1ID int64
	Account2ID int64
	CreateAt   time.Time
}

// {account1_id:int64,account2_id:int64}
func (q *Queries) GetApplicationByID(ctx context.Context, arg *GetApplicationByIDParams) (*Application, error) {
	row := q.queryRow(ctx, q.getApplicationByIDStmt, getApplicationByID, arg.Account1ID, arg.Account2ID, arg.CreateAt)
	var i Application
	err := row.Scan(
		&i.Account1ID,
		&i.Account2ID,
		&i.ApplyMsg,
		&i.RefuseMsg,
		&i.Status,
		&i.CreateAt,
		&i.UpdateAt,
	)
	return &i, err
}

const getApplications = `-- name: GetApplications :many
select app.account1_id, app.account2_id, app.apply_msg, app.refuse_msg, app.status, app.create_at, app.update_at, app.total,
       a1.name as account1_name,
       a1.avatar as account1_avatar,
       a2.name as account2_name,
       a2.avatar as account2_avatar
from accounts a1,
     accounts a2,
     (select account1_id, account2_id, apply_msg, refuse_msg, status, create_at, update_at, count(*) over () as total
      from applications
      where account1_id = ?
         or account2_id = ?
      order by create_at desc
      limit ? offset ?) as app
where a1.id = app.account1_id
  and a2.id = app.account2_id
`

type GetApplicationsParams struct {
	Account1ID int64
	Account2ID int64
	Limit      int32
	Offset     int32
}

type GetApplicationsRow struct {
	Account1ID     int64
	Account2ID     int64
	ApplyMsg       string
	RefuseMsg      string
	Status         ApplicationsStatus
	CreateAt       time.Time
	UpdateAt       time.Time
	Total          interface{}
	Account1Name   string
	Account1Avatar string
	Account2Name   string
	Account2Avatar string
}

// {account1_id:int64,account2_id:int64,limit:int32,offset:int32,total:int64}
func (q *Queries) GetApplications(ctx context.Context, arg *GetApplicationsParams) ([]*GetApplicationsRow, error) {
	rows, err := q.query(ctx, q.getApplicationsStmt, getApplications,
		arg.Account1ID,
		arg.Account2ID,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetApplicationsRow{}
	for rows.Next() {
		var i GetApplicationsRow
		if err := rows.Scan(
			&i.Account1ID,
			&i.Account2ID,
			&i.ApplyMsg,
			&i.RefuseMsg,
			&i.Status,
			&i.CreateAt,
			&i.UpdateAt,
			&i.Total,
			&i.Account1Name,
			&i.Account1Avatar,
			&i.Account2Name,
			&i.Account2Avatar,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getApplicationsCreatTime = `-- name: GetApplicationsCreatTime :one
select create_at
from applications
where (account1_id = ? and account2_id = ?)
   or (account1_id = ? and account2_id = ?)
order by create_at desc
limit 1
`

type GetApplicationsCreatTimeParams struct {
	Account1ID   int64
	Account2ID   int64
	Account1ID_2 int64
	Account2ID_2 int64
}

func (q *Queries) GetApplicationsCreatTime(ctx context.Context, arg *GetApplicationsCreatTimeParams) (time.Time, error) {
	row := q.queryRow(ctx, q.getApplicationsCreatTimeStmt, getApplicationsCreatTime,
		arg.Account1ID,
		arg.Account2ID,
		arg.Account1ID_2,
		arg.Account2ID_2,
	)
	var create_at time.Time
	err := row.Scan(&create_at)
	return create_at, err
}

const getApplicationsStatus = `-- name: GetApplicationsStatus :one
select status
from applications
where (account1_id = ? and account2_id = ?)
    or (account1_id = ? and account2_id = ?)
order by create_at desc
limit 1
`

type GetApplicationsStatusParams struct {
	Account1ID   int64
	Account2ID   int64
	Account1ID_2 int64
	Account2ID_2 int64
}

func (q *Queries) GetApplicationsStatus(ctx context.Context, arg *GetApplicationsStatusParams) (ApplicationsStatus, error) {
	row := q.queryRow(ctx, q.getApplicationsStatusStmt, getApplicationsStatus,
		arg.Account1ID,
		arg.Account2ID,
		arg.Account1ID_2,
		arg.Account2ID_2,
	)
	var status ApplicationsStatus
	err := row.Scan(&status)
	return status, err
}

const updateApplication = `-- name: UpdateApplication :exec
update applications
set status = ?,
    refuse_msg = ?,
    update_at = now()
where account1_id = ?
  and account2_id = ?
    and create_at = ?
`

type UpdateApplicationParams struct {
	Status     ApplicationsStatus
	RefuseMsg  string
	Account1ID int64
	Account2ID int64
	CreateAt   time.Time
}

// {status:string,refuse_msg:string,account1_id:int64,account2_id:int64}
func (q *Queries) UpdateApplication(ctx context.Context, arg *UpdateApplicationParams) error {
	_, err := q.exec(ctx, q.updateApplicationStmt, updateApplication,
		arg.Status,
		arg.RefuseMsg,
		arg.Account1ID,
		arg.Account2ID,
		arg.CreateAt,
	)
	return err
}
