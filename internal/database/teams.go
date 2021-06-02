// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

type Team struct {
	ID                 string       `db:"id" json:"id"`
	BackgroundImageUrl string       `db:"background_image_url" json:"background_image_url"`
	Banner             string       `db:"banner" json:"banner"`
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

func (q *Query) GetTeam(t Team) ([]Team, error) {
	var r []Team
	sql := generateSQL("select * from teams", t, SEP_AND)
	rows, err := q.DB.NamedQuery(sql, t)
	if err != nil {
		return r, err
	}

	for rows.Next() {
		var t Team
		err := rows.StructScan(&t)
		if err != nil {
			return r, err
		}
		p := TeamMember{TeamID: t.ID}
		tms := []TeamMember{}

		err = q.DB.Select(&tms, "select u1.id as user_id, u1.user_login as user_login, u1.display_name as user_name FROM team_members as tm JOIN users u1 ON tm.user_id = u1.id where tm.team_id=$1", p.TeamID)
		if err != nil {
			return r, err
		}
		t.Users = tms
		r = append(r, t)
	}

	return r, err
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
