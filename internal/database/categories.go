// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import "log"

type Category struct {
	ID   string `db:"id"`
	Name string `db:"category_name"`
}

func (c CLIDatabase) GetCategoryByID(id string) (Category, error) {
	var r Category
	err := c.DB.Get(&r, "select * from categories where id = $1", id)
	if err != nil {
		return r, err
	}
	log.Printf("%#v", r)

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
