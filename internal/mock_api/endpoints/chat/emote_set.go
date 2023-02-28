// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package chat

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var emoteSetsMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var emoteSetsScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type EmoteSets struct{}

func (e EmoteSets) Path() string { return "/chat/emotes/set" }

func (e EmoteSets) GetRequiredScopes(method string) []string {
	return emoteSetsScopesByMethod[method]
}

func (e EmoteSets) ValidMethod(method string) bool {
	return emoteSetsMethodsSupported[method]
}

func (e EmoteSets) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getEmoteSets(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
func getEmoteSets(w http.ResponseWriter, r *http.Request) {
	emotes := []EmotesResponse{}
	setID := r.URL.Query().Get("emote_set_id")
	if setID == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter emote_set_id")
		return
	}

	for _, v := range defaultEmoteTypes {
		ownerID := util.RandomUserID()
		emoteType := v
		for i := 0; i < 5; i++ {
			id := util.RandomInt(10 * 1000)
			name := util.RandomGUID()
			er := EmotesResponse{
				ID:   fmt.Sprint(id),
				Name: name,
				Images: EmotesImages{
					ImageURL1X: fmt.Sprintf("https://static-cdn.jtvnw.net/emoticons/v2/%v/static/light/1.0", id),
					ImageURL2X: fmt.Sprintf("https://static-cdn.jtvnw.net/emoticons/v2/%v/static/light/2.0", id),
					ImageURL4X: fmt.Sprintf("https://static-cdn.jtvnw.net/emoticons/v2/%v/static/light/3.0", id),
				},
				EmoteType:  &emoteType,
				EmoteSetID: &setID,
				OwnerID:    &ownerID,
				Format: []string{
					"static",
					"animated",
				},
				Scale: []string{
					"1.0",
					"2.0",
					"3.0",
				},
				ThemeMode: []string{
					"light",
					"dark",
				},
			}
			if emoteType == "subscription" {
				thousand := "1000"
				er.Tier = &thousand
			} else {
				es := ""
				er.Tier = &es
			}

			emotes = append(emotes, er)
		}
	}

	bytes, _ := json.Marshal(
		models.APIResponse{
			Data:     emotes,
			Template: templateEmoteURL,
		},
	)
	w.Write(bytes)
}
