// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import "log"

type Principle struct {
}

func (c CLIDatabase) GetPrinciple(p Principle) (Principle, error) {
	var r Principle

	sql := generateSQL("select * from principle", u, SEP_AND)
	sql = fmt.Sprintf("%v LIMIT 1", sql)
	rows, err := c.DB.NamedQuery(sql, u)
	if err != nil {
		return r, err
	}

	for rows.Next() {
		err := rows.StructScan(&r)
		if err != nil {
			return r, err
		}
	}

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
