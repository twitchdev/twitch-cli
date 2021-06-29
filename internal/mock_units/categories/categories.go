// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package categories

import (
	"encoding/json"
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
)

type Endpoint struct{}

var db database.CLIDatabase

func (e Endpoint) Path() string { return "/categories" }

func (e Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getCategories(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getCategories(w http.ResponseWriter, r *http.Request) {
	c, err := db.NewQuery(r, 100).GetCategories(database.Category{})
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}
	j, err := json.Marshal(c)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(j)
}
