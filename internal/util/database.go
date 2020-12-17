// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package util

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

var dbFileName = "eventCache.db"

// EventCacheParameters is used to define required parameters when writing into the database
type EventCacheParameters struct {
	ID        string
	Event     string
	JSON      string
	FromUser  string
	ToUser    string
	Transport string
	Timestamp string
}

// EventCacheResponse is used to define the response coming from a SELECT-based function
type EventCacheResponse struct {
	ID        string
	Event     string
	JSON      string
	Transport string
	Timestamp string
}

type migrateMap struct {
	SQL     string
	Message string
}

var migrateSQL = map[int]migrateMap{
	0: {
		SQL:     "",
		Message: "",
	},
	1: {
		SQL:     `drop table events; create table events (id text not null primary key, event text not null, json text not null, from_user text not null, to_user text not null, transport text not null, timestamp text not null);`,
		Message: "Previously executed events are incompatible with new versions of the CLI.",
	},
}

const currentVersion = 1

func getDatabase() (sql.DB, error) {
	home, err := GetApplicationDir()
	if err != nil {
		return sql.DB{}, err
	}

	var path = filepath.Join(home, dbFileName)
	var needToInit = false
	if _, err := os.Stat(path); os.IsNotExist(err) {
		needToInit = true
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return sql.DB{}, err
	}

	if needToInit == true {
		initDatabase(*db)
	}

	checkAndUpdate(*db)
	return *db, nil
}

func initDatabase(db sql.DB) error {
	createSQL := `create table events (id text not null primary key, event text not null, json text not null, from_user text not null, to_user text not null, transport text not null, timestamp text not null);`

	_, err := db.Exec(createSQL)
	if err != nil {
		return err
	}

	_, err = db.Exec("PRAGMA user_version=" + strconv.Itoa(currentVersion))
	if err != nil {
		return err
	}

	return nil
}

// InsertIntoDB inserts an event into the database for replay functions later.
func InsertIntoDB(p EventCacheParameters) error {
	db, err := getDatabase()
	if err != nil {
		return err
	}

	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`insert into events(id, event, json, from_user, to_user, transport, timestamp) values(?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.ID, p.Event, p.JSON, p.FromUser, p.ToUser, p.Transport, p.Timestamp)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

// GetEventByID returns an event based on an ID provided for replay.
func GetEventByID(id string) (EventCacheResponse, error) {
	db, err := getDatabase()
	var r EventCacheResponse

	if err != nil {
		return r, err
	}

	defer db.Close()

	stmt, err := db.Prepare("select id, json, transport, event from events where id = ?")
	if err != nil {
		return r, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&r.ID, &r.JSON, &r.Transport, &r.Event)
	if err != nil {
		return r, err
	}

	return r, err
}

func checkAndUpdate(db sql.DB) error {
	var v int
	rows, err := db.Query(`PRAGMA user_version`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&v)
		if err != nil {
			return err
		}
	}

	v++
	for i := v; i < len(migrateSQL); i++ {
		_, err = db.Exec(migrateSQL[i].SQL)
		if err != nil {
			return err
		}

		fmt.Println(migrateSQL[i].Message)
	}

	_, err = db.Exec("PRAGMA user_version=" + strconv.Itoa(currentVersion))
	if err != nil {
		return err
	}
	return nil
}
