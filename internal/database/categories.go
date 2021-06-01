// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

type Category struct {
	ID   string `db:"id" json:"id"`
	Name string `db:"category_name" json:"name"`
}

func (c CLIDatabase) GetCategory(cat Category) (Category, error) {
	var r Category
	rows, err := c.DB.NamedQuery(generateSQL("select * from categories", cat, SEP_AND), cat)
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

func (c CLIDatabase) InsertCategory(category Category, upsert bool) error {
	_, err := c.DB.NamedExec(`insert into categories values(:id, :category_name)`, category)
	return err
}

func (c CLIDatabase) SearchCategories(query string) ([]Category, error) {
	categories := []Category{}
	err := c.DB.Select(&categories, `select * from categories where category_name like '%$1%'`, query)
	if err != nil {
		return categories, err
	}
	return categories, nil
}
