// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package mock_auth

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/twitchdev/twitch-cli/internal/database"
)

type ValidateTokenEndpoint struct{}

type ValidateTokenEndpointResponse struct {
	ClientID  string   `json:"client_id"`
	UserID    string   `json:"user_id,omitempty"`
	UserLogin string   `json:"login,omitempty"`
	ExpiresIn int      `json:"expires_in"`
	Scopes    []string `json:"scopes"`
}

func (e ValidateTokenEndpoint) Path() string { return "/validate" }

func (e ValidateTokenEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	tokenHeader := r.Header.Get("Authorization")
	if tokenHeader == "" || len(tokenHeader) < 7 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token := ""
	// handle prefixes
	h := strings.ToLower(tokenHeader)
	if strings.HasPrefix(h, "oauth ") {
		token = tokenHeader[6:]
	} else if strings.HasPrefix(h, "bearer ") {
		token = tokenHeader[7:]
	}
	println(token)
	auth, err := db.NewQuery(r, 100).GetAuthorizationByToken(token)
	if err != nil || auth.ID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expiresAt, _ := time.Parse(time.RFC3339, auth.ExpiresAt)

	diff := expiresAt.Sub(time.Now())

	scopes := []string{}
	for _, s := range strings.Split(auth.Scopes, " ") {
		if s != "" {
			scopes = append(scopes, s)
		}
	}
	resp := ValidateTokenEndpointResponse{
		ClientID:  auth.ClientID,
		UserID:    auth.UserID,
		Scopes:    scopes,
		ExpiresIn: int(diff.Seconds()),
	}
	if auth.UserID != "" {
		user, err := db.NewQuery(r, 100).GetUser(database.User{ID: auth.UserID})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp.UserLogin = user.UserLogin
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
