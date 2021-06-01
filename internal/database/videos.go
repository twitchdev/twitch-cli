// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

type Video struct {
	ID               string `db:"id" json:"id"`
	StreamID         string `db:"stream_id" json:"stream_id"`
	BroadcasterID    string `db:"broadcaster_id" json:"user_id"`
	BroadcasterLogin string `json:"user_login"`
	BroadcasterName  string `json:"user_name"`
	Title            string `db:"title" json:"title"`
	VideoDescription string `db:"video_description" json:"video_description"`
	CreatedAt        string `db:"created_at" json:"created_at"`
	PublishedAt      string `db:"published_at" json:"published_at"`
	Viewable         string `db:"viewable" json:"viewable"`
	ViewCount        int    `db:"view_count" json:"view_count"`
	Duration         string `db:"duration" json:"duration"`
	VideoLanguage    string `db:"video_language" json:"video_language"`
}

func (c CLIDatabase) GetVideos(v Video) ([]Video, error) {
	var r []Video

	sql := generateSQL("select * from videos", v, SEP_AND)
	rows, err := c.DB.NamedQuery(sql, v)
	if err != nil {
		return r, err
	}

	for rows.Next() {
		err := rows.StructScan(&r)
		if err != nil {
			return r, err
		}
	}

	return r, nil
}

func (c CLIDatabase) InsertVideo(v Video) error {
	stmt := generateInsertSQL("videos", "id", v, false)
	_, err := c.DB.NamedExec(stmt, v)
	return err
}
