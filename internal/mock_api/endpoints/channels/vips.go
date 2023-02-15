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

var vipsMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   true,
	http.MethodDelete: true,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var vipsScopesByMethod = map[string][]string{
	http.MethodGet:    {"channel:read:vips", "channel:manage:vips"},
	http.MethodPost:   {"channel:manage:vips"},
	http.MethodDelete: {"channel:manage:vips"},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type GetVIPsResponseBody struct {
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	UserLogin string `json:"user_login"`
}

type Vips struct{}

func (e Vips) Path() string { return "/channels/vips" }

func (e Vips) GetRequiredScopes(method string) []string {
	return vipsScopesByMethod[method]
}

func (e Vips) ValidMethod(method string) bool {
	return vipsMethodsSupported[method]
}

func (e Vips) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getVIPs(w, r)
		break
	case http.MethodPost:
		postVIPs(w, r)
		break
	case http.MethodDelete:
		deleteVIPs(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getVIPs(w http.ResponseWriter, r *http.Request) {
	broadcasterID := r.URL.Query().Get("broadcaster_id")
	if broadcasterID == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter broadcaster_id")
		return
	}

	userIDs := r.URL.Query()["user_id"]

	dbr, err := db.NewQuery(r, 100).GetVIPsByBroadcaster(broadcasterID)
	if err != nil {
		mock_errors.WriteServerError(w, "error fetching vips: "+err.Error())
		return
	}

	vips := []GetVIPsResponseBody{}

	for _, vip := range dbr.Data.([]database.VIP) {
		if len(userIDs) != 0 {
			// One or more user_id given in query parameters. Grab only these VIPs.
			found := false
			for _, user := range userIDs {
				if user == vip.UserID {
					found = true
				}
			}
			if !found {
				continue
			}
		}

		userDbr, err := db.NewQuery(r, 100).GetUser(database.User{ID: vip.UserID})
		if err != nil {
			mock_errors.WriteServerError(w, "error fetching user: "+err.Error())
			return
		}

		vips = append(vips, GetVIPsResponseBody{
			UserID:    vip.UserID,
			UserName:  userDbr.DisplayName,
			UserLogin: userDbr.UserLogin,
		})
	}

	bytes, _ := json.Marshal(models.APIResponse{Data: vips})
	w.Write(bytes)
}

func postVIPs(w http.ResponseWriter, r *http.Request) {
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

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter user_id")
		return
	}

	userDbr, err := db.NewQuery(r, 100).GetUser(database.User{ID: userID})
	if err != nil {
		mock_errors.WriteServerError(w, "error fetching user: "+err.Error())
		return
	}

	if userDbr.ID == "" {
		mock_errors.WriteNotFound(w, "The ID in user_id was not found")
		return
	}

	// Check if user is a moderator on the channel
	isModerator := false
	moderatorListDbr, err := db.NewQuery(r, 1000).GetModeratorsForBroadcaster(broadcasterID)
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}
	for _, mod := range moderatorListDbr.Data.([]database.Moderator) {
		if mod.UserID == userID {
			isModerator = true
		}
	}
	if isModerator {
		mock_errors.WriteUnprocessableEntity(w, "The specified user is a moderator. To make them a VIP, you must first remove them as a moderator.")
		return
	}

	// Get VIPs
	isVIP := false
	vipListDbr, err := db.NewQuery(r, 100).GetVIPsByBroadcaster(broadcasterID)
	if err != nil {
		mock_errors.WriteServerError(w, "error fetching vips: "+err.Error())
		return
	}
	for _, vip := range vipListDbr.Data.([]database.VIP) {
		if vip.UserID == userID {
			isVIP = true
		}
	}
	if isVIP {
		mock_errors.WriteUnprocessableEntity(w, "User is already a VIP")
		return
	}

	err = db.NewQuery(r, 100).AddVIP(database.UserRequestParams{BroadcasterID: broadcasterID, UserID: userID})
	if err != nil {
		mock_errors.WriteServerError(w, "error adding vip: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteVIPs(w http.ResponseWriter, r *http.Request) {
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

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter user_id")
		return
	}

	// Get VIPs
	isVIP := false
	vipListDbr, err := db.NewQuery(r, 100).GetVIPsByBroadcaster(broadcasterID)
	if err != nil {
		mock_errors.WriteServerError(w, "error fetching vips: "+err.Error())
		return
	}
	for _, vip := range vipListDbr.Data.([]database.VIP) {
		if vip.UserID == userID {
			isVIP = true
		}
	}
	if !isVIP {
		mock_errors.WriteUnprocessableEntity(w, "User is not a VIP in broadcaster's channel")
		return
	}

	err = db.NewQuery(r, 100).DeleteVIP(broadcasterID, userID)
	if err != nil {
		mock_errors.WriteServerError(w, "error removing vip: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
