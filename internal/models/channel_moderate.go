// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type ChannelModerateFollowersAction struct {
	FollowDurationMinutes int `json:"follow_duration_minutes"`
}

type ChannelModerateSlowAction struct {
	WaitTimeSeconds int `json:"wait_time_seconds"`
}

type ChannelModerateVipAction struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
}

type ChannelModerateUnvipAction struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
}

type ChannelModerateModAction struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
}

type ChannelModerateUnmodAction struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
}

type ChannelModerateBanAction struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
	Reason    string `json:"reason"`
}

type ChannelModerateUnbanAction struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
}

type ChannelModerateTimeoutAction struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
	Reason    string `json:"reason"`
	// TODO: This should be a timestamp type or something
	ExpiresAt string `json:"expires_at"`
}

type ChannelModerateUntimeoutAction struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
}

type ChannelModerateRaidAction struct {
	UserID      string `json:"user_id"`
	UserLogin   string `json:"user_login"`
	UserName    string `json:"user_name"`
	ViewerCount int    `json:"viewer_count"`
}

type ChannelModerateUnraidAction struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
}

type ChannelModerateDeleteAction struct {
	UserID      string `json:"user_id"`
	UserLogin   string `json:"user_login"`
	UserName    string `json:"user_name"`
	MessageID   string `json:"message_id"`
	MessageBody string `json:"message_body"`
}

type ChannelModerateAutomodTermsAction struct {
	Action      string   `json:"action"`
	List        string   `json:"list"`
	Terms       []string `json:"terms"`
	FromAutomod bool     `json:"from_automod"`
}

type ChannelModerateUnbanRequestAction struct {
	IsApproved       bool   `json:"is_approved"`
	UserID           string `json:"user_id"`
	UserLogin        string `json:"user_login"`
	UserName         string `json:"user_name"`
	ModeratorMessage string `json:"moderator_message"`
}

type ChannelModerateWarnAction struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
	Reason    string `json:"reason"`

	ChatRulesCited []string `json:"chat_rules_cited"`
}

type ChannelModerateEventSubEvent struct {
	BroadcasterUserID          string                             `json:"broadcaster_user_id"`
	BroadcasterUserLogin       string                             `json:"broadcaster_user_login"`
	BroadcasterUserName        string                             `json:"broadcaster_user_name"`
	SourceBroadcasterUserID    *string                            `json:"source_broadcaster_user_id"`
	SourceBroadcasterUserLogin *string                            `json:"source_broadcaster_user_login"`
	SourceBroadcasterUserName  *string                            `json:"source_broadcaster_user_name"`
	ModeratorUserID            string                             `json:"moderator_user_id"`
	ModeratorUserLogin         string                             `json:"moderator_user_login"`
	ModeratorUserName          string                             `json:"moderator_user_name"`
	Action                     string                             `json:"action"`
	Followers                  *ChannelModerateFollowersAction    `json:"followers"`
	Slow                       *ChannelModerateSlowAction         `json:"slow"`
	Vip                        *ChannelModerateVipAction          `json:"vip"`
	Unvip                      *ChannelModerateUnvipAction        `json:"unvip"`
	Mod                        *ChannelModerateModAction          `json:"mod"`
	Unmod                      *ChannelModerateUnmodAction        `json:"unmod"`
	Ban                        *ChannelModerateBanAction          `json:"ban"`
	Unban                      *ChannelModerateUnbanAction        `json:"unban"`
	Timeout                    *ChannelModerateTimeoutAction      `json:"timeout"`
	Untimeout                  *ChannelModerateUntimeoutAction    `json:"untimeout"`
	Raid                       *ChannelModerateRaidAction         `json:"raid"`
	Unraid                     *ChannelModerateUnraidAction       `json:"unraid"`
	Delete                     *ChannelModerateDeleteAction       `json:"delete"`
	AutomodTerms               *ChannelModerateAutomodTermsAction `json:"automod_terms"`
	UnbanRequest               *ChannelModerateUnbanRequestAction `json:"unban_request"`
	Warn                       *ChannelModerateWarnAction         `json:"warn"`
	SharedChatBan              *ChannelModerateBanAction          `json:"shared_chat_ban"`
	SharedChatUnban            *ChannelModerateUnbanAction        `json:"shared_chat_unban"`
	SharedChatTimeout          *ChannelModerateTimeoutAction      `json:"shared_chat_timeout"`
	SharedChatUntimeout        *ChannelModerateUntimeoutAction    `json:"shared_chat_untimeout"`
	SharedChatDelete           *ChannelModerateDeleteAction       `json:"shared_chat_delete"`
}

type ChannelModerateEventSubResponse struct {
	Subscription EventsubSubscription         `json:"subscription"`
	Event        ChannelModerateEventSubEvent `json:"event"`
}
