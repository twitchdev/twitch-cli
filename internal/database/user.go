// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/twitchdev/twitch-cli/internal/util"
)

type User struct {
	ID              string `db:"id"`
	UserLogin       string `db:"user_login"`
	DisplayName     string `db:"display_name"`
	Email           string `db:"email"`
	UserType        string `db:"user_type"`
	BroadcasterType string `db:"broadcaster_type"`
	UserDescription string `db:"user_description"`
	CreatedAt       string `db:"created_at"`
	ModifiedAt      string `db:"modified_at"`
}

type Follow struct {
	BroadcasterID    string `db:"to_id" json:"to_id"`
	BroadcasterLogin string `db:"to_login" json:"to_login"`
	BroadcasterName  string `db:"to_name" json:"to_name"`
	ViewerID         string `db:"from_id" json:"from_id"`
	ViewerLogin      string `db:"from_login" json:"from_login"`
	ViewerName       string `db:"from_name" json:"from_name"`
	FollowedAt       string `db:"created_at" json:"followed_at"`
}

type Block struct {
	UserID    string `db:"user_id" son:"user_id"`
	UserLogin string `db:"user_login" json:"user_login"`
	UserName  string `db:"user_name" json:"display_name"`
}

func (c CLIDatabase) GetUserByID(id string) (User, error) {
	db := c.DB
	var r User

	err := db.Get(&r, "select * from users where id = $1", id)
	if errors.Is(err, sql.ErrNoRows) {
		return r, nil
	} else if err != nil {
		return r, err
	}

	return r, err
}

func (c CLIDatabase) GetUserByLogin(login string) (User, error) {
	db := c.DB
	var r User

	err := db.Get(&r, "select * from users where user_login = $1", login)
	if errors.Is(err, sql.ErrNoRows) {
		return r, nil
	} else if err != nil {
		return r, err
	}

	return r, err
}

func (c CLIDatabase) InsertUser(u User, upsert bool) error {
	db := c.DB

	stmt := `insert into users values(:id, :user_login, :display_name, :email, :user_type, :broadcaster_type, :user_description, :created_at, :modified_at)`
	if upsert == true {
		stmt += ` on conflict(id) do update set user_login=:user_login, display_name=:display_name, email=:email, user_type=:user_type, broadcaster_type=:broadcaster_type, user_description=:user_description, modified_at=:modified_at`
	}
	_, err := db.NamedExec(stmt, u)
	return err
}

func (c CLIDatabase) AddFollow(from_user string, to_user string) error {
	_, err := c.DB.Exec(`insert into follows values($1, $2, $3)`, to_user, from_user, util.GetTimestamp().UTC().Format(time.RFC3339))
	return err
}

func (c CLIDatabase) GetFollowsByBroadcaster(id string) ([]Follow, error) {
	db := c.DB
	var r []Follow

	err := db.Select(&r, "SELECT u1.id as to_id, u1.user_login as to_login, u1.display_name as to_name, u2.id as from_id, u2.user_login as from_login, u2.display_name as from_name, f.created_at as created_at FROM follows as f JOIN users u1 ON f.broadcaster_id = u1.id JOIN users u2 ON f.user_id = u2.id where broadcaster_id = $1 ORDER BY f.created_at DESC", id)
	if errors.Is(err, sql.ErrNoRows) {
		return r, nil
	} else if err != nil {
		return r, err
	}

	return r, err
}

func (c CLIDatabase) GetFollowsByViewer(id string) ([]Follow, error) {
	db := c.DB
	var r []Follow

	err := db.Select(&r, "SELECT u1.id as to_id, u1.user_login as to_login, u1.display_name as to_name, u2.id as from_id, u2.user_login as from_login, u2.display_name as from_name, f.created_at as created_at FROM follows as f JOIN users u1 ON f.broadcaster_id = u1.id JOIN users u2 ON f.user_id = u2.id where user_id = $1 ORDER BY f.created_at DESC", id)
	if errors.Is(err, sql.ErrNoRows) {
		return r, nil
	} else if err != nil {
		return r, err
	}

	return r, err
}

func (c CLIDatabase) GetFollowsByBroadcasterAndUser(b string, u string) ([]Follow, error) {
	db := c.DB
	var r []Follow

	err := db.Select(&r, "SELECT u1.id as to_id, u1.user_login as to_login, u1.display_name as to_name, u2.id as from_id, u2.user_login as from_login, u2.display_name as from_name, f.created_at as created_at FROM follows as f JOIN users u1 ON f.broadcaster_id = u1.id JOIN users u2 ON f.user_id = u2.id where broadcaster_id = $1 and user_id = $2 ORDER BY f.created_at DESC", b, u)
	if errors.Is(err, sql.ErrNoRows) {
		return r, nil
	} else if err != nil {
		return r, err
	}

	return r, err
}

func (c CLIDatabase) DeleteFollow(from_user string, to_user string) error {
	_, err := c.DB.Exec(`delete from follows where broadcaster_id=$1 and user_id=$2`, to_user, from_user)
	return err
}

func (c CLIDatabase) AddBlock(from_user string, to_user string) error {
	_, err := c.DB.Exec(`insert into blocks values($1, $2, $3)`, to_user, from_user, util.GetTimestamp().UTC().Format(time.RFC3339))
	return err
}

func (c CLIDatabase) GetBlocksByBroadcaster(id string) ([]Block, error) {
	var r []Block

	err := c.DB.Select(&r, "SELECT u1.id as user_id, u1.user_login as user_login, u1.display_name as user_name FROM blocks as b JOIN users u1 ON f.user_id = u1.id where broadcaster_id = $1 ORDER BY b.created_at DESC", id)
	if errors.Is(err, sql.ErrNoRows) {
		return r, nil
	} else if err != nil {
		return r, err
	}

	return r, err
}

func (c CLIDatabase) DeleteBlock(from_user string, to_user string) error {
	_, err := c.DB.Exec(`delete from blocks where broadcaster_id=$1 and user_id=$2`, to_user, from_user)
	return err
}
