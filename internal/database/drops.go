// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import "log"

type DropsEntitlement struct {
}

func (c CLIDatabase) GetDropsEntitlementById(id string) (DropsEntitlement, error) {
	var r DropsEntitlement

	err := c.DB.Get(&r, "select * from drops_entitlements where id = $1", id)
	if err != nil {
		return r, err
	}
	log.Printf("%#v", r)

	return r, err
}

func (c CLIDatabase) InsertDropsEntitlement(d DropsEntitlement, upsert bool) error {
	tx := c.DB.MustBegin()
	tx.NamedExec(`insert into drops_entitlements values(:id, :values...)`, d)
	err := tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
