// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import "log"

type Principle struct {
}

func (q *Query) GetPrinciple(p Principle) (*DBResponse, error) {
	var r Principle

	sql := generateSQL("select * from principle", u, SEP_AND)
	sql = fmt.Sprintf("%v LIMIT 1", sql)
	rows, err := q.DB.NamedQuery(sql, u)
	if err != nil {
		return r, err
	}

	for rows.Next() {
		err := rows.StructScan(&r)
		if err != nil {
			return r, err
		}
	}

	dbr := DBResponse{
		Data:  r,
		Limit: q.Limit,
		Total: len(r),
	}

	if len(r) != q.Limit {
		q.PaginationCursor = ""
	}

	dbr.Cursor = q.PaginationCursor

	return &dbr, err
}

func (q *Query) InsertPrinciple(p Principle, upsert bool) error {
	tx := q.DB.MustBegin()
	tx.NamedExec(`insert into principle values(:id, :values...)`, p)
	err := tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
