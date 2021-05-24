// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package endpoint_name

import (
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
)

var methodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var scopesByMethod = map[string][]string{
	http.MethodGet:    []string{},
	http.MethodPost:   []string{},
	http.MethodDelete: []string{},
	http.MethodPatch:  []string{},
	http.MethodPut:    []string{},
}

var db database.CLIDatabase

type Endpoint struct{}

func (e Endpoint) GetPath() string { return "/endpoint" }

func (e Endpoint) GetRequiredScopes(method string) []string {
	return scopesByMethod[method]
}

func (e Endpoint) ValidMethod(method string) bool {
	return methodsSupported[method]
}

func (e Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	w.WriteHeader(200)
}
