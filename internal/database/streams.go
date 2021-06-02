// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import "log"

type Stream struct {
	ID             string   `db:"id" json:"id" dbs:"s.id"`
	UserID         string   `db:"broadcaster_id" json:"user_id"`
	UserLogin      string   `db:"broadcaster_login" json:"user_login" dbi:"false"`
	UserName       string   `db:"broadcaster_name" json:"user_name" dbi:"false"`
	CategoryID     string   `db:"category_id" json:"game_id"`
	CategoryName   string   `db:"category_name" json:"game_name" dbi:"false"`
	StreamType     string   `db:"stream_type" json:"stream_type"`
	Title          string   `db:"title" json:"title"`
	ViewerCount    int      `db:"viewer_count" json:"viewer_count"`
	StartedAt      string   `db:"started_at" json:"started_at"`
	StreamLanguage string   `db:"stream_language" json:"stream_language"`
	IsMature       bool     `db:"is_mature" json:"is_mature"`
	TagIDs         []string `json:"tag_ids" dbi:"false"`
}

type StreamTag struct {
	TagID    string `db:"tag_id" json:"tag_id"`
	StreamID string `db:"stream_id" json:"-"`
}

type Tag struct {
	ID     string `db:"id" json:"id"`
	IsAuto bool   `db:"is_auto" dbi:"false" json:"is_auto"`
	Name   string `db:"tag_name" json:"tag_name"`
}

func (q *Query) GetStream(s Stream) (*DBResposne, error) {
	var r = []Stream{}
	sql := generateSQL("select s.*, u1.user_login as broadcaster_login, u1.display_name as broadcaster_name, c.category_name from streams s join users u1 on s.broadcaster_id = u1.id join categories c on c.id = s.category_id", s, SEP_AND)
	rows, err := q.DB.NamedQuery(sql, s)
	log.Print(sql)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	for rows.Next() {
		var s Stream
		var st []string
		err := rows.StructScan(&s)
		if err != nil {
			return nil, err
		}
		err = q.DB.Select(&st, "select tag_id from stream_tags where stream_id=$1", s.ID)
		if err != nil {
			return nil, err
		}
		s.TagIDs = st
		r = append(r, s)
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

func (q *Query) InsertStream(p Stream, upsert bool) error {
	stmt := generateInsertSQL("streams", "id", p, upsert)
	_, err := q.DB.NamedExec(stmt, p)
	return err
}

func (q *Query) GetTags(t Tag) ([]Tag, error) {
	r := []Tag{}
	sql := generateSQL("select * from tags", t, SEP_AND)
	rows, err := q.DB.NamedQuery(sql, t)

	for rows.Next() {
		var t Tag
		err := rows.StructScan(&t)
		if err != nil {
			return r, err
		}
		r = append(r, t)
	}
	return r, err
}

func (q *Query) InsertTag(t Tag) error {
	stmt := generateInsertSQL("tags", "", t, false)
	_, err := q.DB.NamedExec(stmt, t)
	return err
}

func (q *Query) InsertStreamTag(st StreamTag) error {
	stmt := generateInsertSQL("stream_tags", "", st, false)
	_, err := q.DB.NamedExec(stmt, st)
	return err
}
