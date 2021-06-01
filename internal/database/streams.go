// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

type Stream struct {
	ID             string      `db:"id" json:"id"`
	UserID         string      `db:"broadcaster_id" json:"user_id"`
	UserLogin      string      `db:"broadcaster_login" json:"user_login"`
	UserName       string      `db:"broadcaster_name" json:"user_name"`
	CategoryID     string      `db:"category_id" json:"game_id"`
	CategoryName   string      `db:"category_name" json:"game_name"`
	StreamType     string      `db:"stream_type" json:"stream_type"`
	Title          string      `db:"title" json:"title"`
	ViewerCount    int         `db:"viewer_count" json:"viewer_count"`
	StartedAt      string      `db:"started_at" json:"started_at"`
	StreamLanguage string      `db:"stream_language" json:"stream_language"`
	IsMature       bool        `db:"is_mature" json:"is_mature"`
	TagIDs         []StreamTag `json:"tag_ids"`
}

type StreamInsert struct {
	ID             string `db:"id" json:"id"`
	UserID         string `db:"broadcaster_id" json:"user_id"`
	CategoryID     string `db:"category_id" json:"game_id"`
	StreamType     string `db:"stream_type" json:"stream_type"`
	Title          string `db:"title" json:"title"`
	ViewerCount    int    `db:"viewer_count" json:"viewer_count"`
	StartedAt      string `db:"started_at" json:"started_at"`
	StreamLanguage string `db:"stream_language" json:"stream_language"`
	IsMature       bool   `db:"is_mature" json:"is_mature"`
}

type StreamTag struct {
	TagID    string `db:"tag_id" json:"is_mature"`
	StreamID string `db:"stream_id" json:"stream_id,omit"`
}

type Tag struct {
	ID   string `db:"id" json:"id"`
	Name string `db:"tag_name" json:"tag_name"`
}

func (c CLIDatabase) GetStream(s Stream) ([]Stream, error) {
	var r []Stream
	sql := generateSQL("select * from streams", s, SEP_AND)
	rows, err := c.DB.NamedQuery(sql, s)
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

func (c CLIDatabase) InsertStream(p StreamInsert, upsert bool) error {
	stmt := generateInsertSQL("streams", "id", p, upsert)
	_, err := c.DB.NamedExec(stmt, p)
	return err
}

func (c CLIDatabase) InsertTag(t Tag) error {
	stmt := generateInsertSQL("tags", "", t, false)
	_, err := c.DB.NamedExec(stmt, t)
	return err
}

func (c CLIDatabase) InsertStreamTag(st StreamTag) error {
	stmt := generateInsertSQL("stream_tags", "", st, false)
	_, err := c.DB.NamedExec(stmt, st)
	return err
}
