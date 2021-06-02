// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/twitchdev/twitch-cli/internal/util"
)

type AuthenticationClient struct {
	ID          string `db:"id"`
	Secret      string `db:"secret"`
	Name        string `db:"name"`
	IsExtension bool   `db:"is_extension"`
}

type Authorization struct {
	ID        string         `db:"id"`
	ClientID  string         `db:"client_id"`
	UserID    sql.NullString `db:"user_id"`
	Token     string         `db:"token"`
	ExpiresAt string         `db:"expires_at"`
	Scopes    sql.NullString `db:"scopes"`
}

func (q *Query) GetAuthorizationByToken(token string) (Authorization, error) {
	var r Authorization
	db := q.DB

	err := db.Get(&r, "select * from authorizations where token = $1", token)
	if errors.Is(err, sql.ErrNoRows) {
		return r, nil
	} else if err != nil {
		return r, err
	}

	return r, err
}

func (q *Query) InsertOrUpdateAuthenticationClient(client AuthenticationClient, upsert bool) (AuthenticationClient, error) {
	db := q.DB

	stmt := `insert into clients values(:id, :secret, :is_extension, :name)`
	if upsert == true {
		stmt += ` on conflict(id) do update set secret=:secret, is_extension:is_extension, name=:name`
	}

	client.Secret = generateString(30)

	for {
		tx := db.MustBegin()
		tx.NamedExec(stmt, client)
		err := tx.Commit()
		if err == nil {
			return client, err
		}

		client.ID = util.RandomClientID()
	}
}

func (q *Query) CreateAuthorization(a Authorization) (Authorization, error) {
	db := q.DB

	a.Token = generateString(15)
	a.ExpiresAt = util.GetTimestamp().Add(24 * 30 * time.Hour).Format(time.RFC3339Nano)

	for {
		// loop to create unique tokens; likely won't happen, but is worth handling regardless
		tx := db.MustBegin()
		tx.NamedExec(`insert into authorizations (client_id, user_id, token, expires_at) values (:client_id, :user_id, :token, :expires_at)`, a)
		err := tx.Commit()
		if err == nil {
			return a, nil
		}
		a.Token = generateString(15)
	}
}

func generateString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}
