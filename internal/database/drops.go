// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import "log"

type DropsEntitlement struct {
}

func (q *Query) GetDropsEntitlementById(id string) (*DBResposne, error) {
	var r []DropsEntitlement

	err := q.DB.Get(&r, "select * from drops_entitlements where id = $1", id)
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

func (q *Query) InsertDropsEntitlement(d DropsEntitlement, upsert bool) error {
	tx := q.DB.MustBegin()
	tx.NamedExec(`insert into drops_entitlements values(:id, :values...)`, d)
	err := tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
