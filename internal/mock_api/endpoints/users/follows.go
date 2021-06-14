// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package users

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mattn/go-sqlite3"
	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
)

var followMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   true,
	http.MethodDelete: true,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var followScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {"user:edit:follows"},
	http.MethodDelete: {"user:edit:follows"},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type PostFollowBody struct {
	ToID   string `json:"to_id"`
	FromID string `json:"from_id"`
}

type FollowsEndpoint struct{}

func (e FollowsEndpoint) Path() string { return "/users/follows" }

func (e FollowsEndpoint) GetRequiredScopes(method string) []string {
	return followScopesByMethod[method]
}

func (e FollowsEndpoint) ValidMethod(method string) bool {
	return followMethodsSupported[method]
}

func (e FollowsEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getFollows(w, r)
	case http.MethodPost:
		postFollows(w, r)
	case http.MethodDelete:
		deleteFollows(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func getFollows(w http.ResponseWriter, r *http.Request) {
	to := r.URL.Query().Get("to_id")
	from := r.URL.Query().Get("from_id")

	if len(to) == 0 && len(from) == 0 {
		mock_errors.WriteBadRequest(w, "one of to_id or from_id is required")
		return
	}

	req := database.UserRequestParams{
		UserID:        from,
		BroadcasterID: to,
	}

	dbr, err := db.NewQuery(r, 100).GetFollows(req)
	if dbr == nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	f := dbr.Data.([]database.Follow)

	if len(f) == 0 {
		f = []database.Follow{}
	}

	body := models.APIResponse{
		Data:  f,
		Total: &dbr.Total,
	}
	if dbr != nil && dbr.Cursor != "" {
		log.Printf("%#v", &dbr)
		body.Pagination = &models.APIPagination{
			Cursor: dbr.Cursor,
		}
	}

	json, _ := json.Marshal(body)
	w.Write(json)
}

func deleteFollows(w http.ResponseWriter, r *http.Request) {
	to := r.URL.Query().Get("to_id")
	from := r.URL.Query().Get("from_id")

	if len(to) == 0 || len(from) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := db.NewQuery(r, 100).DeleteFollow(from, to)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func postFollows(w http.ResponseWriter, r *http.Request) {
	var body PostFollowBody

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		mock_errors.WriteBadRequest(w, "error reading body")
		return
	}
	if body.FromID == "" || body.ToID == "" {
		mock_errors.WriteBadRequest(w, "from_id and to_id are required")
		return
	}

	err = db.NewQuery(r, 100).AddFollow(database.UserRequestParams{UserID: body.FromID, BroadcasterID: body.ToID})
	if err != nil {
		if database.DatabaseErrorIs(err, sqlite3.ErrConstraintForeignKey) || database.DatabaseErrorIs(err, sqlite3.ErrConstraintUnique) || database.DatabaseErrorIs(err, sqlite3.ErrConstraintPrimaryKey) {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		log.Printf("%#v\n%#v", err, sqlite3.ErrConstraintUnique)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}
