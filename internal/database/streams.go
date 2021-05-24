// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import "log"

type Stream struct {
}

func (c CLIDatabase) GetStreamById(id string) (Stream, error) {
	var r Stream

	err := c.DB.Get(&r, "select * from streams where id = $1", id)
	if err != nil {
		return r, err
	}
	log.Printf("%#v", r)

	return r, err
}

func (c CLIDatabase) InsertStream(p Stream, upsert bool) error {
	tx := c.DB.MustBegin()
	tx.NamedExec(`insert into streams values(:id, :values...)`, p)
	err := tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
