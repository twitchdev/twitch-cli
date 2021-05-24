// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/twitchdev/twitch-cli/internal/util"
)

const MOD_ADD = "moderation.moderator.add"
const MOD_REMOVE = "moderation.moderator.remove"

type Moderator struct {
	UserID    string `db:"user_id" json:"user_id"`
	UserLogin string `db:"user_login" json:"user_login"`
	UserName  string `db:"user_name" json:"user_name"`
}

type ModeratorAction struct {
	ID                   string `db:"id" json:"id"`
	EventType            string `db:"event_type" json:"event_type"`
	EventTimestamp       string `db:"event_timestamp" json:"event_timestamp"`
	EventVersion         string `db:"event_version" json:"event_version"`
	ModeratorActionEvent `json:"event_data"`
}

type ModeratorActionEvent struct {
	BroadcasterID    string `db:"broadcaster_id" json:"broadcaster_id"`
	BroadcasterLogin string `db:"broadcaster_login" json:"broadcaster_login"`
	BroadcasterName  string `db:"broadcaster_name" json:"broadcaster_name"`
	UserID           string `db:"user_id" json:"user_id"`
	UserLogin        string `db:"user_login" json:"user_login"`
	UserName         string `db:"user_name" json:"user_name"`
}

func (c CLIDatabase) GetModerationActionsByBroadcaster(broadcaster string) ([]ModeratorAction, error) {
	var r []ModeratorAction

	err := c.DB.Select(&r, "SELECT u1.id as user_id, u1.user_login as user_login, u1.display_name as user_name, u2.id as broadcaster_id, u2.user_login as broadcaster_login, u2.display_name as broadcaster_name, ma.event_type, ma.event_version, ma.event_timestamp, ma.id FROM moderator_actions as ma JOIN users u1 ON ma.user_id = u1.id JOIN users u2 ON ma.broadcaster_id = u2.id where broadcaster_id = $1 ORDER BY ma.event_timestamp DESC", broadcaster)
	if err != nil {
		return r, err
	}
	log.Printf("%#v", r)

	return r, err
}

func (c CLIDatabase) AddModerator(broadcaster string, user string) error {
	ma := ModeratorAction{
		ID:             util.RandomGUID(),
		EventType:      MOD_ADD,
		EventVersion:   "1.0",
		EventTimestamp: util.GetTimestamp().Format(time.RFC3339),
		ModeratorActionEvent: ModeratorActionEvent{
			UserID:        user,
			BroadcasterID: broadcaster,
		},
	}

	tx := c.DB.MustBegin()
	tx.Exec(`insert into moderators values($1, $2, $3)`, broadcaster, user, util.GetTimestamp().UTC().Format(time.RFC3339))

	tx.NamedExec(`INSERT INTO moderator_actions VALUES(:id, :event_type, :event_timestamp, :event_version, :broadcaster_id, :user_id)`, ma)
	return tx.Commit()
}

func (c CLIDatabase) GetModeratorsForBroadcaster(broadcasterID string, userID string) ([]Moderator, error) {
	var r []Moderator

	err := c.DB.Select(&r, "SELECT u1.id as user_id, u1.user_login as user_login, u1.display_name as user_name FROM moderators as m JOIN users u1 ON m.user_id = u1.id where broadcaster_id = $1 ORDER BY m.created_at DESC", broadcasterID)
	if errors.Is(err, sql.ErrNoRows) {
		return r, nil
	} else if err != nil {
		return r, err
	}

	return r, err
}

func (c CLIDatabase) RemoveModerator(broadcaster string, user string) error {
	ma := ModeratorAction{
		ID:             util.RandomGUID(),
		EventType:      MOD_REMOVE,
		EventVersion:   "1.0",
		EventTimestamp: util.GetTimestamp().Format(time.RFC3339),
		ModeratorActionEvent: ModeratorActionEvent{
			UserID:        user,
			BroadcasterID: broadcaster,
		},
	}

	tx := c.DB.MustBegin()
	tx.Exec(`delete from moderators where broadcaster_id=$1 and user_id=$2`, broadcaster, user)
	tx.NamedExec(`INSERT INTO moderator_actions VALUES(:id, :event_type, :event_timestamp, :event_version, :broadcaster_id, :user_id)`, ma)
	return tx.Commit()
}
