// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

type Poll struct {
	ID                         string        `db:"id" dbs:"p.id" json:"id"`
	BroadcasterID              string        `db:"broadcaster_id" dbs:"u1.id" json:"broadcaster_id"`
	BroadcasterLogin           string        `db:"broadcaster_login" dbi:"false" json:"broadcaster_login"`
	BroadcasterName            string        `db:"broadcaster_name" dbi:"false" json:"broadcaster_name"`
	Title                      string        `db:"title" dbs:"p.title" json:"title"`
	BitsVotingEnabled          bool          `db:"bits_voting_enabled" json:"bits_voting_enabled"`
	BitsPerVote                int           `db:"bits_per_vote" json:"bits_per_vote"`
	ChannelPointsVotingEnabled bool          `db:"channel_points_voting_enabled" json:"channel_points_voting_enabled"`
	ChannelPointsPerVote       int           `db:"channel_points_per_vote" json:"channel_points_per_vote"`
	Status                     string        `db:"status" json:"status"`
	Duration                   int           `db:"duration" json:"duration"`
	StartedAt                  string        `db:"started_at" json:"started_at"`
	EndedAt                    string        `db:"ended_at" json:"ended_at,omitempty"`
	Choices                    []PollsChoice `json:"choices"`
}

type PollsChoice struct {
	ID                 string `db:"id" json:"id"`
	Title              string `db:"title" json:"title"`
	Votes              int    `db:"votes" json:"votes"`
	ChannelPointsVotes int    `db:"channel_points_votes" json:"channel_points_votes"`
	BitsVotes          int    `db:"bits_votes" json:"bits_votes"`
	PollID             string `db:"poll_id" json:"-"`
}

func (q *Query) GetPolls(p Poll) (*DBResponse, error) {
	r := []Poll{}

	sql := generateSQL("select p.*, u1.user_login as broadcaster_login, u1.display_name as broadcaster_name from polls p join users u1 on p.broadcaster_id = u1.id", p, SEP_AND)
	rows, err := q.DB.NamedQuery(sql, p)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var p Poll
		err := rows.StructScan(&p)
		if err != nil {
			return nil, err
		}
		r = append(r, p)
	}

	for i, p := range r {
		var pc []PollsChoice
		err = q.DB.Select(&pc, "select * from poll_choices where poll_id=$1", p.ID)
		if err != nil {
			return nil, err
		}

		r[i].Choices = pc
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

func (q *Query) InsertPoll(p Poll) error {
	tx := q.DB.MustBegin()
	tx.NamedExec(generateInsertSQL("polls", "id", p, false), p)
	for _, c := range p.Choices {
		tx.NamedExec(generateInsertSQL("poll_choices", "id", c, false), c)
	}
	return tx.Commit()
}

func (q *Query) UpdatePoll(p Poll) error {
	_, err := q.DB.NamedExec(generateUpdateSQL("polls", []string{"id"}, p), p)
	return err
}

func (q *Query) UpdatePollChoice(p PollsChoice) error {
	_, err := q.DB.Exec("update poll_choices set votes = votes + 1 where id = $1", p.ID)
	return err
}
