// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package schedule

import (
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
)

var scheduleICalMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var scheduleICalScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type ScheduleICal struct{}

func (e ScheduleICal) Path() string { return "/schedule/icalendar" }

func (e ScheduleICal) GetRequiredScopes(method string) []string {
	return scheduleICalScopesByMethod[method]
}

func (e ScheduleICal) ValidMethod(method string) bool {
	return scheduleICalMethodsSupported[method]
}

func (e ScheduleICal) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		e.getIcal(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// stubbed with fake data for now, since .ics generation libraries are far and few between for golang
// and it's just useful for mock data
func (e ScheduleICal) getIcal(w http.ResponseWriter, r *http.Request) {
	broadcaster := r.URL.Query().Get("broadcaster_id")
	if broadcaster == "" {
		mock_errors.WriteBadRequest(w, "Missing required paramater broadaster_id")
		return
	}

	body :=
		`BEGIN:VCALENDAR
PRODID:-//twitch.tv//StreamSchedule//1.0
VERSION:2.0
CALSCALE:GREGORIAN
REFRESH-INTERVAL;VALUE=DURATION:PT1H
NAME:TwitchDev
BEGIN:VEVENT
UID:e4acc724-371f-402c-81ca-23ada79759d4
DTSTAMP:20210323T040131Z
DTSTART;TZID=/America/New_York:20210701T140000
DTEND;TZID=/America/New_York:20210701T150000
SUMMARY:TwitchDev Monthly Update // July 1, 2021
DESCRIPTION:Science & Technology.
CATEGORIES:Science & Technology
END:VEVENT
END:VCALENDAR`
	w.Header().Set("Content-Type", "text/calendar")
	w.Write([]byte(body))
}
