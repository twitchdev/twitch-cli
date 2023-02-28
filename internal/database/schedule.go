// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import (
	"database/sql"
	"time"
)

type Schedule struct {
	Segments  []ScheduleSegment `json:"segments"`
	UserID    string            `db:"broadcaster_id" json:"broadcaster_id"`
	UserLogin string            `db:"broadcaster_login" json:"broadcaster_login" dbi:"false"`
	UserName  string            `db:"broadcaster_name" json:"broadcaster_name" dbi:"false"`
	Vacation  *ScheduleVacation `json:"vacation"`
}

type ScheduleSegment struct {
	ID            string           `db:"id" json:"id" dbs:"s.id"`
	Title         string           `db:"title" json:"title"`
	StartTime     string           `db:"starttime" json:"start_time"`
	EndTime       string           `db:"endtime" json:"end_time"`
	IsRecurring   bool             `db:"is_recurring" json:"is_recurring"`
	IsVacation    bool             `db:"is_vacation" json:"-"`
	Category      *SegmentCategory `json:"category"`
	UserID        string           `db:"broadcaster_id" json:"-"`
	Timezone      string           `db:"timezone" json:"timezone,omitempty"`
	CategoryID    *string          `db:"category_id" json:"-"`
	CategoryName  *string          `db:"category_name" dbi:"false" json:"-"`
	IsCanceled    *bool            `db:"is_canceled" json:"-"`
	CanceledUntil *string          `json:"canceled_until"`
}
type ScheduleVacation struct {
	ID        string `db:"id" json:"-"`
	StartTime string `db:"starttime" json:"start_time"`
	EndTime   string `db:"endtime" json:"end_time"`
}

type SegmentCategory struct {
	ID           *string `db:"category_id" json:"id" dbs:"category_id"`
	CategoryName *string `db:"category_name" json:"name" dbi:"false"`
}

func (q *Query) GetSchedule(p ScheduleSegment, startTime time.Time) (*DBResponse, error) {
	r := Schedule{}

	u, err := q.GetUser(User{ID: p.UserID})
	if err != nil {
		return nil, err
	}
	r.UserID = u.ID
	r.UserLogin = u.UserLogin
	r.UserName = u.DisplayName

	sql := generateSQL("select s.*, c.category_name from stream_schedule s left join categories c on s.category_id = c.id", p, SEP_AND)
	p.StartTime = startTime.Format(time.RFC3339)
	sql += " and datetime(starttime) >= datetime(:starttime) " + q.SQL
	rows, err := q.DB.NamedQuery(sql, p)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var s ScheduleSegment
		err := rows.StructScan(&s)
		if err != nil {
			return nil, err
		}
		if s.CategoryID != nil {
			s.Category = &SegmentCategory{
				ID:           s.CategoryID,
				CategoryName: s.CategoryName,
			}
		}
		if s.IsVacation {
			r.Vacation = &ScheduleVacation{
				StartTime: s.StartTime,
				EndTime:   s.EndTime,
			}
		} else {
			r.Segments = append(r.Segments, s)
		}
	}
	v, err := q.GetVacations(ScheduleSegment{UserID: p.UserID})
	if err != nil {
		return nil, err
	}
	r.Vacation = &v
	dbr := DBResponse{
		Data:  r,
		Limit: q.Limit,
		Total: len(r.Segments),
	}

	if len(r.Segments) != q.Limit {
		q.PaginationCursor = ""
	}

	dbr.Cursor = q.PaginationCursor

	return &dbr, err
}

func (q *Query) InsertSchedule(p ScheduleSegment) error {
	tx := q.DB.MustBegin()
	_, err := tx.NamedExec(generateInsertSQL("stream_schedule", "id", p, false), p)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (q *Query) DeleteSegment(id string, broadcasterID string) error {
	_, err := q.DB.Exec("delete from stream_schedule where id=$1 and broadcaster_id=$2", id, broadcasterID)
	return err
}

func (q *Query) UpdateSegment(p ScheduleSegment) error {
	_, err := q.DB.NamedExec(generateUpdateSQL("stream_schedule", []string{"id"}, p), p)
	return err
}

func (q *Query) GetVacations(p ScheduleSegment) (ScheduleVacation, error) {
	v := ScheduleVacation{}
	err := q.DB.Get(&v, "select id,starttime,endtime from stream_schedule where is_vacation=true and datetime(endtime) > datetime('now') and broadcaster_id= $1 limit 1", p.UserID)
	if err == sql.ErrNoRows {
		return v, nil
	}
	return v, err
}
