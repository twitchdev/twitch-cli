// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package moderation

import (
	"encoding/json"
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var automodStatusMethodsSupported = map[string]bool{
	http.MethodGet:    false,
	http.MethodPost:   true,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var automodStatusScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {"moderation:read"},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type AutomodStatus struct{}

type PostAutomodStatusBody struct {
	Data []PostAutomodStatusBodyData `json:"data"`
}
type PostAutomodStatusBodyData struct {
	MessageID   string `json:"msg_id"`
	MessageText string `json:"msg_text"`
	UserID      string `json:"user_id"`
}

type PostAutomodStatusResponse struct {
	MessageID   string `json:"msg_id"`
	IsPermitted bool   `json:"is_permitted"`
}

func (e AutomodStatus) Path() string { return "/moderation/enforcements/status" }

func (e AutomodStatus) GetRequiredScopes(method string) []string {
	return automodStatusScopesByMethod[method]
}

func (e AutomodStatus) ValidMethod(method string) bool {
	return automodStatusMethodsSupported[method]
}

func (e AutomodStatus) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodPost:
		postAutomodStatus(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func postAutomodStatus(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	var body PostAutomodStatusBody
	response := []PostAutomodStatusResponse{}
	if !userCtx.MatchesBroadcasterIDParam(r) {
		mock_errors.WriteUnauthorized(w, "broadcaster_id does not match token")
		return
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		mock_errors.WriteBadRequest(w, "error reading body")
		return
	}

	for _, data := range body.Data {
		if data.MessageID == "" || data.MessageText == "" {
			mock_errors.WriteBadRequest(w, "msg_id and msg_text are required")
			return
		}

		shouldPermit := util.RandomInt(2) == 0

		response = append(response, PostAutomodStatusResponse{MessageID: data.MessageID, IsPermitted: shouldPermit})
	}

	bytes, _ := json.Marshal(models.APIResponse{Data: response})
	w.Write(bytes)
}
