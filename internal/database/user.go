// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/twitchdev/twitch-cli/internal/util"
)

type User struct {
	ID              string         `db:"id" json:"id" dbs:"u1.id"`
	UserLogin       string         `db:"user_login" json:"login"`
	DisplayName     string         `db:"display_name" json:"display_name"`
	Email           string         `db:"email" json:"email,omitempty"`
	UserType        string         `db:"user_type" json:"type"`
	BroadcasterType string         `db:"broadcaster_type" json:"broadcaster_type"`
	UserDescription string         `db:"user_description" json:"description"`
	CreatedAt       string         `db:"created_at" json:"created_at"`
	ModifiedAt      string         `db:"modified_at" json:"-"`
	ProfileImageURL string         `dbi:"false" json:"profile_image_url" `
	OfflineImageURL string         `dbi:"false" json:"offline_image_url" `
	ViewCount       int            `dbi:"false" json:"view_count"`
	CategoryID      sql.NullString `db:"category_id" json:"game_id" dbi:"force"`
	CategoryName    sql.NullString `db:"category_name" json:"game_name" dbi:"false"`
	Title           string         `db:"title" json:"title"`
	Language        string         `db:"stream_language" json:"stream_language"`
	Delay           int            `db:"delay" json:"delay" dbi:"force"`
	ChatColor       string         `db:"chat_color" json:"-"`
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
	UserID    string `db:"user_id" json:"user_id"`
	UserLogin string `db:"user_login" json:"user_login"`
	UserName  string `db:"user_name" json:"display_name"`
}

type Editor struct {
	UserID    string `db:"user_id" json:"user_id"`
	UserLogin string `db:"user_login" json:"-"`
	UserName  string `db:"user_name" json:"user_name"`
	CreatedAt string `db:"created_at" json:"created_at"`
}

type SearchChannel struct {
	ID           string   `db:"id" json:"id" dbs:"u1.id"`
	UserLogin    string   `db:"user_login" json:"broadcaster_login"`
	DisplayName  string   `db:"display_name" json:"display_name"`
	CategoryID   *string  `db:"category_id" json:"game_id" dbi:"false"`
	CategoryName *string  `db:"category_name" json:"game_name" dbi:"false"`
	Title        string   `db:"title" json:"title"`
	Language     string   `db:"stream_language" json:"broadcaster_language"`
	TagIDs       []string `json:"tag_ids" dbi:"false"`
	IsLive       bool     `json:"is_live" db:"is_live"`
	StartedAt    *string  `db:"started_at" json:"started_at"`
	// calculated fields
	ThumbNailURL string `json:"thumbnail_url"`
}

type VIP struct {
	BroadcasterID string `db:"broadcaster_id"`
	UserID        string `db:"user_id"`
	CreatedAt     string `db:"created_at"`
}

func (q *Query) GetUser(u User) (User, error) {
	var r User
	sql := generateSQL("select * from users u1", u, SEP_AND)
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
		r.OfflineImageURL = "https://static-cdn.jtvnw.net/jtv_user_pictures/3f13ab61-ec78-4fe6-8481-8682cb3b0ac2-channel_offline_image-1920x1080.png"
		r.ProfileImageURL = "https://static-cdn.jtvnw.net/jtv_user_pictures/8a6381c7-d0c0-4576-b179-38bd5ce1d6af-profile_image-300x300.png"
	}

	return r, err
}

func (q *Query) GetUsers(u User) (*DBResponse, error) {
	var r []User
	sql := generateSQL("select * from users u1", u, SEP_AND)
	rows, err := q.DB.NamedQuery(sql+q.SQL, u)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var u User
		err := rows.StructScan(&u)
		if err != nil {
			return nil, err
		}
		r = append(r, u)
	}

	dbr := DBResponse{
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

func (q *Query) GetChannels(u User) (*DBResponse, error) {
	var r []User
	sql := generateSQL("select u1.*, c.category_name from users u1 left join categories c on u1.category_id = c.id", u, SEP_AND)
	rows, err := q.DB.NamedQuery(sql+q.SQL, u)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var u User
		err := rows.StructScan(&u)
		if err != nil {
			return nil, err
		}
		r = append(r, u)
	}

	dbr := DBResponse{
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

func (q *Query) InsertUser(u User, upsert bool) error {
	stmt := generateInsertSQL("users", "id", u, upsert)
	_, err := q.DB.NamedExec(stmt, u)
	return err
}

func (q *Query) AddFollow(p UserRequestParams) error {
	stmt := generateInsertSQL("follows", "", p, false)
	p.CreatedAt = util.GetTimestamp().UTC().Format(time.RFC3339)
	_, err := q.DB.NamedExec(stmt, p)
	return err
}

func (q *Query) GetFollows(p UserRequestParams) (*DBResponse, error) {
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

	dbr := DBResponse{
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

func (q *Query) GetBlocks(p UserRequestParams) (*DBResponse, error) {
	var r []Block
	sql := generateSQL("SELECT u1.id as user_id, u1.user_login as user_login, u1.display_name as user_name FROM blocks as b JOIN users u1 ON b.user_id = u1.id", p, SEP_AND)
	sql = fmt.Sprintf("%v ORDER BY b.created_at DESC", sql)

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
	dbr := DBResponse{
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

func (q *Query) UpdateChannel(id string, u User) error {
	sql := generateUpdateSQL("users", []string{"id"}, u)
	_, err := q.DB.NamedExec(sql, u)
	return err
}

func (q *Query) AddEditor(p UserRequestParams) error {
	stmt := generateInsertSQL("editors", "id", p, false)
	p.CreatedAt = util.GetTimestamp().UTC().Format(time.RFC3339)
	_, err := q.DB.NamedExec(stmt, p)
	return err
}

func (q *Query) GetEditors(u User) (*DBResponse, error) {
	var r []Editor

	err := q.DB.Select(&r, "SELECT u1.id as user_id, u1.user_login as user_login, u1.display_name as user_name, e.created_at FROM editors as e JOIN users u1 ON e.user_id = u1.id where broadcaster_id = $1 ORDER BY e.created_at DESC", u.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	dbr := DBResponse{
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

func (q *Query) SearchChannels(query string, live_only bool) (*DBResponse, error) {
	r := []SearchChannel{}
	stmt := `select u1.id, u1.user_login, u1.display_name, u1.category_id, u1.title, u1.stream_language, c.category_name, case when s.id is null then 'false' else 'true' end is_live, s.started_at from users u1 left join streams s on u1.id = s.broadcaster_id left join categories c on u1.category_id = c.id where lower(u1.user_login) like lower($1)`

	if live_only {
		stmt = `select u1.id, u1.user_login, u1.display_name, u1.category_id, u1.title, u1.stream_language, c.category_name, case when s.id is null then 'false' else 'true' end is_live, s.started_at from users u1 left join streams s on u1.id = s.broadcaster_id left join categories c on u1.category_id = c.id where lower(u1.user_login) like lower($1) and is_live='true'`
	}

	err := q.DB.Select(&r, stmt+q.SQL, fmt.Sprintf("%%%v%%", query))
	if err != nil {
		return nil, err
	}

	for i, c := range r {
		st := []string{}
		err = q.DB.Select(&st, "select tag_id from stream_tags where user_id=$1", c.ID)
		if err != nil {
			return nil, err
		}

		emptyString := ""
		if c.StartedAt == nil {
			r[i].StartedAt = &emptyString
		}
		r[i].TagIDs = st
		r[i].ThumbNailURL = "https://static-cdn.jtvnw.net/jtv_user_pictures/3f13ab61-ec78-4fe6-8481-8682cb3b0ac2-channel_offline_image-300x300.png"
	}

	dbr := DBResponse{
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

func (q *Query) GetVIPsByBroadcaster(broadcaster string) (*DBResponse, error) {
	var r []VIP

	err := q.DB.Select(&r, "SELECT * FROM vips WHERE broadcaster_id=$1", broadcaster)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	dbr := DBResponse{
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

func (q *Query) AddVIP(p UserRequestParams) error {
	stmt := generateInsertSQL("vips", "user_id", p, false)
	p.CreatedAt = util.GetTimestamp().UTC().Format(time.RFC3339)
	_, err := q.DB.NamedExec(stmt, p)
	return err
}

func (q *Query) DeleteVIP(broadcaster string, user string) error {
	_, err := q.DB.Exec(`delete from vips where broadcaster_id=$1 and user_id=$2`, broadcaster, user)
	return err
}
