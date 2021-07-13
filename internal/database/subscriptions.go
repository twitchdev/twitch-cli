// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import (
	"database/sql"
	"log"
)

type Subscription struct {
	BroadcasterID    string          `db:"broadcaster_id" json:"broadcaster_id"`
	BroadcasterLogin string          `db:"broadcaster_login" json:"broadcaster_login"`
	BroadcasterName  string          `db:"broadcaster_name" json:"broadcaster_name"`
	UserID           string          `db:"user_id" json:"user_id"`
	UserLogin        string          `db:"user_login" json:"user_login"`
	UserName         string          `db:"user_name" json:"user_name"`
	IsGift           bool            `db:"is_gift" json:"is_gift"`
	GifterID         *sql.NullString `db:"gifter_id" json:"gifter_id,omitempty"`
	GifterName       *sql.NullString `db:"gifter_name" json:"gifter_name,omitempty"`
	GifterLogin      *sql.NullString `db:"gifter_login" json:"gifter_login,omitempty"`
	Tier             string          `db:"tier" json:"tier"`
	CreatedAt        string          `db:"created_at" json:"-"`
}

type SubscriptionInsert struct {
	BroadcasterID string          `db:"broadcaster_id" json:"broadcaster_id"`
	UserID        string          `db:"user_id" json:"user_id"`
	IsGift        bool            `db:"is_gift" json:"is_gift"`
	GifterID      *sql.NullString `db:"gifter_id" json:"gifter_id,omitempty"`
	Tier          string          `db:"tier" json:"tier"`
	CreatedAt     string          `db:"created_at" json:"-"`
}

func (q *Query) GetSubscriptions(s Subscription) (*DBResponse, error) {
	r := []Subscription{}
	sql := generateSQL("SELECT u1.id as user_id, u1.user_login as user_login, u1.display_name as user_name, u2.id as broadcaster_id, u2.user_login as broadcaster_login, u2.display_name as broadcaster_name, u3.id as gifter_id, u3.user_login as gifter_login, u3.display_name as gifter_name, s.tier as tier, s.is_gift as is_gift FROM subscriptions as s JOIN users u1 ON s.user_id = u1.id JOIN users u2 ON s.broadcaster_id = u2.id LEFT JOIN users u3 ON s.gifter_id = u3.id", s, SEP_AND)
	sql += " order by s.created_at desc"
	sql += q.SQL

	rows, err := q.DB.NamedQuery(sql, s)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	for rows.Next() {
		s := Subscription{}
		err := rows.StructScan(&s)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		r = append(r, s)
	}

	var total int
	rows, err = q.DB.NamedQuery(generateSQL("select count(*) from subscriptions", s, SEP_AND), s)
	for rows.Next() {
		err := rows.Scan(&total)
		if err != nil {
			log.Print(err)
			return nil, err
		}
	}

	dbr := DBResponse{
		Data:  r,
		Limit: q.Limit,
		Total: total,
	}

	if len(r) != q.Limit {
		q.PaginationCursor = ""
	}

	dbr.Cursor = q.PaginationCursor

	return &dbr, err
}

func (q *Query) InsertSubscription(s SubscriptionInsert) error {
	stmt := generateInsertSQL("subscriptions", "", s, false)
	_, err := q.DB.NamedExec(stmt, s)
	return err
}
