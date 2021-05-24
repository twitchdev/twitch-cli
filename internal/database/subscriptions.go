// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import "log"

type Subscription struct {
}

func (c CLIDatabase) GetSubscriptionById(id string) (Subscription, error) {
	var r Subscription

	err := c.DB.Get(&r, "select * from subscriptions where id = $1", id)
	if err != nil {
		return r, err
	}
	log.Printf("%#v", r)

	return r, err
}

func (c CLIDatabase) InsertSubscription(p Subscription, upsert bool) error {
	tx := c.DB.MustBegin()
	tx.NamedExec(`insert into subscriptions values(:id, :values...)`, p)
	err := tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
