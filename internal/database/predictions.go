// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

type Prediction struct {
	ID               string              `db:"id" dbs:"p.id" json:"id"`
	BroadcasterID    string              `db:"broadcaster_id" dbs:"u1.id" json:"broadcaster_id"`
	BroadcasterLogin string              `db:"broadcaster_login" dbi:"false" json:"broadcaster_login"`
	BroadcasterName  string              `db:"broadcaster_name" dbi:"false" json:"broadcaster_name"`
	Title            string              `db:"title" json:"title"`
	WinningOutcomeID *string             `db:"winning_outcome_id" json:"winning_outcome_id"`
	PredictionWindow int                 `db:"prediction_window" json:"prediction_window"`
	Status           string              `db:"status" json:"status"`
	StartedAt        string              `db:"created_at" json:"created_at"`
	EndedAt          *string             `db:"ended_at" json:"ended_at"`
	LockedAt         *string             `db:"locked_at" json:"locked_at"`
	Outcomes         []PredictionOutcome `json:"outcomes"`
}

type PredictionOutcome struct {
	ID            string                  `db:"id" dbs:"po.id" json:"id"`
	Title         string                  `db:"title" json:"title"`
	Users         int                     `db:"users" json:"users"`
	ChannelPoints int                     `db:"channel_points" json:"channel_points"`
	TopPredictors []*PredictionPrediction `json:"top_predictors"`
	Color         string                  `db:"color" json:"color"`
	PredictionID  string                  `db:"prediction_id" json:"-"`
}
type PredictionPrediction struct {
	PredictionID string `db:"prediction_id" json:"-"`
	UserID       string `db:"user_id" json:"user_id"`
	UserLogin    string `db:"user_login" dbi:"false" json:"user_login"`
	UserName     string `db:"user_name" dbi:"false" json:"user_name"`
	Amount       int    `db:"amount" json:"channel_points_used"`
	OutcomeID    string `db:"outcome_id" json:"-"`
	// calculated fields
	AmountWon int `json:"channel_points_won"`
}

func (q *Query) GetPredictions(p Prediction) (*DBResponse, error) {
	r := []Prediction{}

	sql := generateSQL("select p.*, u1.user_login as broadcaster_login, u1.display_name as broadcaster_name from predictions p join users u1 on p.broadcaster_id = u1.id", p, SEP_AND)
	rows, err := q.DB.NamedQuery(sql, p)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		p := Prediction{}
		err := rows.StructScan(&p)
		if err != nil {
			return nil, err
		}

		r = append(r, p)
	}

	for i, p := range r {
		outcomes := []PredictionOutcome{}
		err = q.DB.Select(&outcomes, "select po.id, po.title, po.color, COUNT(pp.prediction_id) as users, IFNULL(SUM(pp.amount),0) as channel_points from prediction_outcomes po left join prediction_predictions pp on po.id = pp.outcome_id where po.prediction_id = $1 group by po.id, po.title, po.color", p.ID)
		if err != nil {
			return nil, err
		}

		for i, o := range outcomes {
			if p.WinningOutcomeID != nil {
				topPredictors := []PredictionPrediction{}
				err = q.DB.Select(&topPredictors, "select pp.*, u1.user_login, u1.display_name as user_name from prediction_predictions pp join users u1 on pp.user_id = u1.id where pp.outcome_id = $1 order by pp.amount desc limit 2", o.ID)
				if err != nil {
					return nil, err
				}
				tp := []*PredictionPrediction{}
				for _, p := range topPredictors {
					tp = append(tp, &p)
				}
				outcomes[i].TopPredictors = tp

			}
			if len(outcomes[i].TopPredictors) == 0 {
				outcomes[i].TopPredictors = nil
			}
		}
		r[i].Outcomes = outcomes
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

func (q *Query) InsertPrediction(p Prediction) error {
	tx := q.DB.MustBegin()
	tx.NamedExec(generateInsertSQL("predictions", "id", p, false), p)

	for _, o := range p.Outcomes {
		tx.NamedExec(generateInsertSQL("prediction_outcomes", "id", o, false), o)
	}
	return tx.Commit()
}

func (q *Query) InsertPredictionPrediction(p PredictionPrediction) error {
	_, err := q.DB.NamedExec(generateInsertSQL("prediction_predictions", "user_id", p, false), p)
	return err
}

func (q *Query) UpdatePrediction(p Prediction) error {
	_, err := q.DB.NamedExec(generateUpdateSQL("predictions", []string{"id", "broadcaster_id"}, p), p)
	return err
}
