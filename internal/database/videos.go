// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import "log"

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
	Type             string              `json:"type"`
	URL              string              `json:"url"`
	ThumbnailURL     string              `json:"thumbnail_url"`
}

type VideoMutedSegment struct {
	VideoID     string `db:"video_id" json:"-"`
	VideoOffset int    `db:"video_offset" json:"video_offset"`
	Duration    int    `db:"duration" json:"duration"`
}

func (q *Query) GetVideos(v Video) ([]Video, error) {
	var r []Video
	sql := generateSQL("select v.*, u1.user_login as broadcaster_login, u1.display_name as broadcaster_name from videos v join users u1 on v.broadcaster_id = u1.id", v, SEP_AND)
	rows, err := q.DB.NamedQuery(sql, v)
	if err != nil {
		log.Print(err)
		return r, err
	}

	for rows.Next() {
		var v Video
		var vms []VideoMutedSegment
		err := rows.StructScan(&v)
		if err != nil {
			log.Print(err)
			return r, err
		}
		err = q.DB.Select(&vms, "select * from video_muted_segments where video_id=$1", v.ID)
		if err != nil {
			log.Print(err)
			return r, err
		}
		v.MutedSegments = vms
		r = append(r, v)
	}
	return r, nil
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
