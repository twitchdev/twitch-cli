// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import "log"

type Extension struct {
}

func (c CLIDatabase) GetExtensionById(id string) (Extension, error) {
	var r Extension

	err := c.DB.Get(&r, "select * from principle where id = $1", id)
	if err != nil {
		return r, err
	}
	log.Printf("%#v", r)

	return r, err
}

func (c CLIDatabase) InsertExtension(p Extension, upsert bool) error {
	tx := c.DB.MustBegin()
	tx.NamedExec(`insert into principle values(:id, :values...)`, p)
	return tx.Commit()
}
