// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

// EventCacheParameters is used to define required parameters when writing into the database
type EventCacheParameters struct {
	ID        string `db:"id"`
	Event     string `db:"event"`
	JSON      string `db:"json"`
	FromUser  string `db:"from_user"`
	ToUser    string `db:"to_user"`
	Transport string `db:"transport"`
	Timestamp string `db:"timestamp"`
}

// EventCacheResponse is used to define the response coming from a SELECT-based function
type EventCacheResponse struct {
	ID        string
	Event     string
	JSON      string
	Transport string
	Timestamp string
}

// InsertIntoDB inserts an event into the database for replay functions later.
func (q *Query) InsertIntoDB(p EventCacheParameters) error {
	db := q.DB

	tx := db.MustBegin()
	tx.NamedExec(`insert into events(id, event, json, from_user, to_user, transport, timestamp) values(:id, :event, :json, :from_user, :to_user, :transport, :timestamp)`, p)
	err := tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// GetEventByID returns an event based on an ID provided for replay.
func (q *Query) GetEventByID(id string) (EventCacheResponse, error) {
	db := q.DB
	var r EventCacheResponse

	err := db.Get(&r, "select id, json, transport, event from events where id = $1", id)
	if err != nil {
		return r, err
	}

	return r, err
}
