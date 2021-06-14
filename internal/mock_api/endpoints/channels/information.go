// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package channels

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/mattn/go-sqlite3"
	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
)

type Channel struct {
	ID           string `db:"id" json:"broadcaster_id"`
	UserLogin    string `db:"user_login" json:"broadcaster_login"`
	DisplayName  string `db:"display_name" json:"broadcaster_name"`
	CategoryID   string `db:"category_id" json:"game_id"`
	CategoryName string `db:"category_name" json:"game_name" dbi:"false"`
	Title        string `db:"title" json:"title"`
	Language     string `db:"stream_language" json:"stream_language"`
	Delay        int    `dbi:"false" json:"delay"`
}

var informationMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  true,
	http.MethodPut:    false,
}

var informationScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {"channel:manage:broadcast"},
	http.MethodPut:    {},
}

type InformationEndpoint struct{}

type PatchInformationEndpointRequest struct {
	GameID              string `json:"game_id"`
	BroadcasterLanguage string `json:"broadcaster_language"`
	Title               string `json:"title"`
	Delay               *int   `json:"delay"`
}

func (e InformationEndpoint) Path() string { return "/channels" }

func (e InformationEndpoint) GetRequiredScopes(method string) []string {
	return informationScopesByMethod[method]
}

func (e InformationEndpoint) ValidMethod(method string) bool {
	return informationMethodsSupported[method]
}

func (e InformationEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getInformation(w, r)
		break
	case http.MethodPatch:
		patchInformation(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		break
	}
}

func getInformation(w http.ResponseWriter, r *http.Request) {
	broadcasterID := r.URL.Query().Get("broadcaster_id")

	if broadcasterID == "" {
		mock_errors.WriteBadRequest(w, "Broacaster ID is required")
		return
	}

	dbr, err := db.NewQuery(r, 100).GetChannels(database.User{ID: broadcasterID})
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}

	channels := dbr.Data.([]database.User)
	c := convertUsers(channels)

	bytes, _ := json.Marshal(models.APIResponse{
		Data: c,
	})

	w.Write(bytes)
}

func patchInformation(w http.ResponseWriter, r *http.Request) {
	broadcasterID := r.URL.Query().Get("broadcaster_id")
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)

	if broadcasterID == "" {
		mock_errors.WriteBadRequest(w, "Broacaster ID is required")
		return
	}

	if broadcasterID != userCtx.UserID {
		mock_errors.WriteUnauthorized(w, "Broadcaster ID does not match token")
		return
	}

	log.Print("here")
	var params PatchInformationEndpointRequest
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		mock_errors.WriteBadRequest(w, "Error reading body")
		return
	}

	u, err := db.NewQuery(r, 100).GetUser(database.User{ID: broadcasterID})

	var gameID = u.CategoryID
	if params.GameID == "" || params.GameID == "0" {
		gameID = sql.NullString{}
	} else if params.GameID != "" {
		gameID = sql.NullString{String: params.GameID, Valid: true}
	}

	if params.Delay != nil && u.BroadcasterType != "partner" {
		mock_errors.WriteBadRequest(w, "Delay is partner only")
		return
	}

	var delay int
	if params.Delay == nil {
		delay = u.Delay
	} else {
		delay = *params.Delay
	}
	err = db.NewQuery(r, 100).UpdateChannel(broadcasterID, database.User{
		ID:         broadcasterID,
		Title:      params.Title,
		Language:   params.BroadcasterLanguage,
		CategoryID: gameID,
		Delay:      delay,
	})
	if err != nil {
		if database.DatabaseErrorIs(err, sqlite3.ErrConstraintForeignKey) {
			mock_errors.WriteBadRequest(w, "Game ID is invalid")
			return
		}
		mock_errors.WriteServerError(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func convertUsers(users []database.User) []Channel {
	response := []Channel{}
	for _, u := range users {
		response = append(response, Channel{
			ID:           u.ID,
			UserLogin:    u.UserLogin,
			DisplayName:  u.DisplayName,
			Title:        u.Title,
			Language:     u.Language,
			CategoryID:   u.CategoryID.String,
			CategoryName: u.CategoryName.String,
			Delay:        u.Delay,
		})
	}
	return response
}
