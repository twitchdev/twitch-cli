// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package streams

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var markersMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   true,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var markersScopesByMethod = map[string][]string{
	http.MethodGet:    {"user:read:broadcast", "channel:manage:broadcast"},
	http.MethodPost:   {"channel:manage:broadcast"},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type Markers struct{}

type MarkerPostBody struct {
	UserID      string `json:"user_id"`
	Description string `json:"description"`
}

type MarkerPostResponse struct {
	ID              string `json:"id"`
	CreatedAt       string `json:"created_at"`
	Description     string `json:"description"`
	PositionSeconds int    `json:"position_seconds"`
}

func (e Markers) Path() string { return "/streams/markers" }

func (e Markers) GetRequiredScopes(method string) []string {
	return markersScopesByMethod[method]
}

func (e Markers) ValidMethod(method string) bool {
	return markersMethodsSupported[method]
}

func (e Markers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getMarkers(w, r)
		break
	case http.MethodPost:
		postMarkers(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getMarkers(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	userID := r.URL.Query().Get("user_id")
	videoID := r.URL.Query().Get("video_id")

	if userID == "" && videoID == "" {
		mock_errors.WriteBadRequest(w, "one of user_id and video_id is required")
		return
	}
	if userID != "" && videoID != "" {
		mock_errors.WriteBadRequest(w, "only one of user_id and video_id is allowed")
		return
	}
	if videoID == "" && userCtx.UserID != userID {
		mock_errors.WriteBadRequest(w, "user_id must match token")
		return
	}

	dbr, err := db.NewQuery(r, 100).GetStreamMarkers(database.StreamMarker{BroadcasterID: userID, VideoID: videoID})
	if err != nil {
		println(err.Error())
		mock_errors.WriteServerError(w, "error fetching markers")
		return
	}
	markerResponse := dbr.Data.([]database.StreamMarkerUser)

	json.NewEncoder(w).Encode(models.APIResponse{Data: markerResponse})
}

func postMarkers(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)

	body := MarkerPostBody{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		mock_errors.WriteBadRequest(w, "error reading body")
		return
	}
	if userCtx.UserID != body.UserID {
		mock_errors.WriteBadRequest(w, "user_id does not match token")
		return
	}

	dbr, err := db.NewQuery(r, 100).GetStream(database.Stream{UserID: body.UserID})
	if err != nil {
		mock_errors.WriteServerError(w, "error fetching videos")
		return
	}
	streams := dbr.Data.([]database.Stream)
	if len(streams) == 0 {
		mock_errors.WriteBadRequest(w, "user is not live")
		return
	}

	dbr, err = db.NewQuery(r, 100).GetVideos(database.Video{BroadcasterID: body.UserID}, "", "time")
	if err != nil {
		mock_errors.WriteServerError(w, "error retrieving videos")
		return
	}

	videos := dbr.Data.([]database.Video)
	if len(videos) == 0 {
		mock_errors.WriteBadRequest(w, "user does not have a recent VOD")
		return
	}
	sm := database.StreamMarker{
		ID:              util.RandomGUID(),
		VideoID:         videos[0].ID,
		BroadcasterID:   userCtx.UserID,
		PositionSeconds: int(util.RandomInt(60 * 60)),
		Description:     body.Description,
		CreatedAt:       util.GetTimestamp().Format(time.RFC3339),
	}

	err = db.NewQuery(r, 100).InsertStreamMarker(sm)
	if err != nil {
		println(err.Error())
		mock_errors.WriteServerError(w, "error inserting marker")
		return
	}

	json.NewEncoder(w).Encode(models.APIResponse{Data: []MarkerPostResponse{{
		ID:              sm.ID,
		PositionSeconds: sm.PositionSeconds,
		Description:     sm.Description,
		CreatedAt:       sm.CreatedAt,
	}}})
}
