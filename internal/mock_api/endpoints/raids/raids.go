// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package raids

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var raidsMethodsSupported = map[string]bool{
	http.MethodGet:    false,
	http.MethodPost:   true,
	http.MethodDelete: true,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var raidsScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {"channel:manage:raids"},
	http.MethodDelete: {"channel:manage:raids"},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type GetVIPsResponseBody struct {
	CreatedAt string `json:"created_at"`
	IsMature  bool   `json:"is_mature"`
}

type Raids struct{}

func (e Raids) Path() string { return "/raids" }

func (e Raids) GetRequiredScopes(method string) []string {
	return raidsScopesByMethod[method]
}

func (e Raids) ValidMethod(method string) bool {
	return raidsMethodsSupported[method]
}

func (e Raids) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodPost:
		postRaids(w, r)
		break
	case http.MethodDelete:
		deleteRaids(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func postRaids(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	if !userCtx.MatchesSpecifiedIDParam(r, "from_broadcaster_id") {
		mock_errors.WriteUnauthorized(w, "from_broadcaster_id does not match token")
		return
	}

	fromBroadcasterID := r.URL.Query().Get("from_broadcaster_id")
	if fromBroadcasterID == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter from_broadcaster_id")
		return
	}

	toBroadcasterID := r.URL.Query().Get("to_broadcaster_id")
	if toBroadcasterID == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter to_broadcaster_id")
		return
	}

	if fromBroadcasterID == toBroadcasterID {
		mock_errors.WriteBadRequest(w, "The IDs on from_broadcaster_id and to_broadcaster_id cannot be the same ID")
		return
	}

	// Check if user exists
	user, err := db.NewQuery(r, 100).GetUser(database.User{ID: toBroadcasterID})
	if err != nil {
		mock_errors.WriteServerError(w, "error pulling to_broadcaster_id from user database: "+err.Error())
		return
	}
	if user.ID == "" {
		mock_errors.WriteBadRequest(w, "User specified in to_broadcaster_id doesn't exist")
		return
	}

	rand.Seed(util.GetTimestamp().UnixNano())
	isMature := rand.Float32() < 0.5

	bytes, _ := json.Marshal(models.APIResponse{
		Data: GetVIPsResponseBody{
			CreatedAt: util.GetTimestamp().Format(time.RFC3339Nano),
			IsMature:  isMature,
		},
	})
	w.Write(bytes)

	// There's no real channel handling in the mock API, so we'll just ingest this and say it happened.
	// Right now this means no 409 Conflict handling
}

func deleteRaids(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	if !userCtx.MatchesBroadcasterIDParam(r) {
		mock_errors.WriteUnauthorized(w, "broadcaster_id does not match token")
		return
	}

	broadcasterID := r.URL.Query().Get("broadcaster_id")
	if broadcasterID == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter broadcaster_id")
		return
	}

	// There's no real channel handling in the mock API, so we'll just ingest this and say it happened.
	// Right now this means no 404 Not Found handling

	w.WriteHeader(http.StatusNoContent)
}
