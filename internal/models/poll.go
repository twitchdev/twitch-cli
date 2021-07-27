// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type PollEventSubResponse struct {
	Subscription EventsubSubscription `json:"subscription"`
	Event        PollEventSubEvent    `json:"event"`
}

type PollEventSubEvent struct {
	ID                   string                      `json:"id"`
	BroadcasterUserID    string                      `json:"broadcaster_user_id"`
	BroadcasterUserLogin string                      `json:"broadcaster_user_login"`
	BroadcasterUserName  string                      `json:"broadcaster_user_name"`
	Title                string                      `json:"title"`
	Choices              []PollEventSubEventChoice   `json:"choices"`
	BitsVoting           PollEventSubEventGoodVoting `json:"bits_voting"`
	ChannelPointsVoting  PollEventSubEventGoodVoting `json:"channel_points_voting"`
	Status               string                      `json:"status,omitempty"`
	StartedAt            string                      `json:"started_at"`
	EndsAt               string                      `json:"ends_at,omitempty"`
	EndedAt              string                      `json:"ended_at,omitempty"`
}

type PollEventSubEventChoice struct {
	ID                 string `json:"id"`
	Title              string `json:"title"`
	BitsVotes          *int   `json:"bits_votes,omitempty"`
	ChannelPointsVotes *int   `json:"channel_points_votes,omitempty"`
	Votes              *int   `json:"votes,omitempty"`
}

type PollEventSubEventGoodVoting struct {
	IsEnabled     bool `json:"is_enabled"`
	AmountPerVote int  `json:"amount_per_vote"`
}
