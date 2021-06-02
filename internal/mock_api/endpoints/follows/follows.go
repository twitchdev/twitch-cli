// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package follows

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mattn/go-sqlite3"
	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/models"
)

var methodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   true,
	http.MethodDelete: true,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var scopesByMethod = map[string][]string{
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

var db database.CLIDatabase

type Endpoint struct{}

func (e Endpoint) Path() string { return "/users/follows" }

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
	q := r.URL.Query()
	to := q["to_id"]
	from := q["from_id"]

	if len(to) == 0 && len(from) == 0 {
		w.WriteHeader(400)
		return
	}

	// adds a blank string to the end of the array- so will always have at least a 0 index attribute
	to = append(to, "")
	from = append(from, "")

	req := database.UserRequestParams{
		UserID:        from[0],
		BroadcasterID: to[0],
	}

	dbr, err := db.NewQuery(r, 100).GetFollows(req)
	if dbr == nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	f := dbr.Data.([]database.Follow)

	log.Printf("%v", err)

	if len(f) == 0 {
		f = []database.Follow{}
	}

	body := models.APIResponse{
		Data:  f,
		Total: &dbr.Total,
	}
	if dbr.Cursor != "" {
		log.Printf("%#v", &dbr)
		body.Pagination = &models.APIPagination{
			Cursor: &dbr.Cursor,
		}
	}

	json, _ := json.Marshal(body)
	w.Write(json)
}

func deleteFollows(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	to := q["to_id"]
	from := q["from_id"]

	if len(to) == 0 || len(from) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := db.NewQuery(r, 100).DeleteFollow(from[0], to[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func postFollows(w http.ResponseWriter, r *http.Request) {
	var body PostFollowBody
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(b, &body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if body.FromID == "" || body.ToID == "" {
		log.Printf("%#v", body)
		w.WriteHeader(http.StatusBadRequest)
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
