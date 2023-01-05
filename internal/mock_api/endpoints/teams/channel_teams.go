// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package teams

import (
	"encoding/json"
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
)

var channelTeamsMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var channelTeamsScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type ChannelTeams struct{}
type ChannelTeamResponse struct {
	ID                 string  `json:"id"`
	BackgroundImageUrl *string `json:"background_image_url"`
	Banner             *string `json:"banner"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
	Info               string  `json:"info"`
	ThumbnailURL       string  `json:"thumbnail_url"`
	TeamName           string  `json:"team_name"`
	TeamDisplayName    string  `json:"team_display_name"`
	BroadcasterID      string  `json:"broadcaster_id"`
	BroadcasterName    string  `json:"broadcaster_name"`
	BroadcasterLogin   string  `json:"broadcaster_login"`
}

func (e ChannelTeams) Path() string { return "/teams/channel" }

func (e ChannelTeams) GetRequiredScopes(method string) []string {
	return channelTeamsScopesByMethod[method]
}

func (e ChannelTeams) ValidMethod(method string) bool {
	return channelTeamsMethodsSupported[method]
}

func (e ChannelTeams) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getChannelTeams(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getChannelTeams(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query().Get("broadcaster_id")) == 0 {
		mock_errors.WriteBadRequest(w, "broadcaster_id is required")
		return
	}

	// Get user information
	userdbr, err := db.NewQuery(r, 1).GetUser(database.User{ID: r.URL.Query().Get("broadcaster_id")})
	if err != nil {
		mock_errors.WriteServerError(w, "error fetching user")
		return
	}
	broadcasterID := r.URL.Query().Get("broadcaster_id")
	broadcasterLogin := userdbr.UserLogin
	broadcasterName := userdbr.DisplayName

	// Get team information
	dbr, err := db.NewQuery(r, 100).GetTeamByBroadcaster(r.URL.Query().Get("broadcaster_id"))
	if err != nil {
		mock_errors.WriteServerError(w, "error fetching team")
		return
	}
	team := dbr.Data.([]database.Team)
	if len(team) == 0 {
		dbr.Data = make([]database.Team, 0)
	}
	response := []ChannelTeamResponse{}
	for _, t := range team {
		response = append(response, ChannelTeamResponse{
			ID:                 t.ID,
			Info:               t.Info,
			BackgroundImageUrl: t.BackgroundImageUrl,
			Banner:             t.Banner,
			CreatedAt:          t.CreatedAt,
			UpdatedAt:          t.UpdatedAt,
			ThumbnailURL:       t.ThumbnailURL,
			TeamName:           t.TeamName,
			TeamDisplayName:    t.TeamDisplayName,
			BroadcasterID:      broadcasterID,
			BroadcasterLogin:   broadcasterLogin,
			BroadcasterName:    broadcasterName,
		})
	}

	bytes, _ := json.Marshal(models.APIResponse{Data: response})
	w.Write(bytes)
}
