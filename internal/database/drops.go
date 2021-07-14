// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import "log"

type DropsEntitlement struct {
	ID        string `db:"id" json:"id" dbs:"de.id"`
	UserID    string `db:"user_id" json:"user_id"`
	BenefitID string `db:"benefit_id" json:"benefit_id"`
	GameID    string `db:"game_id" json:"game_id"`
	Timestamp string `db:"timestamp" json:"timestamp"`
	Status    string `db:"status" json:"fulfillment_status"`
}

func (q *Query) GetDropsEntitlements(de DropsEntitlement) (*DBResponse, error) {
	var r []DropsEntitlement
	stmt := generateSQL("select * from drops_entitlements de", de, SEP_AND)
	stmt += " order by timestamp desc " + q.SQL
	rows, err := q.DB.NamedQuery(stmt, de)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	for rows.Next() {
		var de DropsEntitlement
		err := rows.StructScan(&de)
		if err != nil {
			log.Print(err)
			return nil, err
		}

		r = append(r, de)
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
func (q *Query) InsertDropsEntitlement(d DropsEntitlement) error {
	stmt := generateInsertSQL("drops_entitlements", "id", d, false)
	_, err := q.DB.NamedExec(stmt, d)
	return err
}
