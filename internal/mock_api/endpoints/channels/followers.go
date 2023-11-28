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

var followersMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var followersScopesByMethod = map[string][]string{
	http.MethodGet:    {"moderator:read:followers"},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type FollowersEndpoint struct{}

type GetFollowersEndpointResponseData struct {
	UserID     string `json:"user_id"`
	UserLogin  string `json:"user_login"`
	UserName   string `json:"user_name"`
	FollowedAt string `json:"followed_at"`
}

func (e FollowersEndpoint) Path() string { return "/channels/followers" }

func (e FollowersEndpoint) GetRequiredScopes(method string) []string {
	return followersScopesByMethod[method]
}

func (e FollowersEndpoint) ValidMethod(method string) bool {
	return followersMethodsSupported[method]
}

func (e FollowersEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getFollowers(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func getFollowers(w http.ResponseWriter, r *http.Request) {
	user_id := r.URL.Query().Get("user_id")
	broadcaster_id := r.URL.Query().Get("broadcaster_id")

	if broadcaster_id == "" {
		mock_errors.WriteBadRequest(w, "The broadcaster_id query parameter is required")
		return
	}

	/// If user_id used:
	/// - Check for moderator:read:followers scope is used and if broadcaster_id == access token user id
	/// - If false on above, check if user_id is a moderator for broadcaster_id

	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	trustedUser := false

	if userCtx.MatchesBroadcasterIDParam(r) {
		trustedUser = true
	} else {
		// Check if user a moderator instead
		moderatorListDbr, err := db.NewQuery(r, 1000).GetModeratorsForBroadcaster(broadcaster_id)
		if err != nil {
			mock_errors.WriteServerError(w, err.Error())
			return
		}
		modFound := false
		for _, mod := range moderatorListDbr.Data.([]database.Moderator) {
			if mod.UserID == user_id {
				modFound = true // Moderator found
			}
		}
		if modFound {
			trustedUser = true
		}
	}

	if !userCtx.HasScope("moderator:read:followers") {
		// Doesn't matter if they are broadcaster, moderator, or regular: if they don't have the scope they don't get the data.
		trustedUser = false
	}

	if user_id != "" && !trustedUser {
		mock_errors.WriteUnauthorized(w, "When user_id param is provided, user_id must match the User ID in the access token, and the access token must have the scope moderator:read:followers")
		return
	}

	req := database.UserRequestParams{
		UserID:        user_id,
		BroadcasterID: broadcaster_id,
	}

	dbr, err := db.NewQuery(r, 100).GetFollows(req, false)
	if dbr == nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}

	// Build list of who the user is following
	follows := []GetFollowersEndpointResponseData{}
	if trustedUser {
		for _, f := range dbr.Data.([]database.Follow) {
			follows = append(follows, GetFollowersEndpointResponseData{
				UserID:     f.ViewerID,
				UserLogin:  f.ViewerLogin,
				UserName:   f.ViewerName,
				FollowedAt: f.FollowedAt,
			})
		}
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
