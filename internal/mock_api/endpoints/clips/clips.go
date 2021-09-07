// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package clips

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var clipsMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   true,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var clipsScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {"clips:edit"},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type Clips struct{}

type CreateClipsResponse struct {
	ID      string `json:"id"`
	EditURL string `json:"edit_url"`
}

func (e Clips) Path() string { return "/clips" }

func (e Clips) GetRequiredScopes(method string) []string {
	return clipsScopesByMethod[method]
}

func (e Clips) ValidMethod(method string) bool {
	return clipsMethodsSupported[method]
}

func (e Clips) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)
	switch r.Method {
	case http.MethodGet:
		getClips(w, r)
		break
	case http.MethodPost:
		postClips(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getClips(w http.ResponseWriter, r *http.Request) {
	broadcasterID := r.URL.Query().Get("broadcaster_id")
	gameID := r.URL.Query().Get("game_id")
	id := r.URL.Query().Get("id")

	startedAt := r.URL.Query().Get("started_at")
	endedAt := r.URL.Query().Get("ended_at")

	if broadcasterID == "" && gameID == "" && id == "" {
		mock_errors.WriteBadRequest(w, "one of broadcaster_id, game_id, or id is required")
		return
	}

	if endedAt != "" && startedAt == "" {
		mock_errors.WriteBadRequest(w, "started_at is required if ended_at is provided")
		return
	}

	if startedAt != "" && endedAt == "" {
		sa, _ := time.Parse(time.RFC3339, startedAt)
		endedAt = sa.Add(7 * 24 * time.Hour).Format(time.RFC3339)
	}

	dbr, err := db.NewQuery(r, 100).GetClips(database.Clip{ID: id, BroadcasterID: broadcasterID, GameID: gameID}, startedAt, endedAt)
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}

	clips := dbr.Data.([]database.Clip)
	apiResponse := models.APIResponse{
		Data: clips,
	}

	if len(apiResponse.Data.([]database.Clip)) == 0 {
		apiResponse.Data = []string{}
	}

	if dbr.Limit == len(dbr.Data.([]database.Clip)) {
		apiResponse.Pagination = &models.APIPagination{
			Cursor: dbr.Cursor,
		}
	}

	bytes, _ := json.Marshal(apiResponse)
	w.Write(bytes)
}

func postClips(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	broadcasterID := r.URL.Query().Get("broadcaster_id")
	// has_delay has no effect in the mock API
	if broadcasterID == "" {
		mock_errors.WriteBadRequest(w, "broadcaster_id is required")
		return
	}

	dbr, err := db.NewQuery(r, 100).GetStream(database.Stream{UserID: broadcasterID})
	if err != nil {
		mock_errors.WriteServerError(w, "Error fetching stream status for user")
		return
	}

	streams := dbr.Data.([]database.Stream)

	if len(streams) == 0 {
		mock_errors.WriteBadRequest(w, "Requested broadcaster is not live")
		return
	}

	id := util.RandomGUID()

	clip := database.Clip{
		ID:            id,
		BroadcasterID: broadcasterID,
		CreatorID:     userCtx.UserID,
		Title:         streams[0].Title,
		Language:      streams[0].Language,
		GameID:        streams[0].RealCategoryID,
		ViewCount:     0,
		VideoID:       "",
		Duration:      33.3,
		CreatedAt:     util.GetTimestamp().Format(time.RFC3339),
	}

	err = db.NewQuery(r, 100).InsertClip(clip)
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}

	bytes, _ := json.Marshal(models.APIResponse{Data: []CreateClipsResponse{
		{
			ID:      id,
			EditURL: fmt.Sprintf("http://clips.twitch.tv/%v/edit", id),
		},
	}})

	w.Write(bytes)
}
