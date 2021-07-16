// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package schedule

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var scheduleSettingsMethodsSupported = map[string]bool{
	http.MethodGet:    false,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  true,
	http.MethodPut:    false,
}

var scheduleSettingsScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {"channel:manage:schedule"},
	http.MethodPut:    {},
}

type ScheduleSettings struct{}

type PatchSettingsBody struct {
	IsVacationEnabled *bool  `json:"is_vacation_enabled"`
	VacationStartTime string `json:"vacation_start_time"`
	VacationEndTime   string `json:"vacation_end_time"`
	Timezone          string `json:"timezone"`
}

func (e ScheduleSettings) Path() string { return "/schedule/settings" }

func (e ScheduleSettings) GetRequiredScopes(method string) []string {
	return scheduleSettingsScopesByMethod[method]
}

func (e ScheduleSettings) ValidMethod(method string) bool {
	return scheduleSettingsMethodsSupported[method]
}

func (e ScheduleSettings) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodPatch:
		e.patchSchedule(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (e ScheduleSettings) patchSchedule(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	if !userCtx.MatchesBroadcasterIDParam(r) {
		mock_errors.WriteUnauthorized(w, "User token does not match broadcaster_id parameter")
		return
	}

	vacation, err := db.NewQuery(r, 100).GetVacations(database.ScheduleSegment{UserID: userCtx.UserID})
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}

	var body PatchSettingsBody
	err = json.NewDecoder(r.Body).Decode(&body)

	if body.IsVacationEnabled == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if *body.IsVacationEnabled == false {
		if vacation.ID != "" {
			err := db.NewQuery(r, 100).DeleteSegment(vacation.ID, userCtx.UserID)
			if err != nil {
				mock_errors.WriteServerError(w, err.Error())
				return
			}
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if vacation.ID != "" && *body.IsVacationEnabled == true {
		mock_errors.WriteBadRequest(w, "Existing vacation already exists")
		return
	}

	if body.Timezone == "" || body.VacationStartTime == "" || body.VacationEndTime == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter")
		return
	}

	_, err = time.LoadLocation(body.Timezone)
	if err != nil {
		mock_errors.WriteBadRequest(w, "Invalid timezone requested")
		return
	}

	st, err := time.Parse(time.RFC3339, body.VacationStartTime)
	if err != nil {
		mock_errors.WriteBadRequest(w, "Invalid vacation_start_time requested")
		return
	}

	et, err := time.Parse(time.RFC3339, body.VacationEndTime)
	if err != nil {
		mock_errors.WriteBadRequest(w, "Invalid vacation_end_time requested")
		return
	}
	f := false
	err = db.NewQuery(r, 100).InsertSchedule(database.ScheduleSegment{
		ID:          base64.RawStdEncoding.EncodeToString([]byte(fmt.Sprintf("%v\\%v", util.RandomGUID(), st))),
		StartTime:   st.UTC().Format(time.RFC3339),
		EndTime:     et.UTC().Format(time.RFC3339),
		IsVacation:  true,
		IsRecurring: false,
		IsCanceled:  &f,
		UserID:      userCtx.UserID,
	})
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
