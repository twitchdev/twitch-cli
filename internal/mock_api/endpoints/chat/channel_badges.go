// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package chat

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var channelBadgesMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var channelBadgesScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type ChannelBadges struct{}

func (e ChannelBadges) Path() string { return "/chat/badges" }

var defaultChannelBadgeVersions = []string{"0", "12", "2", "3", "6", "9"}

func (e ChannelBadges) GetRequiredScopes(method string) []string {
	return channelBadgesScopesByMethod[method]
}

func (e ChannelBadges) ValidMethod(method string) bool {
	return channelBadgesMethodsSupported[method]
}

func (e ChannelBadges) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getChannelBadges(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
func getChannelBadges(w http.ResponseWriter, r *http.Request) {
	badges := []BadgesResponse{}
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)

	if !userCtx.MatchesBroadcasterIDParam(r) {
		mock_errors.WriteUnauthorized(w, "Broadcaster ID does not match token.")
		return
	}

	b := BadgesResponse{
		SetID:    "subscriber",
		Versions: []BadgesVersion{},
	}

	for _, v := range defaultChannelBadgeVersions {
		id := util.RandomGUID()
		b.Versions = append(b.Versions, BadgesVersion{
			ID:         v,
			ImageURL1X: fmt.Sprintf("https://static-cdn.jtvnw.net/badges/v1/%v/1", id),
			ImageURL2X: fmt.Sprintf("https://static-cdn.jtvnw.net/badges/v1/%v/2", id),
			ImageURL4X: fmt.Sprintf("https://static-cdn.jtvnw.net/badges/v1/%v/3", id),
		})
	}
	badges = append(badges, b)

	bytes, _ := json.Marshal(models.APIResponse{Data: badges})
	w.Write(bytes)
}
