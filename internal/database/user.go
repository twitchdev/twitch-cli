// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import (
	"fmt"
	"log"
	"time"

	"github.com/twitchdev/twitch-cli/internal/util"
)

type User struct {
	ID              string `db:"id"`
	UserLogin       string `db:"user_login"`
	DisplayName     string `db:"display_name"`
	Email           string `db:"email"`
	UserType        string `db:"user_type"`
	BroadcasterType string `db:"broadcaster_type"`
	UserDescription string `db:"user_description"`
	CreatedAt       string `db:"created_at"`
	ModifiedAt      string `db:"modified_at"`
}

type Follow struct {
	BroadcasterID    string `db:"to_id" json:"to_id"`
	BroadcasterLogin string `db:"to_login" json:"to_login"`
	BroadcasterName  string `db:"to_name" json:"to_name"`
	ViewerID         string `db:"from_id" json:"from_id"`
	ViewerLogin      string `db:"from_login" json:"from_login"`
	ViewerName       string `db:"from_name" json:"from_name"`
	FollowedAt       string `db:"created_at" json:"followed_at"`
}

type UserRequestParams struct {
	BroadcasterID string `db:"broadcaster_id"`
	UserID        string `db:"user_id"`
	CreatedAt     string `db:"created_at"`
}

type Block struct {
	UserID    string `db:"user_id" son:"user_id"`
	UserLogin string `db:"user_login" json:"user_login"`
	UserName  string `db:"user_name" json:"display_name"`
}

func (q *Query) GetUser(u User) (User, error) {
	var r User
	sql := generateSQL("select * from users", u, SEP_AND)
	sql = fmt.Sprintf("%v LIMIT 1", sql)
	rows, err := q.DB.NamedQuery(sql, u)
	if err != nil {
		return r, err
	}

	for rows.Next() {
		err := rows.StructScan(&r)
		if err != nil {
			return r, err
		}
	}

	return r, err
}

func (q *Query) GetUsers(u User) ([]User, error) {
	var r []User
	sql := generateSQL("select * from users", u, SEP_AND)
	rows, err := q.DB.NamedQuery(sql, u)
	if err != nil {
		return r, err
	}

	for rows.Next() {
		var u User
		err := rows.StructScan(&u)
		if err != nil {
			return r, err
		}
		r = append(r, u)
	}

	return r, err
}

func (q *Query) InsertUser(u User, upsert bool) error {
	stmt := generateInsertSQL("users", "id", u, true)
	_, err := q.DB.NamedExec(stmt, u)
	return err
}

func (q *Query) AddFollow(p UserRequestParams) error {
	stmt := generateInsertSQL("follows", "", p, false)
	p.CreatedAt = util.GetTimestamp().UTC().Format(time.RFC3339)
	_, err := q.DB.NamedExec(stmt, p)
	return err
}

func (q *Query) GetFollows(p UserRequestParams) (*DBResposne, error) {
	db := q.DB
	var r []Follow
	var f Follow
	sql := generateSQL("SELECT u1.id as to_id, u1.user_login as to_login, u1.display_name as to_name, u2.id as from_id, u2.user_login as from_login, u2.display_name as from_name, f.created_at as created_at FROM follows as f JOIN users u1 ON f.broadcaster_id = u1.id JOIN users u2 ON f.user_id = u2.id", p, SEP_AND)
	sql = fmt.Sprintf("%v ORDER BY f.created_at DESC %v", sql, q.SQL)
	rows, err := db.NamedQuery(sql, p)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err := rows.StructScan(&f)
		if err != nil {
			return nil, err
		}
		r = append(r, f)
	}
	var total int
	rows, err = q.DB.NamedQuery(generateSQL("select count(*) from follows", p, SEP_AND), p)
	for rows.Next() {
		err := rows.Scan(&total)
		if err != nil {
			log.Print(err)
			return nil, err
		}
	}

	dbr := DBResposne{
		Data:  r,
		Limit: q.Limit,
		Total: total,
	}

	if len(r) != q.Limit {
		q.PaginationCursor = ""
	}

	dbr.Cursor = q.PaginationCursor

	return &dbr, err
}

func (q *Query) DeleteFollow(from_user string, to_user string) error {
	_, err := q.DB.Exec(`delete from follows where broadcaster_id=$1 and user_id=$2`, to_user, from_user)
	return err
}

func (q *Query) AddBlock(p UserRequestParams) error {
	stmt := generateInsertSQL("blocks", "id", p, false)
	p.CreatedAt = util.GetTimestamp().UTC().Format(time.RFC3339)
	_, err := q.DB.NamedExec(stmt, p)
	return err
}

func (q *Query) GetBlocks(p UserRequestParams) (*DBResposne, error) {
	var r []Block
	sql := generateSQL("SELECT u1.id as user_id, u1.user_login as user_login, u1.display_name as user_name FROM blocks as b JOIN users u1 ON b.user_id = u1.id", p, SEP_AND)
	sql = fmt.Sprintf("%v ORDER BY f.created_at DESC", sql)

	rows, err := q.DB.NamedQuery(sql, p)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		b := Block{}
		err := rows.StructScan(&b)
		if err != nil {
			return nil, err
		}
		r = append(r, b)
	}
	dbr := DBResposne{
		Data:  r,
		Limit: q.Limit,
		Total: len(r),
	}

	if len(r) != q.Limit {
		q.PaginationCursor = ""
	}

	dbr.Cursor = q.PaginationCursor

	return &dbr, err
}

func (q *Query) DeleteBlock(from_user string, to_user string) error {
	_, err := q.DB.Exec(`delete from blocks where broadcaster_id=$1 and user_id=$2`, to_user, from_user)
	return err
}
