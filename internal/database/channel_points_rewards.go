// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import "log"

type ChannelPointsReward struct {
}

func (q *Query) GetChannelPointsRewardById(id string) (ChannelPointsReward, error) {
	var r ChannelPointsReward

	err := q.DB.Get(&r, "select * from channel_points_rewards where id = $1", id)
	if err != nil {
		return r, err
	}
	log.Printf("%#v", r)

	return r, err
}

func (q *Query) InsertChannelPointsReward(r ChannelPointsReward, upsert bool) error {
	tx := q.DB.MustBegin()
	tx.NamedExec(`insert into channel_points_rewards values(:id, :values...)`, r)
	err := tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
