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

var channelEmotesMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var channelEmotesScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type ChannelEmotes struct{}

func (e ChannelEmotes) Path() string { return "/chat/emotes/channel" }

func (e ChannelEmotes) GetRequiredScopes(method string) []string {
	return channelEmotesScopesByMethod[method]
}

func (e ChannelEmotes) ValidMethod(method string) bool {
	return channelEmotesMethodsSupported[method]
}

func (e ChannelEmotes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getChannelEmotes(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
func getChannelEmotes(w http.ResponseWriter, r *http.Request) {
	emotes := []EmotesResponse{}
	broadcaster := r.URL.Query().Get("broadcaster_id")
	if broadcaster == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter broadcaster_id")
		return
	}

	setID := fmt.Sprint(util.RandomInt(10 * 1000))
	ownerID := util.RandomUserID()
	for _, v := range defaultEmoteTypes {
		emoteType := v
		for i := 0; i < 5; i++ {
			id := util.RandomInt(10 * 1000)
			name := util.RandomGUID()
			er := EmotesResponse{
				ID:   fmt.Sprint(id),
				Name: name,
				Images: EmotesImages{
					ImageURL1X: fmt.Sprintf("https://static-cdn.jtvnw.net/emoticons/v1/%v/1.0", id),
					ImageURL2X: fmt.Sprintf("https://static-cdn.jtvnw.net/emoticons/v1/%v/2.0", id),
					ImageURL4X: fmt.Sprintf("https://static-cdn.jtvnw.net/emoticons/v1/%v/4.0", id),
				},
				EmoteType:  &emoteType,
				EmoteSetID: &setID,
				OwnerID:    &ownerID,
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

	bytes, _ := json.Marshal(models.APIResponse{Data: emotes})
	w.Write(bytes)
}
