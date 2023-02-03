// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import (
	"database/sql"
	"fmt"
	"log"
)

type Stream struct {
	ID          string   `db:"id" json:"id" dbs:"s.id"`
	UserID      string   `db:"broadcaster_id" json:"user_id"`
	UserLogin   string   `db:"broadcaster_login" json:"user_login" dbi:"false"`
	UserName    string   `db:"broadcaster_name" json:"user_name" dbi:"false"`
	StreamType  string   `db:"stream_type" json:"type"`
	ViewerCount int      `db:"viewer_count" json:"viewer_count"`
	StartedAt   string   `db:"started_at" json:"started_at"`
	IsMature    bool     `db:"is_mature" json:"is_mature"`
	TagIDs      []string `json:"tag_ids" dbi:"false"`
	Tags        []string `json:"tags" dbi:"false"`
	// stored in users, but pulled here for json parsing
	CategoryID       sql.NullString `db:"category_id" json:"-" dbi:"false"`
	RealCategoryID   string         `json:"game_id"`
	CategoryName     sql.NullString `db:"category_name" json:"-" dbi:"false"`
	RealCategoryName string         `json:"game_name"`
	Title            string         `db:"title" json:"title" dbi:"false"`
	Language         string         `db:"stream_language" json:"language" dbi:"false"`
	// calculated fields
	ThumbnailURL string `json:"thumbnail_url"`
}

type StreamTag struct {
	TagID  string `db:"tag_id" json:"tag_id"`
	UserID string `db:"user_id" json:"-"`
}

type Tag struct {
	ID     string `db:"id" json:"id"`
	IsAuto bool   `db:"is_auto" dbi:"false" json:"is_auto"`
	Name   string `db:"tag_name" json:"tag_name"`
}

type StreamMarkerUser struct {
	BroadcasterID    string              `db:"broadcaster_id" json:"user_id"`
	BroadcasterLogin string              `db:"broadcaster_login" json:"user_login" dbi:"false"`
	BroadcasterName  string              `db:"broadcaster_name" json:"user_name" dbi:"false"`
	Videos           []StreamMarkerVideo `json:"videos"`
}

type StreamMarkerVideo struct {
	VideoID string         `db:"video_id" dbs:"v.id" json:"video_id"`
	Markers []StreamMarker `json:"markers"`
}

type StreamMarker struct {
	ID              string `db:"id" dbs:"sm.id" json:"id"`
	CreatedAt       string `db:"created_at" json:"created_at"`
	PositionSeconds int    `db:"position_seconds" json:"position_seconds"`
	Description     string `db:"description" json:"description"`
	BroadcasterID   string `db:"broadcaster_id" json:"-"`
	VideoID         string `db:"video_id" dbs:"v.id" json:"-"`
	URL             string `json:"URL"`
}

func (q *Query) GetStream(s Stream) (*DBResponse, error) {
	var r = []Stream{}
	sql := generateSQL("select s.*, u1.user_login as broadcaster_login, u1.display_name as broadcaster_name, u1.category_id as category_id, c.category_name, u1.stream_language as stream_language, u1.title as title from streams s join users u1 on s.broadcaster_id = u1.id left join categories c on c.id = u1.category_id", s, SEP_AND)
	rows, err := q.DB.NamedQuery(sql+q.SQL, s)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	for rows.Next() {
		var s Stream
		err := rows.StructScan(&s)

		if err != nil {
			return nil, err
		}
		if s.CategoryID.Valid {
			s.RealCategoryID = s.CategoryID.String
		}
		if s.CategoryName.Valid {
			s.RealCategoryName = s.CategoryName.String
		}
		s.ThumbnailURL = fmt.Sprintf("https://static-cdn.jtvnw.net/previews-ttv/live_user_%v-{width}x{height}.jpg", s.UserLogin)
		r = append(r, s)
	}

	for i := range r {
		r[i].TagIDs = []string{} // Needs to be removed from db when this is fully removed from API
		r[i].Tags = []string{"English", "CLI Tag"}
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

func (q *Query) InsertStream(p Stream, upsert bool) error {
	stmt := generateInsertSQL("streams", "id", p, upsert)
	_, err := q.DB.NamedExec(stmt, p)
	return err
}

func (q *Query) GetTags(t Tag) (*DBResponse, error) {
	r := []Tag{}
	sql := generateSQL("select * from tags", t, SEP_AND)
	rows, err := q.DB.NamedQuery(sql+q.SQL, t)

	for rows.Next() {
		var t Tag
		err := rows.StructScan(&t)
		if err != nil {
			return nil, err
		}
		r = append(r, t)
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

func (q *Query) GetStreamTags(id string) (*DBResponse, error) {
	r := []Tag{}
	err := q.DB.Select(&r, "select t.* from tags t join stream_tags st on st.tag_id = t.id where st.user_id=$1", id)
	if err != nil {
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

func (q *Query) DeleteAllStreamTags(userID string) error {
	_, err := q.DB.Exec("delete from stream_tags where user_id = $1", userID)
	return err
}

func (q *Query) GetFollowedStreams(userID string) (*DBResponse, error) {
	var r = []Stream{}
	sql := "select s.*, u1.user_login as broadcaster_login, u1.display_name as broadcaster_name, u1.category_id as category_id, c.category_name, u1.stream_language as stream_language, u1.title as title from streams s join users u1 on s.broadcaster_id = u1.id left join categories c on c.id = u1.category_id join follows f on f.broadcaster_id = s.broadcaster_id where f.user_id = $1"

	err := q.DB.Select(&r, sql+q.SQL, userID)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	for i, s := range r {
		var st []string
		if err != nil {
			return nil, err
		}
		if s.CategoryID.Valid {
			r[i].RealCategoryID = s.CategoryID.String
		}
		if s.CategoryName.Valid {
			r[i].RealCategoryName = s.CategoryName.String
		}
		err = q.DB.Select(&st, "select tag_id from stream_tags where user_id=$1", s.UserID)
		if err != nil {
			return nil, err
		}
		r[i].TagIDs = st
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

func (q *Query) GetStreamMarkers(sm StreamMarker) (*DBResponse, error) {
	r := []StreamMarkerUser{}
	stmt := generateSQL("select u1.id as broadcaster_id, u1.user_login as broadcaster_login, u1.display_name as broadcaster_name from users u1 join videos v on v.broadcaster_id = u1.id", sm, SEP_AND)
	rows, err := q.DB.NamedQuery(stmt+" limit 1", sm)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var u StreamMarkerUser
		err := rows.StructScan(&u)
		if err != nil {
			return nil, err
		}
		r = append(r, u)
	}

	for i, u := range r {
		sm := []StreamMarker{}
		video := []Video{}
		err := q.DB.Select(&video, "select v.id from videos v where v.broadcaster_id=$1 order by v.created_at desc limit 1", r[i].BroadcasterID)
		if err != nil {
			return nil, err
		}

		for _, v := range video {
			err := q.DB.Select(&sm, "select sm.*, v.id as video_id from stream_markers sm join videos v on sm.video_id = v.id where v.id = $1 order by sm.position_seconds asc", v.ID)
			if err != nil {
				return nil, err
			}
			for i := range sm {
				sm[i].URL = fmt.Sprintf("https://twitch.tv/%v/manager/highlighter/%v?t=%v", u.BroadcasterLogin, v.ID, calcTOffset(sm[i].PositionSeconds))
			}

			r[i].Videos = append(r[i].Videos, StreamMarkerVideo{VideoID: v.ID, Markers: sm})
		}

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

func (q *Query) InsertStreamMarker(sm StreamMarker) error {
	stmt := generateInsertSQL("stream_markers", "", sm, false)
	_, err := q.DB.NamedExec(stmt, sm)
	return err
}

func calcTOffset(offset int) string {
	hours := offset / (60 * 60)
	minutes := (offset % (60 * 60)) / 60
	seconds := (offset % (60 * 60)) % 60
	return fmt.Sprintf("%vh%vm%vs", hours, minutes, seconds)
}
