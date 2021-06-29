// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package videos

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var videosMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: true,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var videosScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodDelete: {"channel:manage:videos"},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type Videos struct{}

var periodDurationMapping = map[string]time.Duration{
	"all":   0,
	"day":   24 * time.Hour,
	"week":  7 * 24 * time.Hour,
	"month": 30 * 7 * 24 * time.Hour,
}

func (e Videos) Path() string { return "/videos" }

func (e Videos) GetRequiredScopes(method string) []string {
	return videosScopesByMethod[method]
}

func (e Videos) ValidMethod(method string) bool {
	return videosMethodsSupported[method]
}

func (e Videos) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getVideos(w, r)
	case http.MethodDelete:
		deleteVideos(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func getVideos(w http.ResponseWriter, r *http.Request) {
	videos := []database.Video{}
	validPeriods := []string{"all", "day", "week", "month"}
	validSorts := []string{"time", "trending", "views"}
	validTypes := []string{"all", "upload", "archive", "highlight"}

	// required
	ids := r.URL.Query()["id"]
	userID := r.URL.Query().Get("user_id")
	gameID := r.URL.Query().Get("game_id")
	if len(ids) == 0 && userID == "" && gameID == "" {
		mock_errors.WriteBadRequest(w, "one of id or user_id or game_id is required")
		return
	}
	if len(ids) > 0 && (userID != "" || gameID != "") {
		mock_errors.WriteBadRequest(w, "if id is provided, then user_id and game_id are unavailable")
		return
	}

	//optional
	language := r.URL.Query().Get("language")
	period := r.URL.Query().Get("period")
	sort := r.URL.Query().Get("sort")
	videoType := r.URL.Query().Get("type")

	if period == "" {
		period = validPeriods[0]
	}
	if sort == "" {
		sort = validSorts[0]
	}
	if videoType == "" {
		videoType = validTypes[0]
	}

	if !isOneOf(validPeriods, period) {
		mock_errors.WriteBadRequest(w, fmt.Sprintf("period must be one of %v", strings.Join(validPeriods, " or ")))
		return
	}
	if !isOneOf(validSorts, sort) {
		mock_errors.WriteBadRequest(w, fmt.Sprintf("sort must be one of %v", strings.Join(validSorts, " or ")))
		return
	}
	if !isOneOf(validTypes, videoType) {
		mock_errors.WriteBadRequest(w, fmt.Sprintf("type must be one of %v", strings.Join(validTypes, " or ")))
		return
	}

	if videoType == "all" {
		// turn to empty string to not select based on type
		videoType = ""
	}

	for _, id := range ids {
		dbr, err := db.NewQuery(r, 100).GetVideos(database.Video{ID: id}, util.GetTimestamp().Add(-1*periodDurationMapping[period]).Format(time.RFC3339), "")
		if err != nil {
			mock_errors.WriteServerError(w, "error fetching videos")
			return
		}
		videos = append(videos, dbr.Data.([]database.Video)...)
	}
	params := database.Video{
		BroadcasterID: userID,
		VideoLanguage: language,
		Type:          videoType,
	}
	if gameID != "" {
		params.CategoryID = &gameID
	}
	timePeriod := ""
	if period != "all" {
		timePeriod = util.GetTimestamp().Add(-1 * periodDurationMapping[period]).Format(time.RFC3339)
	}
	dbr, err := db.NewQuery(r, 100).GetVideos(params, timePeriod, sort)
	if err != nil {
		mock_errors.WriteServerError(w, "error fetching videos")
		return
	}

	videos = append(videos, dbr.Data.([]database.Video)...)

	body := models.APIResponse{
		Data: videos,
	}

	if dbr.Cursor != "" {
		body.Pagination = &models.APIPagination{
			Cursor: dbr.Cursor,
		}
	}

	j, _ := json.Marshal(body)

	w.Write(j)
}

func deleteVideos(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query()["id"]

	if len(ids) == 0 || len(ids) > 5 {
		mock_errors.WriteBadRequest(w, "id is required")
	}

	for _, id := range ids {
		err := db.NewQuery(r, 100).DeleteVideo(id)
		if err != nil {
			mock_errors.WriteServerError(w, "error deleting videos")
			return
		}
	}

	bytes, _ := json.Marshal(models.APIResponse{Data: ids})
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
