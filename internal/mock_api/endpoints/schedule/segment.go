// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package schedule

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/twitchdev/twitch-cli/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var scheduleSegmentMethodsSupported = map[string]bool{
	http.MethodGet:    false,
	http.MethodPost:   true,
	http.MethodDelete: true,
	http.MethodPatch:  true,
	http.MethodPut:    false,
}

var scheduleSegmentScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {"channel:manage:schedule"},
	http.MethodDelete: {"channel:manage:schedule"},
	http.MethodPatch:  {"channel:manage:schedule"},
	http.MethodPut:    {},
}

var f = false

type ScheduleSegment struct{}

type SegmentPatchAndPostBody struct {
	StartTime   string  `json:"start_time"`
	Timezone    string  `json:"timezone"`
	IsRecurring *bool   `json:"is_recurring"`
	Duration    string  `json:"duration"`
	CategoryID  *string `json:"category_id"`
	Title       string  `json:"title"`
	IsCanceled  *bool   `json:"is_canceled"`
}

func (e ScheduleSegment) Path() string { return "/schedule/segment" }

func (e ScheduleSegment) GetRequiredScopes(method string) []string {
	return scheduleSegmentScopesByMethod[method]
}

func (e ScheduleSegment) ValidMethod(method string) bool {
	return scheduleSegmentMethodsSupported[method]
}

func (e ScheduleSegment) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodPost:
		e.postSegment(w, r)
	case http.MethodDelete:
		e.deleteSegment(w, r)
	case http.MethodPatch:
		e.patchSegment(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (e ScheduleSegment) postSegment(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	duration := 240

	if !userCtx.MatchesBroadcasterIDParam(r) {
		mock_errors.WriteUnauthorized(w, "User token does not match broadcaster_id parameter")
		return
	}
	var body SegmentPatchAndPostBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		mock_errors.WriteBadRequest(w, "Error parsing body")
		return
	}

	if body.StartTime == "" {
		mock_errors.WriteBadRequest(w, "Missing start_time")
		return
	}
	st, err := time.Parse(time.RFC3339, body.StartTime)
	if err != nil {
		mock_errors.WriteBadRequest(w, "Invalid/malformed start_time provided")
		return
	}

	if body.Timezone == "" {
		mock_errors.WriteBadRequest(w, "Missing timezone")
		return
	}
	_, err = time.LoadLocation(body.Timezone)
	if err != nil {
		mock_errors.WriteBadRequest(w, "Invalid timezone provided")
		return
	}

	var isRecurring bool

	if body.IsRecurring == nil {
		isRecurring = false
	} else {
		isRecurring = *body.IsRecurring
	}

	if len(body.Title) > 140 {
		mock_errors.WriteBadRequest(w, "Title must be less than 140 characters")
		return
	}

	if body.Duration != "" {
		duration, err = strconv.Atoi(body.Duration)
		if err != nil {
			mock_errors.WriteBadRequest(w, "Invalid duration provided")
			return
		}
	}
	et := st.Add(time.Duration(duration) * time.Minute)

	segmentID := util.RandomGUID()
	eventID := base64.RawStdEncoding.EncodeToString([]byte(fmt.Sprintf("%v\\%v", segmentID, st)))
	segment := database.ScheduleSegment{
		ID:          eventID,
		StartTime:   st.UTC().Format(time.RFC3339),
		EndTime:     et.UTC().Format(time.RFC3339),
		IsRecurring: isRecurring,
		IsVacation:  false,
		CategoryID:  body.CategoryID,
		Title:       body.Title,
		UserID:      userCtx.UserID,
		IsCanceled:  &f,
	}
	err = db.NewQuery(nil, 100).InsertSchedule(segment)
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}
	if isRecurring {
		// just a years worth of recurring events; mock data
		for i := 0; i < 52; i++ {
			weekAdd := (i + 1) * 7 * 24
			startTime := time.Now().Add(time.Duration(weekAdd) * time.Hour).UTC()
			endTime := time.Now().Add(time.Duration(weekAdd) * time.Hour).UTC()
			eventID := base64.RawStdEncoding.EncodeToString([]byte(fmt.Sprintf("%v\\%v", segmentID, startTime)))

			s := database.ScheduleSegment{
				ID:          eventID,
				StartTime:   startTime.Format(time.RFC3339),
				EndTime:     endTime.Format(time.RFC3339),
				IsRecurring: isRecurring,
				IsVacation:  false,
				CategoryID:  body.CategoryID,
				Title:       body.Title,
				UserID:      userCtx.UserID,
				IsCanceled:  &f,
			}

			err := db.NewQuery(nil, 100).InsertSchedule(s)
			if err != nil {
				mock_errors.WriteServerError(w, err.Error())
				return
			}
		}
	}
	dbr, err := db.NewQuery(nil, 100).GetSchedule(database.ScheduleSegment{ID: eventID}, time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC))
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}
	b := dbr.Data.(database.Schedule)

	if b.Vacation.StartTime == "" && b.Vacation.EndTime == "" {
		b.Vacation = nil
	}

	apiResponse := models.APIResponse{
		Data: b,
	}

	bytes, _ := json.Marshal(apiResponse)
	w.Write(bytes)
}

func (e ScheduleSegment) deleteSegment(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	id := r.URL.Query().Get("id")
	if !userCtx.MatchesBroadcasterIDParam(r) {
		mock_errors.WriteUnauthorized(w, "User token does not match broadcaster_id parameter")
		return
	}

	if id == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter id")
		return
	}

	err := db.NewQuery(nil, 100).DeleteSegment(id, userCtx.UserID)
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (e ScheduleSegment) patchSegment(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	id := r.URL.Query().Get("id")
	if !userCtx.MatchesBroadcasterIDParam(r) {
		mock_errors.WriteUnauthorized(w, "User token does not match broadcaster_id parameter")
		return
	}
	if id == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter id")
		return
	}

	dbr, err := db.NewQuery(nil, 100).GetSchedule(database.ScheduleSegment{ID: id, UserID: userCtx.UserID}, time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC))
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}
	b := dbr.Data.(database.Schedule)

	if len(b.Segments) == 0 {
		mock_errors.WriteBadRequest(w, "Invalid ID requested")
		return
	}
	segment := b.Segments[0]

	var body SegmentPatchAndPostBody
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		mock_errors.WriteBadRequest(w, "Error parsing body")
		return
	}

	// start_time
	st, err := time.Parse(time.RFC3339, segment.StartTime)
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}
	if body.StartTime != "" {
		st, err = time.Parse(time.RFC3339, body.StartTime)
		if err != nil {
			mock_errors.WriteBadRequest(w, "Error parsing start_time")
			return
		}
	}

	// is_canceled
	isCanceled := false
	if body.IsCanceled != nil {
		isCanceled = *body.IsCanceled
	}

	// timezone
	if body.Timezone != "" {
		_, err := time.LoadLocation(body.Timezone)
		if err != nil {
			mock_errors.WriteBadRequest(w, "Error parsing timezone")
			return
		}
	}

	// title
	title := segment.Title
	if body.Title != "" {
		if len(body.Title) > 140 {
			mock_errors.WriteBadRequest(w, "Title must be less than 140 characters")
			return
		}
		title = body.Title
	}

	// duration
	et, err := time.Parse(time.RFC3339, segment.EndTime)
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}
	if body.Duration != "" {
		duration, err := strconv.Atoi(body.Duration)
		if err != nil {
			mock_errors.WriteBadRequest(w, "Invalid duration provided")
			return
		}

		et = st.Add(time.Duration(duration) * time.Minute)
	}

	s := database.ScheduleSegment{
		ID:         segment.ID,
		StartTime:  st.UTC().Format(time.RFC3339),
		EndTime:    et.UTC().Format(time.RFC3339),
		IsCanceled: &isCanceled,
		Title:      title,
	}

	err = db.NewQuery(r, 20).UpdateSegment(s)
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}

	dbr, err = db.NewQuery(nil, 100).GetSchedule(database.ScheduleSegment{ID: segment.ID}, time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC))
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}
	b = dbr.Data.(database.Schedule)

	if b.Vacation.StartTime == "" && b.Vacation.EndTime == "" {
		b.Vacation = nil
	}

	apiResponse := models.APIResponse{
		Data: b,
	}

	bytes, _ := json.Marshal(apiResponse)
	w.Write(bytes)
}
