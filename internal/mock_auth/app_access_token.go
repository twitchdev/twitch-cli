// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package mock_auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/util"
)

type AppAccessTokenEndpoint struct{}

type AppAccessTokenEndpointResposne struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int      `json:"expires_in"`
	Scope        []string `json:"scope"`
	TokenType    string   `json:"token_type"`
}

func (e AppAccessTokenEndpoint) Path() string { return "/token" }

func (e AppAccessTokenEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	clientID := r.URL.Query().Get("client_id")
	clientSecret := r.URL.Query().Get("client_secret")
	grantType := r.URL.Query().Get("grant_type")
	scope := r.URL.Query().Get("scope")
	scopes := strings.Split(scope, " ")
	if clientID == "" || clientSecret == "" || grantType != "client_credentials" {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if areValidScopes(scopes, APP_ACCES_TOKEN) != true {
		w.Write(mock_errors.GetErrorBytes(http.StatusBadRequest, errors.New("Unauthorized"), "Invalid scopes requested"))
		return
	}

	res, err := db.NewQuery(r, 10).GetAuthenticationClient(database.AuthenticationClient{ID: clientID, Secret: clientSecret})
	if err != nil {
		w.Write(mock_errors.GetErrorBytes(http.StatusInternalServerError, err, err.Error()))
	}

	ac := res.Data.([]database.AuthenticationClient)
	if len(ac) == 0 {
		w.Write(mock_errors.GetErrorBytes(http.StatusBadRequest, errors.New("Unauthorized"), "Client ID/Secret invalid"))
		return
	}

	a := database.Authorization{
		ClientID:  ac[0].ID,
		UserID:    "",
		ExpiresAt: util.GetTimestamp().Add(24 * time.Hour).Format(time.RFC3339),
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
