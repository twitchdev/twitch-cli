// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import "log"

type Team struct {
}

func (c CLIDatabase) GetTeamById(id string) (Team, error) {
	db, err := getDatabase()
	var r Team

	if err != nil {
		return r, err
	}

	defer db.Close()

	err = db.Get(&r, "select * from teams where id = $1", id)
	if err != nil {
		return r, err
	}
	log.Printf("%#v", r)

	return r, err
}

func (c CLIDatabase) InsertTeam(p Team, upsert bool) error {
	tx := c.DB.MustBegin()
	tx.NamedExec(`insert into teams values(:id, :values...)`, p)
	err := tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
