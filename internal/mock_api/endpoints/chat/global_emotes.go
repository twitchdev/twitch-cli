// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package chat

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var globalEmotesMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var globalEmotesScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type GlobalEmotes struct{}

func (e GlobalEmotes) Path() string { return "/chat/emotes/global" }

func (e GlobalEmotes) GetRequiredScopes(method string) []string {
	return globalEmotesScopesByMethod[method]
}

func (e GlobalEmotes) ValidMethod(method string) bool {
	return globalEmotesMethodsSupported[method]
}

func (e GlobalEmotes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getGlobalEmotes(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getGlobalEmotes(w http.ResponseWriter, r *http.Request) {
	emotes := []EmotesResponse{}

	for i := 0; i < 100; i++ {
		id := util.RandomInt(10 * 1000)
		name := util.RandomGUID()
		emotes = append(emotes, EmotesResponse{
			ID:   fmt.Sprintf("%v", id),
			Name: name,
			Images: EmotesImages{
				ImageURL1X: fmt.Sprintf("https://static-cdn.jtvnw.net/emoticons/v1/%v/1.0", id),
				ImageURL2X: fmt.Sprintf("https://static-cdn.jtvnw.net/emoticons/v1/%v/2.0", id),
				ImageURL4X: fmt.Sprintf("https://static-cdn.jtvnw.net/emoticons/v1/%v/4.0", id),
			},
		})
	}

	bytes, _ := json.Marshal(models.APIResponse{Data: emotes})
	w.Write(bytes)
}
