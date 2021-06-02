// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

type Category struct {
	ID   string `db:"id" json:"id"`
	Name string `db:"category_name" json:"name"`
}

func (q *Query) GetCategories(cat Category) (*DBResposne, error) {
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

func (q *Query) InsertCategory(category Category, upsert bool) error {
	_, err := q.DB.NamedExec(`insert into categories values(:id, :category_name)`, category)
	return err
}

func (q *Query) SearchCategories(query string) ([]Category, error) {
	categories := []Category{}
	err := q.DB.Select(&categories, `select * from categories where category_name like '%$1%'`, query)
	if err != nil {
		return categories, err
	}
	return categories, nil
}
