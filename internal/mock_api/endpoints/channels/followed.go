// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package channels

import (
	"encoding/json"
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
)

var followedMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var followedScopesByMethod = map[string][]string{
	http.MethodGet:    {"user:read:follows"},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type FollowedEndpoint struct{}

type GetFollowedEndpointResponseData struct {
	BroadcasterID    string `json:"broadcaster_id"`
	BroadcasterLogin string `json:"broadcaster_login"`
	BroadcasterName  string `json:"broadcaster_name"`
	FollowedAt       string `json:"followed_at"`
}

func (e FollowedEndpoint) Path() string { return "/channels/followed" }

func (e FollowedEndpoint) GetRequiredScopes(method string) []string {
	return followedScopesByMethod[method]
}

func (e FollowedEndpoint) ValidMethod(method string) bool {
	return followedMethodsSupported[method]
}

func (e FollowedEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getFollowed(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func getFollowed(w http.ResponseWriter, r *http.Request) {
	user_id := r.URL.Query().Get("user_id")
	broadcaster_id := r.URL.Query().Get("broadcaster_id")

	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)

	if user_id == "" {
		mock_errors.WriteBadRequest(w, "The user_id query parameter is required")
		return
	}

	if user_id != userCtx.UserID {
		mock_errors.WriteUnauthorized(w, "user_id does not match User ID in the access token")
		return
	}

	req := database.UserRequestParams{
		UserID:        user_id,
		BroadcasterID: broadcaster_id,
	}

	dbr, err := db.NewQuery(r, 100).GetFollows(req, true)
	if dbr == nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}

	// Build list of who the user is following
	follows := []GetFollowedEndpointResponseData{}
	for _, f := range dbr.Data.([]database.Follow) {
		follows = append(follows, GetFollowedEndpointResponseData{
			BroadcasterID:    f.BroadcasterID,
			BroadcasterLogin: f.BroadcasterLogin,
			BroadcasterName:  f.BroadcasterName,
			FollowedAt:       f.FollowedAt,
		})
	}

	body := models.APIResponse{
		Data:  follows,
		Total: &dbr.Total,
	}
	if dbr != nil && dbr.Cursor != "" {
		body.Pagination = &models.APIPagination{
			Cursor: dbr.Cursor,
		}
	}

	bytes, _ := json.Marshal(body)
	w.Write(bytes)
}
