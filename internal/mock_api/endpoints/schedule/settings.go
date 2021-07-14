// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package schedule

import (
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
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

func (e ScheduleSettings) Path() string { return "/schedule/settings" }

func (e ScheduleSettings) GetRequiredScopes(method string) []string {
	return scheduleSettingsScopesByMethod[method]
}

func (e ScheduleSettings) ValidMethod(method string) bool {
	return scheduleSettingsMethodsSupported[method]
}

func (e ScheduleSettings) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	w.WriteHeader(200)
}
