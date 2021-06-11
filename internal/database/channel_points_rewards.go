// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import (
	"database/sql"
)

type ChannelPointsReward struct {
	ID                               string         `db:"id" json:"id" dbs:"cpr.id"`
	BroadcasterID                    string         `db:"broadcaster_id" json:"broadcaster_id"`
	BroadcasterLogin                 string         `db:"broadcaster_login" dbi:"false" json:"broadcaster_login"`
	BroadcasterName                  string         `db:"broadcaster_name" dbi:"false" json:"broadcaster_name"`
	RewardImage                      sql.NullString `db:"reward_image" json:"-"`
	RealRewardImage                  *string        `json:"image"`
	BackgroundColor                  string         `db:"background_color" json:"background_color"`
	IsEnabled                        *bool          `db:"is_enabled" json:"is_enabled"`
	Cost                             int            `db:"cost" json:"cost"`
	Title                            string         `db:"title" dbs:"cpr.title" json:"title"`
	RewardPrompt                     string         `db:"reward_prompt" json:"prompt"`
	IsUserInputRequired              bool           `db:"is_user_input_required" json:"is_user_input_requird"`
	MaxPerStream                     `json:"max_per_stream_setting"`
	MaxPerUserPerStream              `json:"max_per_user_per_stream_setting"`
	GlobalCooldown                   `json:"global_cooldown_setting"`
	IsPaused                         bool `db:"is_paused" json:"is_paused"`
	IsInStock                        bool `db:"is_in_stock" json:"is_in_stock"`
	DefaultImage                     `dbi:"false" json:"default_image"`
	ShouldRedemptionsSkipQueue       bool           `db:"should_redemptions_skip_queue" json:"should_redemptions_skip_request_queue"`
	RedemptionsRedeemedCurrentStream *int           `db:"redemptions_redeemed_current_stream" json:"redemptions_redeemed_current_stream"`
	CooldownExpiresAt                sql.NullString `db:"cooldown_expires_at" json:"-"`
	RealCooldownExpiresAt            *string        `json:"cooldown_expires_at"`
}

type MaxPerStream struct {
	StreamMaxEnabled bool `db:"stream_max_enabled" json:"is_enabled"`
	StreamMaxCount   int  `db:"stream_max_count" json:"max_per_stream"`
}

type MaxPerUserPerStream struct {
	StreamUserMaxEnabled bool `db:"stream_user_max_enabled" json:"is_enabled"`
	StreamMUserMaxCount  int  `db:"stream_user_max_count" json:"max_per_user_per_stream"`
}

type GlobalCooldown struct {
	GlobalCooldownEnabled bool `db:"global_cooldown_enabled" json:"is_enabled"`
	GlobalCooldownSeconds int  `db:"global_cooldown_seconds" json:"global_cooldown_seconds"`
}

type DefaultImage struct {
	URL1x string `json:"url_1x"`
	URL2x string `json:"url_2x"`
	URL4x string `json:"url_4x"`
}

func (q *Query) GetChannelPointsReward(cpr ChannelPointsReward) (*DBResponse, error) {
	var r []ChannelPointsReward
	sql := generateSQL("select cpr.*,  u1.user_login as broadcaster_login, u1.display_name as broadcaster_name from channel_points_rewards cpr join users u1 on cpr.broadcaster_id = u1.id", cpr, SEP_AND)
	rows, err := q.DB.NamedQuery(sql+q.SQL, cpr)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var cpr ChannelPointsReward
		err := rows.StructScan(&cpr)
		if err != nil {
			return nil, err
		}
		if cpr.CooldownExpiresAt.Valid {
			cpr.RealCooldownExpiresAt = &cpr.CooldownExpiresAt.String
		}
		if cpr.RewardImage.Valid {
			cpr.RealRewardImage = &cpr.RewardImage.String
		}
		cpr.DefaultImage = DefaultImage{
			URL1x: "https://static-cdn.jtvnw.net/custom-reward-images/default-1.png",
			URL2x: "https://static-cdn.jtvnw.net/custom-reward-images/default-2.png",
			URL4x: "https://static-cdn.jtvnw.net/custom-reward-images/default-4.png",
		}
		r = append(r, cpr)
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

func (q *Query) InsertChannelPointsReward(r ChannelPointsReward) error {
	sql := generateInsertSQL("channel_points_rewards", "", r, false)
	_, err := q.DB.NamedExec(sql, r)
	return err
}
func (q *Query) UpdateChannelPointsReward(r ChannelPointsReward) error {
	sql := generateUpdateSQL("channel_points_rewards", []string{"id"}, r)
	_, err := q.DB.NamedExec(sql, r)
	return err
}

func (q *Query) DeleteChannelPointsReward(id string) error {
	tx := q.DB.MustBegin()
	tx.Exec("delete from channel_points_rewards where id=$1", id)
	tx.Exec("update channel_points_redemptions set redemption_status='FULFILLED' where redemption_status='UNFULFILLED' and reward_id=$1", id)
	err := tx.Commit()
	return err
}
