// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package mock_auth

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/util"
)

type UserTokenEndpoint struct{}

func (e UserTokenEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	clientID := r.URL.Query().Get("client_id")
	clientSecret := r.URL.Query().Get("client_secret")
	grantType := r.URL.Query().Get("grant_type")
	userID := r.URL.Query().Get("user_id")
	scope := r.URL.Query().Get("scope")
	scopes := strings.Split(scope, " ")

	if clientID == "" || clientSecret == "" || grantType != "user_token" || userID == "" {
		mock_errors.WriteBadRequest(w, "missing required parameter")
		return
	}

	if areValidScopes(scopes, USER_ACCESS_TOKEN) != true {
		mock_errors.WriteBadRequest(w, "Invalid scopes requested")
		return
	}

	res, err := db.NewQuery(r, 10).GetAuthenticationClient(database.AuthenticationClient{ID: clientID, Secret: clientSecret})
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}

	ac := res.Data.([]database.AuthenticationClient)
	if len(ac) == 0 {
		mock_errors.WriteBadRequest(w, "Client ID/Secret invalid")
		return
	}

	res, err = db.NewQuery(r, 10).GetUsers(database.User{ID: userID})
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}

	users := res.Data.([]database.User)
	if len(users) == 0 {
		mock_errors.WriteBadRequest(w, "User ID invalid")
		return
	}

	a := database.Authorization{
		ClientID:  ac[0].ID,
		UserID:    userID,
		ExpiresAt: util.GetTimestamp().Add(24 * time.Hour).Format(time.RFC3339),
		Scopes:    strings.Join(scopes, " "),
	}

	auth, err := db.NewQuery(r, 100).CreateAuthorization(a)
	if err != nil {
		w.Write(mock_errors.GetErrorBytes(http.StatusInternalServerError, err, err.Error()))
		return
	}
	ea, _ := time.Parse(time.RFC3339, a.ExpiresAt)
	ater := AppAccessTokenEndpointResposne{
		AccessToken:  auth.Token,
		RefreshToken: "",
		ExpiresIn:    int(ea.Sub(time.Now().UTC()).Seconds()),
		Scope:        scopes,
		TokenType:    "bearer",
	}
	bytes, _ := json.Marshal(ater)
	w.Write(bytes)
	return
}

func (e UserTokenEndpoint) Path() string { return "/authorize" }
