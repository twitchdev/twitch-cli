// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package streams

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
)

var streamsMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var streamsScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type Streams struct{}

type PaginationCursor struct {
	Limit      int `json:"l"`
	PrevOffset int `json:"o"`
}

func (e Streams) Path() string { return "/streams" }

func (e Streams) GetRequiredScopes(method string) []string {
	return streamsScopesByMethod[method]
}

func (e Streams) ValidMethod(method string) bool {
	return streamsMethodsSupported[method]
}

func (e Streams) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getStreams(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getStreams(w http.ResponseWriter, r *http.Request) {
	gameIDs := r.URL.Query()["game_id"]
	languages := r.URL.Query()["language"]
	userIDs := r.URL.Query()["user_id"]
	userLogins := r.URL.Query()["user_login"]
	// custom pagination for this endpoint
	start := 0
	end := 0
	limit := 0

	isBefore := false

	pc := PaginationCursor{}
	if r.URL.Query().Get("after") != "" {
		pc = paginationStruct(r.URL.Query().Get("after"))
		start = pc.PrevOffset
	}

	if r.URL.Query().Get("before") != "" {
		isBefore = true
		pc = paginationStruct(r.URL.Query().Get("after"))
		start = pc.PrevOffset
	}

	if r.URL.Query().Get("first") != "" {
		limit, _ = strconv.Atoi(r.URL.Query().Get("first"))
		if limit == 0 {
			limit = 20
		}

	} else {
		limit = 20
	}

	if isBefore {
		end = start
		start = end - limit
		if start < 0 {
			start = 0
			end = limit
		}
	} else {
		end = start + limit
	}

	if len(gameIDs) > 100 || len(languages) > 100 || len(userIDs) > 100 || len(userLogins) > 100 {
		mock_errors.WriteBadRequest(w, "you may only send 100 of each parameter")
		return
	}

	streams := []database.Stream{}
	// get all streams filtered here
	for _, id := range userIDs {
		dbr, err := db.NewQuery(nil, 100).GetStream(database.Stream{UserID: id})
		if err != nil {
			mock_errors.WriteServerError(w, "error fetching streams")
			return
		}
		s := dbr.Data.([]database.Stream)
		streams = append(streams, s...)
	}

	for _, login := range userLogins {
		dbr, err := db.NewQuery(nil, 100).GetStream(database.Stream{UserLogin: login})
		if err != nil {
			mock_errors.WriteServerError(w, "error fetching streams")
			return
		}
		s := dbr.Data.([]database.Stream)
		streams = append(streams, s...)
	}

	// if neither, none of the code above will run, so let's get all streams
	if len(userIDs) == 0 && len(userLogins) == 0 {
		dbr, err := db.NewQuery(nil, 100).GetStream(database.Stream{})
		if err != nil {
			mock_errors.WriteServerError(w, "error fetching streams")
			return
		}
		s := dbr.Data.([]database.Stream)
		streams = append(streams, s...)
	}

	d := []database.Stream{}
	// filter out any not in the streams here, matching production behavior
	for _, s := range streams {
		if isOneOf(languages, s.Language) && isOneOf(gameIDs, s.CategoryID.String) {
			d = append(d, s)
		}
	}

	if end > len(d) {
		end = len(d)
	}
	apiResponse := models.APIResponse{
		Data: d[start:end],
	}

	if len(d[start:end]) == limit && end != len(d) {
		apiResponse.Pagination = &models.APIPagination{
			Cursor: paginationString(limit, end),
		}
	}

	bytes, _ := json.Marshal(apiResponse)
	w.Write(bytes)
}

func isOneOf(listOFAllowed []string, s string) bool {
	if len(listOFAllowed) == 0 {
		return true
	}
	for _, i := range listOFAllowed {
		if s == i {
			return true
		}
	}
	return false
}

func paginationString(limit int, prevOffset int) string {
	body, _ := json.Marshal(PaginationCursor{Limit: limit, PrevOffset: prevOffset})
	return base64.RawURLEncoding.EncodeToString(body)
}

func paginationStruct(cursor string) PaginationCursor {
	c, _ := base64.RawStdEncoding.DecodeString(cursor)
	pc := PaginationCursor{}
	json.Unmarshal(c, &pc)
	return pc
}
