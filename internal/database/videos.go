// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import "log"

type Video struct {
}

func (c CLIDatabase) GetVideoById(id string) (Video, error) {
	var r Video

	err := c.DB.Get(&r, "select * from videos where id = $1", id)
	if err != nil {
		return r, err
	}
	log.Printf("%#v", r)

	return r, err
}

func (c CLIDatabase) InsertVideo(v Video, upsert bool) error {
	tx := c.DB.MustBegin()
	tx.NamedExec(`insert into videos values(:id, :values...)`, v)
	err := tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
