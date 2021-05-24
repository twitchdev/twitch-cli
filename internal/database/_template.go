// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import "log"

type Principle struct {
}

func (c CLIDatabase) GetPrincipleById(id string) (Principle, error) {
	var r Principle

	err := c.DB.Get(&r, "select * from principle where id = $1", id)
	if err != nil {
		return r, err
	}
	log.Printf("%#v", r)

	return r, err
}

func (c CLIDatabase) InsertPrinciple(p Principle, upsert bool) error {
	tx := c.DB.MustBegin()
	tx.NamedExec(`insert into principle values(:id, :values...)`, p)
	err := tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
