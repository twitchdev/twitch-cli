// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import (
	"fmt"
	"log"
)

type Video struct {
	ID               string              `db:"id" json:"id" dbs:"v.id"`
	StreamID         string              `db:"stream_id" json:"stream_id"`
	BroadcasterID    string              `db:"broadcaster_id" json:"user_id"`
	BroadcasterLogin string              `db:"broadcaster_login" json:"user_login" dbi:"false"`
	BroadcasterName  string              `db:"broadcaster_name" json:"user_name" dbi:"false"`
	Title            string              `db:"title" json:"title"`
	VideoDescription string              `db:"video_description" json:"video_description"`
	CreatedAt        string              `db:"created_at" json:"created_at"`
	PublishedAt      string              `db:"published_at" json:"published_at"`
	Viewable         string              `db:"viewable" json:"viewable"`
	ViewCount        int                 `db:"view_count" json:"view_count"`
	Duration         string              `db:"duration" json:"duration"`
	VideoLanguage    string              `db:"video_language" json:"video_language"`
	MutedSegments    []VideoMutedSegment `json:"muted_segments"`
	// calculated fields
	Type         string `json:"type"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnail_url"`
}

type VideoMutedSegment struct {
	VideoID     string `db:"video_id" json:"-"`
	VideoOffset int    `db:"video_offset" json:"video_offset"`
	Duration    int    `db:"duration" json:"duration"`
}

type Clip struct {
	ID              string  `db:"id" json:"id" dbs:"c.id"`
	BroadcasterID   string  `db:"broadcaster_id" json:"broadcaster_id"`
	BroadcasterName string  `db:"broadcaster_name" json:"broadcaster_name" dbi:"false"`
	CreatorID       string  `db:"creator_id" json:"creator_id"`
	CreatorName     string  `db:"creator_name" json:"creator_login" dbi:"false"`
	VideoID         string  `db:"video_id" json:"video_id"`
	GameID          string  `db:"game_id" json:"game_id"`
	Language        string  `db:"language" dbi:"false" json:"language"`
	Title           string  `db:"title" json:"title"`
	ViewCount       int     `db:"view_count" json:"view_count"`
	CreatedAt       string  `db:"created_at" json:"created_at"`
	Duration        float64 `db:"duration" json:"duration"`
	// calculated fields
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnail_url"`
	EmbedURL     string `json:"embed_urL"`
	StartedAt    string `db:"started_at" dbi:"false" json:"-"`
	EndedAt      string `db:"ended_at" dbi:"false" json:"-"`
}

func (q *Query) GetVideos(v Video) (*DBResposne, error) {
	var r []Video
	sql := generateSQL("select v.*, u1.user_login as broadcaster_login, u1.display_name as broadcaster_name from videos v join users u1 on v.broadcaster_id = u1.id", v, SEP_AND)
	rows, err := q.DB.NamedQuery(sql, v)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	for rows.Next() {
		var v Video
		var vms []VideoMutedSegment
		err := rows.StructScan(&v)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		err = q.DB.Select(&vms, "select * from video_muted_segments where video_id=$1", v.ID)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		v.MutedSegments = vms
		r = append(r, v)
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

func (q *Query) InsertVideo(v Video) error {
	stmt := generateInsertSQL("videos", "", v, false)
	_, err := q.DB.NamedExec(stmt, v)
	return err
}

func (q *Query) DeleteVideo(id string) error {
	_, err := q.DB.Exec("delete from videos where id = $1", id)
	return err
}

func (q *Query) InsertMutedSegmentsForVideo(vms VideoMutedSegment) error {
	stmt := generateInsertSQL("video_muted_segments", "", vms, false)
	_, err := q.DB.NamedExec(stmt, vms)
	return err
}

func (q *Query) InsertClip(c Clip) error {
	stmt := generateInsertSQL("clips", "", c, false)
	_, err := q.DB.NamedExec(stmt, c)
	return err
}

func (q *Query) GetClips(c Clip, startDate string, endDate string) (*DBResposne, error) {
	var r []Clip
	sql := generateSQL("select c.*,  u1.display_name as broadcaster_name, u1.stream_language as language, u2.display_name as creator_name from clips c join users u1 on c.broadcaster_id = u1.id join users u2 on c.creator_id = u2.id ", c, SEP_AND)
	if startDate != "" {
		sql += fmt.Sprintf(" and c.created_at > :started_at and c.created_at < :ended_at ")
	}
	sql += q.SQL
	println(sql)
	rows, err := q.DB.NamedQuery(sql, c)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	for rows.Next() {
		var c Clip
		err := rows.StructScan(&c)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		c.EmbedURL = fmt.Sprintf("https://clips.twitch.tv/embed?clip=%v", c.ID)
		c.ThumbnailURL = "https://clips-media-assets2.twitch.tv/157589949-preview-480x272.jpg"
		c.URL = fmt.Sprintf("https://clips.twitch.tv/%v", c.ID)
		r = append(r, c)
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
