// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import "log"

type Extension struct {
}

func (q *Query) GetExtensionById(id string) (*DBResposne, error) {
	var r []Extension

	err := q.DB.Get(&r, "select * from principle where id = $1", id)
	if err != nil {
		return nil, err
	}
	log.Printf("%#v", r)

	dbr := DBResposne{
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

func (q *Query) InsertExtension(p Extension, upsert bool) error {
	tx := q.DB.MustBegin()
	tx.NamedExec(`insert into principle values(:id, :values...)`, p)
	return tx.Commit()
}
