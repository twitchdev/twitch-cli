// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import (
	"database/sql"
)

type ChannelPointsRedemption struct {
	ID                                string         `db:"id" json:"id" dbs:"cpr.id"`
	BroadcasterID                     string         `db:"broadcaster_id" json:"broadcaster_id" dbs:"cpr.broadcaster_id"`
	BroadcasterLogin                  string         `db:"broadcaster_login" dbi:"false" json:"broadcaster_login"`
	BroadcasterName                   string         `db:"broadcaster_name" dbi:"false" json:"broadcaster_name"`
	UserID                            string         `db:"user_id" json:"user_id"`
	UserLogin                         string         `db:"user_login" dbi:"false" json:"user_login"`
	UserName                          string         `db:"user_name" dbi:"false" json:"user_name"`
	UserInput                         sql.NullString `db:"user_input" json:"-"`
	RealUserInput                     string         `json:"user_input"`
	RedemptionStatus                  string         `db:"redemption_status" json:"status"`
	RedeemedAt                        string         `db:"redeemed_at" json:"redeemed_at"`
	RewardID                          string         `db:"reward_id" json:"-"`
	ChannelPointsRedemptionRewardInfo `json:"reward"`
}

type ChannelPointsRedemptionRewardInfo struct {
	ID           string `dbi:"false" db:"red_id" json:"id" dbs:"red.id"`
	Title        string `dbi:"false" db:"title" json:"title"`
	RewardPrompt string `dbi:"false" db:"reward_prompt" json:"prompt"`
	Cost         int    `dbi:"false" db:"cost" json:"cost"`
}

func (q *Query) GetChannelPointsRedemption(cpr ChannelPointsRedemption, sort string) (*DBResponse, error) {
	var r []ChannelPointsRedemption
	orderBy := ""
	if sort == "" || sort == "OLDEST" {
		orderBy = "asc"
	} else if sort == "NEWEST" {
		orderBy = "desc"
	}

	sql := generateSQL("select cpr.*, u1.user_login as broadcaster_login, u1.display_name as broadcaster_name, u2.user_login, u2.display_name as user_name, red.id as red_id, red.title, red.cost, red.reward_prompt from channel_points_redemptions cpr join users u1 on cpr.broadcaster_id = u1.id join users u2 on cpr.user_id = u2.id join channel_points_rewards red on cpr.reward_id = red.id", cpr, SEP_AND)
	rows, err := q.DB.NamedQuery(sql+" order by cpr.redeemed_at "+orderBy+q.SQL, cpr)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var red ChannelPointsRedemption
		err := rows.StructScan(&red)
		if err != nil {
			return nil, err
		}
		red.RealUserInput = red.UserInput.String
		if !red.UserInput.Valid {
			red.RealUserInput = ""
		}
		r = append(r, red)
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

func (q *Query) InsertChannelPointsRedemption(r ChannelPointsRedemption) error {
	sql := generateInsertSQL("channel_points_redemptions", "", r, false)
	_, err := q.DB.NamedExec(sql, r)
	return err
}

func (q *Query) UpdateChannelPointsRedemption(r ChannelPointsRedemption) error {
	sql := generateUpdateSQL("channel_points_redemptions", []string{"id"}, r)
	_, err := q.DB.NamedExec(sql, r)
	return err
}
