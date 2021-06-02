// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package subscriptions

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/models"
)

var methodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var scopesByMethod = map[string][]string{
	http.MethodGet:    {"channel:read:subscriptions"},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

var db database.CLIDatabase

type Endpoint struct{}

func (e Endpoint) Path() string { return "/subscriptions" }

func (e Endpoint) GetRequiredScopes(method string) []string {
	return scopesByMethod[method]
}

func (e Endpoint) ValidMethod(method string) bool {
	return methodsSupported[method]
}

func (e Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getSubscriptions(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func getSubscriptions(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	bid := q.Get("broadcaster_id")

	if len(bid) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s := database.Subscription{
		BroadcasterID: bid,
	}

	res, err := db.NewQuery(r, 100).GetSubscriptions(s)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	if len(res.Data.([]database.Subscription)) == 0 {
		res.Data = []database.Subscription{}
	}

	body := models.APIResponse{
		Data: res.Data,
	}

	if res.Cursor != "" {
		pag := &models.APIPagination{
			Cursor: &res.Cursor,
		}
		body.Pagination = pag
	}

	j, _ := json.Marshal(body)

	w.Write(j)
}
