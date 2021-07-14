// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package schedule

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
)

var scheduleMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var scheduleScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type Schedule struct{}

func (e Schedule) Path() string { return "/schedule" }

func (e Schedule) GetRequiredScopes(method string) []string {
	return scheduleScopesByMethod[method]
}

func (e Schedule) ValidMethod(method string) bool {
	return scheduleMethodsSupported[method]
}

func (e Schedule) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		e.getSchedule(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (e Schedule) getSchedule(w http.ResponseWriter, r *http.Request) {
	broadcasterID := r.URL.Query().Get("broadcaster_id")
	queryTime := r.URL.Query().Get("start_time")
	offset := r.URL.Query().Get("utc_offset")
	ids := r.URL.Query()["id"]
	schedule := database.Schedule{}
	startTime := time.Now().UTC()
	apiResponse := models.APIResponse{}

	if broadcasterID == "" {
		mock_errors.WriteBadRequest(w, "Required parameter broadcaster_id is missing")
		return
	}

	if queryTime != "" {
		st, err := time.Parse(time.RFC3339, queryTime)
		if err != nil {
			mock_errors.WriteBadRequest(w, "Parameter start_time is in an invalid format")
			return
		}
		startTime = st.UTC()
	}

	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			mock_errors.WriteBadRequest(w, "Error decoding parameter offset")
			return
		}
		tz := time.FixedZone("", o*60)
		startTime = startTime.In(tz)
	}

	segments := []database.ScheduleSegment{}
	if len(ids) > 0 {
		if len(ids) > 100 {
			mock_errors.WriteBadRequest(w, "Parameter id may only have a maximum of 100 values")
			return
		}
		for _, id := range ids {
			dbr, err := db.NewQuery(r, 25).GetSchedule(database.ScheduleSegment{ID: id, UserID: broadcasterID}, startTime)
			if err != nil {
				mock_errors.WriteServerError(w, err.Error())
			}
			response := dbr.Data.(database.Schedule)
			schedule = response
			segments = append(segments, response.Segments...)
		}
		schedule.Segments = segments
		apiResponse = models.APIResponse{
			Data: schedule,
		}
	} else {
		dbr, err := db.NewQuery(r, 25).GetSchedule(database.ScheduleSegment{UserID: broadcasterID}, startTime)
		if err != nil {
			mock_errors.WriteServerError(w, err.Error())
			return
		}
		response := dbr.Data.(database.Schedule)
		segments = append(segments, response.Segments...)
		schedule = response
		schedule.Segments = segments
		apiResponse = models.APIResponse{
			Data: schedule,
		}

		if len(schedule.Segments) == dbr.Limit {
			apiResponse.Pagination = &models.APIPagination{
				Cursor: dbr.Cursor,
			}
		}
	}

	bytes, _ := json.Marshal(apiResponse)
	w.Write(bytes)
}
