// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import "log"

type ChannelPointsRedemption struct {
}

func (q *Query) GetChannelPointsRedemptionById(id string) (ChannelPointsRedemption, error) {
	db := q.DB
	var r ChannelPointsRedemption

	err := db.Get(&r, "select * from channel_points_redemptions where id = $1", id)
	if err != nil {
		return r, err
	}
	log.Printf("%#v", r)

	return r, err
}

func (q *Query) InsertChannelPointsRedemption(r ChannelPointsRedemption, upsert bool) error {
	db := q.DB

	tx := db.MustBegin()
	tx.NamedExec(`insert into channel_points_redemptions values(:id, :values...)`, r)
	err := tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
