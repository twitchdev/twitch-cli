// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import "fmt"

type Team struct {
	ID                 string       `db:"id" json:"id"`
	BackgroundImageUrl *string      `db:"background_image_url" json:"background_image_url"`
	Banner             *string      `db:"banner" json:"banner"`
	CreatedAt          string       `db:"created_at" json:"created_at"`
	UpdatedAt          string       `db:"updated_at" json:"updated_at"`
	Info               string       `db:"info" json:"info"`
	ThumbnailURL       string       `db:"thumbnail_url" json:"thumbnail_url"`
	TeamName           string       `db:"team_name" json:"team_name"`
	TeamDisplayName    string       `db:"team_display_name" json:"team_display_name"`
	Users              []TeamMember `json:"users"`
}

type TeamMember struct {
	TeamID    string `db:"team_id" json:"-"`
	UserID    string `db:"user_id" json:"user_id"`
	UserName  string `db:"user_name" json:"user_name" dbi:"false"`
	UserLogin string `db:"user_login" json:"user_login" dbi:"false"`
}

func (q *Query) GetTeam(t Team) (*DBResponse, error) {
	var r []Team
	sql := generateSQL("select * from teams", t, SEP_AND)
	rows, err := q.DB.NamedQuery(sql, t)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var t Team
		err := rows.StructScan(&t)
		if err != nil {
			return nil, err
		}

		if t.BackgroundImageUrl == nil {
			t.BackgroundImageUrl = nil
		}
		if t.Banner == nil {
			t.Banner = nil
		}

		r = append(r, t)
	}

	for i, t := range r {
		p := TeamMember{TeamID: t.ID}
		tms := []TeamMember{}

		err = q.DB.Select(&tms, "select u1.id as user_id, u1.user_login as user_login, u1.display_name as user_name FROM team_members as tm JOIN users u1 ON tm.user_id = u1.id where tm.team_id=$1", p.TeamID)
		if err != nil {
			return nil, err
		}
		r[i].Users = tms
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

func (q *Query) InsertTeam(t Team) error {
	stmt := generateInsertSQL("teams", "", t, false)
	_, err := q.DB.NamedExec(stmt, t)
	return err
}

func (q *Query) InsertTeamMember(tm TeamMember) error {
	stmt := generateInsertSQL("team_members", "", tm, false)
	_, err := q.DB.NamedExec(stmt, tm)
	return err
}

func (q *Query) GetTeamByBroadcaster(broadcasterID string) (*DBResponse, error) {
	var r []Team
	err := q.DB.Select(&r, "select t.* from teams t join team_members tm on tm.team_id = t.id where tm.user_id = $1", broadcasterID)
	if err != nil {
		return nil, err
	}

	for i, t := range r {
		if t.BackgroundImageUrl == nil {
			r[i].BackgroundImageUrl = nil
		}
		if t.Banner == nil {
			r[i].Banner = nil
		}
		r[i].ThumbnailURL = fmt.Sprintf("https://static-cdn.jtvnw.net/jtv_user_pictures/team-%v-team_logo_image-bf1d9a87ca81432687de60e24ad9593d-600x600.png", t.TeamName)
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
