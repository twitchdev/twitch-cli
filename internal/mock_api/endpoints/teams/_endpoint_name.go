// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package endpoint_name

import (
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
)

var endpointMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var endpointScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type Endpoint struct{}

func (e Endpoint) Path() string { return "/endpoint" }

func (e Endpoint) GetRequiredScopes(method string) []string {
	return endpointScopesByMethod[method]
}

func (e Endpoint) ValidMethod(method string) bool {
	return endpointMethodsSupported[method]
}

func (e Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	w.WriteHeader(200)
}
