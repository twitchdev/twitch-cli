// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import (
	"fmt"
)

type Category struct {
	ID          string `db:"id" json:"id"`
	Name        string `db:"category_name" json:"name"`
	BoxartURL   string `json:"boxart_url"`
	ViewerCount int    `db:"vc" json:"-"`
}

func (q *Query) GetCategories(cat Category) (*DBResponse, error) {
	var r []Category
	rows, err := q.DB.NamedQuery(generateSQL("select * from categories", cat, SEP_AND)+q.SQL, cat)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var cat Category
		err := rows.StructScan(&cat)
		if err != nil {
			return nil, err
		}
		r = append(r, cat)
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

func (q *Query) InsertCategory(category Category, upsert bool) error {
	_, err := q.DB.NamedExec(`insert into categories values(:id, :category_name)`, category)
	return err
}

func (q *Query) SearchCategories(query string) (*DBResponse, error) {
	r := []Category{}
	err := q.DB.Select(&r, `select * from categories where lower(category_name) like lower($1) `+q.SQL, fmt.Sprintf("%%%v%%", query))
	if err != nil {
		return nil, err
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

func (q *Query) GetTopGames() (*DBResponse, error) {
	r := []Category{}

	err := q.DB.Select(&r, "select c.id, c.category_name, SUM(s.viewer_count) as vc from users u join streams s on s.broadcaster_id = u.id left join categories c on u.category_id = c.id group by c.id, c.category_name order by vc desc")
	if err != nil {
		return nil, err
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
