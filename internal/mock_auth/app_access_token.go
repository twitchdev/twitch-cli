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

type AppAccessTokenEndpoint struct{}

type AppAccessTokenRequestBody struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	Scope        string `json:"scope"`
}

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

	params := AppAccessTokenRequestBody{
		ClientID:     r.URL.Query().Get("client_id"),
		ClientSecret: r.URL.Query().Get("client_secret"),
		GrantType:    r.URL.Query().Get("grant_type"),
		Scope:        r.URL.Query().Get("scope"),
	}

	if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		err := r.ParseForm()
		if err != nil {
			mock_errors.WriteServerError(w, err.Error())
			return
		}

		if r.Form.Get("client_id") != "" {
			params.ClientID = r.Form.Get("client_id")
		}
		if r.Form.Get("client_secret") != "" {
			params.ClientSecret = r.Form.Get("client_secret")
		}
		if r.Form.Get("grant_type") != "" {
			params.GrantType = r.Form.Get("grant_type")
		}
		if r.Form.Get("scope") != "" {
			params.Scope = r.Form.Get("scope")
		}
	}

	scopes := strings.Split(params.Scope, " ")

	if params.ClientID == "" || params.ClientSecret == "" || params.GrantType != "client_credentials" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// remove empty entries
	for i, s := range scopes {
		if s == "" {
			scopes = removeFromSlice(scopes, i)
		}
	}

	if !areValidScopes(scopes, APP_ACCES_TOKEN) {
		mock_errors.WriteBadRequest(w, "Invalid scopes requested")
		return
	}

	res, err := db.NewQuery(r, 10).GetAuthenticationClient(database.AuthenticationClient{ID: params.ClientID, Secret: params.ClientSecret})
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}

	ac := res.Data.([]database.AuthenticationClient)
	if len(ac) == 0 {
		mock_errors.WriteBadRequest(w, "Client ID/Secret invalid")
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

func removeFromSlice(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}
