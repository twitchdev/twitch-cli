// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

type ChatSettings struct {
	BroadcasterID                 string `db:"broadcaster_id" json:"broadcaster_id"`
	SlowMode                      *bool  `db:"slow_mode" json:"slow_mode"`
	SlowModeWaitTime              *int   `db:"slow_mode_wait_time" json:"slow_mode_wait_time"`
	FollowerMode                  *bool  `db:"follower_mode" json:"follower_mode"`
	FollowerModeDuration          *int   `db:"follower_mode_duration" json:"follower_mode_duration"`
	SubscriberMode                *bool  `db:"subscriber_mode" json:"subscriber_mode"`
	EmoteMode                     *bool  `db:"emote_mode" json:"emote_mode"`
	UniqueChatMode                *bool  `db:"unique_chat_mode" json:"unique_chat_mode"`
	NonModeratorChatDelay         *bool  `db:"non_moderator_chat_delay" json:"non_moderator_chat_delay"`
	NonModeratorChatDelayDuration *int   `db:"non_moderator_chat_delay_duration" json:"non_moderator_chat_delay_duration"`

	// Shield mode
	ShieldModeIsActive       bool   `db:"shieldmode_is_active" json:"-"`
	ShieldModeModeratorID    string `db:"shieldmode_moderator_id" json:"-"`
	ShieldModeModeratorLogin string `db:"shieldmode_moderator_login" json:"-"`
	ShieldModeModeratorName  string `db:"shieldmode_moderator_name" json:"-"`
	ShieldModeLastActivated  string `db:"shieldmode_last_activated" json:"-"`
}

func (q *Query) GetChatSettingsByBroadcaster(broadcaster string) (*DBResponse, error) {
	var r []ChatSettings

	err := q.DB.Select(&r, "SELECT * FROM chat_settings WHERE broadcaster_id = $1", broadcaster)
	if err != nil {
		return nil, err
	}

	dbr := DBResponse{
		Data:  r,
		Limit: q.Limit,
		Total: len(r),
	}

	// No cursor because there should only ever be one result

	return &dbr, err
}

func (q *Query) InsertChatSettings(s ChatSettings) error {
	stmt := generateInsertSQL("chat_settings", "broadcaster_id", s, true)
	_, err := q.DB.NamedExec(stmt, s)
	return err
}

func (q *Query) UpdateChatSettings(s ChatSettings) error {
	sql := generateUpdateSQL("chat_settings", []string{"broadcaster_id"}, s)
	_, err := q.DB.NamedExec(sql, s)
	return err
}
