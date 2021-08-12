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
const BAN_ADD = "moderation.user.ban"
const BAN_REMOVE = "moderation.user.unban"

type Moderator struct {
	UserID    string `db:"user_id" json:"user_id"`
	UserLogin string `db:"user_login" json:"user_login"`
	UserName  string `db:"user_name" json:"user_name"`
}

type ModeratorAction struct {
	ID                   string `db:"id" json:"id"`
	EventType            string `db:"event_type" json:"event_type"`
	EventTimestamp       string `db:"event_timestamp" json:"event_timestamp"`
	EventVersion         string `db:"event_version" json:"version"`
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

type BanActionEvent struct {
	BroadcasterID      string  `db:"broadcaster_id" json:"broadcaster_id"`
	BroadcasterLogin   string  `db:"broadcaster_login" json:"broadcaster_login"`
	BroadcasterName    string  `db:"broadcaster_name" json:"broadcaster_name"`
	UserID             string  `db:"user_id" json:"user_id"`
	UserLogin          string  `db:"user_login" json:"user_login"`
	UserName           string  `db:"user_name" json:"user_name"`
	ExpiresAt          *string `db:"expires_at" json:"expires_at"`
	Reason             string  `json:"reason"`
	ModeratorID        string  `json:"moderator_id"`
	ModeratorUserLogin string  `json:"moderator_login"`
	ModeratorUserName  string  `json:"moderator_name"`
}

type BanEvent struct {
	ID             string `db:"id" json:"id"`
	EventType      string `db:"event_type" json:"event_type"`
	EventTimestamp string `db:"event_timestamp" json:"event_timestamp"`
	EventVersion   string `db:"event_version" json:"version"`
	BanActionEvent `json:"event_data"`
}
type Ban struct {
	UserID             string  `db:"user_id" json:"user_id"`
	UserLogin          string  `db:"user_login" json:"user_login"`
	UserName           string  `db:"user_name" json:"user_name"`
	ExpiresAt          *string `db:"expires_at" json:"expires_at"`
	Reason             string  `json:"reason"`
	ModeratorID        string  `json:"moderator_id"`
	ModeratorUserLogin string  `json:"moderator_login"`
	ModeratorUserName  string  `json:"moderator_name"`
}

var es = ""

func (q *Query) GetModerationActionsByBroadcaster(broadcaster string) (*DBResponse, error) {
	var r []ModeratorAction

	err := q.DB.Select(&r, "SELECT u1.id as user_id, u1.user_login as user_login, u1.display_name as user_name, u2.id as broadcaster_id, u2.user_login as broadcaster_login, u2.display_name as broadcaster_name, ma.event_type, ma.event_version, ma.event_timestamp, ma.id FROM moderator_actions as ma JOIN users u1 ON ma.user_id = u1.id JOIN users u2 ON ma.broadcaster_id = u2.id where broadcaster_id = $1 ORDER BY ma.event_timestamp DESC", broadcaster)
	if err != nil {
		return nil, err
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

func (q *Query) AddModerator(p UserRequestParams) error {
	stmt := generateInsertSQL("moderators", "id", p, false)
	p.CreatedAt = util.GetTimestamp().UTC().Format(time.RFC3339)

	ma := ModeratorAction{
		ID:             util.RandomGUID(),
		EventType:      MOD_ADD,
		EventVersion:   "1.0",
		EventTimestamp: util.GetTimestamp().Format(time.RFC3339),
		ModeratorActionEvent: ModeratorActionEvent{
			UserID:        p.UserID,
			BroadcasterID: p.BroadcasterID,
		},
	}

	tx := q.DB.MustBegin()
	tx.NamedExec(stmt, p)
	tx.NamedExec(`INSERT INTO moderator_actions VALUES(:id, :event_timestamp, :event_type, :event_version, :broadcaster_id, :user_id)`, ma)
	return tx.Commit()
}

func (q *Query) GetModeratorsForBroadcaster(broadcasterID string, userID string) (*DBResponse, error) {
	var r []Moderator

	err := q.DB.Select(&r, "SELECT u1.id as user_id, u1.user_login as user_login, u1.display_name as user_name FROM moderators as m JOIN users u1 ON m.user_id = u1.id where broadcaster_id = $1 ORDER BY m.created_at DESC", broadcasterID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
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

func (q *Query) RemoveModerator(broadcaster string, user string) error {
	ma := ModeratorAction{
		ID:             util.RandomGUID(),
		EventType:      BAN_ADD,
		EventVersion:   "1.0",
		EventTimestamp: util.GetTimestamp().Format(time.RFC3339),
		ModeratorActionEvent: ModeratorActionEvent{
			UserID:        user,
			BroadcasterID: broadcaster,
		},
	}

	tx := q.DB.MustBegin()
	tx.Exec(`delete from moderators where broadcaster_id=$1 and user_id=$2`, broadcaster, user)
	tx.NamedExec(`INSERT INTO moderator_actions VALUES(:id, :event_timestamp, :event_type, :event_version, :broadcaster_id, :user_id)`, ma)
	return tx.Commit()
}

func (q *Query) InsertBan(p UserRequestParams) error {
	stmt := generateInsertSQL("bans", "id", p, false)
	p.CreatedAt = util.GetTimestamp().UTC().Format(time.RFC3339)

	ma := BanEvent{
		ID:             util.RandomGUID(),
		EventType:      BAN_ADD,
		EventVersion:   "1.0",
		EventTimestamp: util.GetTimestamp().Format(time.RFC3339),
		BanActionEvent: BanActionEvent{
			UserID:        p.UserID,
			BroadcasterID: p.BroadcasterID,
			ExpiresAt:     &es,
		},
	}

	tx := q.DB.MustBegin()
	tx.NamedExec(stmt, p)
	tx.NamedExec(`INSERT INTO ban_events VALUES(:id, :event_timestamp, :event_type, :event_version, :broadcaster_id, :user_id, :expires_at)`, ma)
	return tx.Commit()
}

func (q *Query) GetBans(p UserRequestParams) (*DBResponse, error) {
	r := []Ban{}
	stmt := generateSQL("select b.user_id, b.expires_at, u1.display_name as user_name, u1.user_login from bans b join users u1 on b.user_id=u1.id", p, SEP_AND)
	rows, err := q.DB.NamedQuery(stmt+q.SQL, p)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var b Ban
		err := rows.StructScan(&b)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		if b.ExpiresAt == nil {
			b.ExpiresAt = &es
		}
		b.Reason = "CLI ban"
		r = append(r, b)
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

	return &dbr, nil
}

func (q *Query) GetBanEvents(p UserRequestParams) (*DBResponse, error) {
	r := []BanEvent{}
	stmt := generateSQL("SELECT u1.id as user_id, u1.user_login as user_login, u1.display_name as user_name, u2.id as broadcaster_id, u2.user_login as broadcaster_login, u2.display_name as broadcaster_name, be.event_type, be.event_version, be.event_timestamp, be.id FROM ban_events as be JOIN users u1 ON be.user_id = u1.id JOIN users u2 ON be.broadcaster_id = u2.id", p, SEP_AND)
	rows, err := q.DB.NamedQuery(stmt+" order by be.event_timestamp desc"+q.SQL, p)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var b BanEvent
		err := rows.StructScan(&b)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		es := ""

		if b.ExpiresAt == nil {
			b.ExpiresAt = &es
		}
		// shim for https://github.com/twitchdev/twitch-cli/issues/83
		_, err = time.Parse(time.RFC3339, b.EventTimestamp)
		if err != nil {
			ts := b.EventType
			b.EventType = b.EventTimestamp
			b.EventTimestamp = ts
		}

		b.Reason = "CLI ban"
		r = append(r, b)
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

	return &dbr, nil
}

func (q *Query) GetModeratorEvents(p UserRequestParams) (*DBResponse, error) {
	r := []ModeratorAction{}
	stmt := generateSQL("SELECT u1.id as user_id, u1.user_login as user_login, u1.display_name as user_name, u2.id as broadcaster_id, u2.user_login as broadcaster_login, u2.display_name as broadcaster_name, ma.event_type, ma.event_version, ma.event_timestamp, ma.id FROM moderator_actions as ma JOIN users u1 ON ma.user_id = u1.id JOIN users u2 ON ma.broadcaster_id = u2.id", p, SEP_AND)
	rows, err := q.DB.NamedQuery(stmt+" order by ma.event_timestamp desc"+q.SQL, p)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var ma ModeratorAction
		err := rows.StructScan(&ma)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		// shim for https://github.com/twitchdev/twitch-cli/issues/83
		_, err = time.Parse(time.RFC3339, ma.EventTimestamp)
		if err != nil {
			ts := ma.EventType
			ma.EventType = ma.EventTimestamp
			ma.EventTimestamp = ts
		}

		r = append(r, ma)
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

	return &dbr, nil
}

func (q *Query) GetModerators(p UserRequestParams) (*DBResponse, error) {
	r := []Moderator{}
	stmt := generateSQL("SELECT u1.id as user_id, u1.user_login as user_login, u1.display_name as user_name FROM moderators as m JOIN users u1 ON m.user_id = u1.id", p, SEP_AND)
	rows, err := q.DB.NamedQuery(stmt+" ORDER BY m.created_at DESC "+q.SQL, p)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var m Moderator
		err := rows.StructScan(&m)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		r = append(r, m)
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

	return &dbr, nil
}
