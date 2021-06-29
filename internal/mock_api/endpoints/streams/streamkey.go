// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package streams

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

var streamKeyMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var streamKeyScopesByMethod = map[string][]string{
	http.MethodGet:    {"channel:read:stream_key"},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type StreamKey struct{}

type StreamKeyResposne struct {
	StreamKey string `json:"stream_key"`
}

func (e StreamKey) Path() string { return "/streams/key" }

func (e StreamKey) GetRequiredScopes(method string) []string {
	return streamKeyScopesByMethod[method]
}

func (e StreamKey) ValidMethod(method string) bool {
	return streamKeyMethodsSupported[method]
}

func (e StreamKey) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getStreamKey(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getStreamKey(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)

	if !userCtx.MatchesBroadcasterIDParam(r) {
		mock_errors.WriteUnauthorized(w, "Broadcaster ID does not match token.")
		return
	}

	streamKeys := []StreamKeyResposne{{StreamKey: fmt.Sprintf("live_%v_%v", userCtx.UserID, util.RandomGUID())}}

	bytes, _ := json.Marshal(models.APIResponse{Data: streamKeys})
	w.Write(bytes)
}
