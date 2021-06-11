// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package login

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/twitchdev/twitch-cli/internal/util"
	"github.com/twitchdev/twitch-cli/test_setup"
)

var params = LoginParameters{
	ClientID:     "1234",
	ClientSecret: "4567",
	RedirectURL:  "https://localhost:3000",
	Token:        "890",
	Scopes:       "scope1 scope2",
}

var response = LoginResponse{
	ExpiresAt: util.GetTimestamp().Add(10 * time.Minute),
	Response: AuthorizationResponse{
		TokenType:    "bearer",
		AccessToken:  "890",
		ExpiresIn:    int64(10 * time.Minute / time.Second),
		Scope:        []string{},
		RefreshToken: "890",
	},
}

func TestClientCredentialsLogin(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		qp := r.URL.Query()

		a.Equal(params.ClientID, qp.Get("client_id"), "ClientID mismatch")
		a.Equal(params.ClientSecret, qp.Get("client_secret"), "Secret mismatch")

		body, err := json.Marshal(response.Response)
		a.Nil(err)

		w.Write(body)
	}))
	defer ts.Close()

	params.URL = ts.URL + "?test=test"

	ClientCredentialsLogin(params)
}

func TestCredentialsLogout(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		qp := r.URL.Query()

		a.Equal(params.ClientID, qp.Get("client_id"), "ClientID mismatch")
		a.Equal(params.Token, qp.Get("token"), "Token mismatch")

		if r.URL.Path == "/fail" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer ts.Close()

	params.URL = ts.URL

	r, e := CredentialsLogout(params)
	a.Nil(e)
	a.Empty(r.Response.AccessToken)

	params.URL = ts.URL + "/fail"

	r, e = CredentialsLogout(params)
	a.NotNil(e, "Returned a nil error, expected error %v", e)
}

func TestGenerateState(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	state, err := generateState()

	a.Nil(err)
	a.NotNil(state)
}

func TestStoreInConfig(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	r := response.Response
	storeInConfig(r.AccessToken, r.RefreshToken, r.Scope, response.ExpiresAt)

	a.Equal(r.AccessToken, viper.Get("accesstoken"), "Invalid token in config.")
	a.Equal(r.RefreshToken, viper.Get("refreshtoken"), "Invalid refresh token in config.")
	a.Equal(r.Scope, viper.Get("tokenscopes"), "Invalid scopes in config.")
	a.Equal(response.ExpiresAt.Format(time.RFC3339Nano), viper.GetString("tokenexpiration"), "Invalid expiration in config.")
}

func TestRefreshUserToken(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		qp := r.URL.Query()

		a.Equal(params.ClientID, qp.Get("client_id"), "ClientID mismatch")
		a.Equal(params.ClientSecret, qp.Get("client_secret"), "Secret mismatch")

		body, err := json.Marshal(response.Response)
		a.Nil(err)

		w.Write(body)
	}))
	defer ts.Close()

	resp, err := RefreshUserToken(RefreshParameters{
		RefreshToken: params.Token,
		URL:          ts.URL + "?foo=bar",
		ClientID:     params.ClientID,
		ClientSecret: params.ClientSecret,
	})

	a.Nil(err)
	a.NotNil(resp)

	a.Equal(response.Response.AccessToken, resp.Response.AccessToken, "Access token mismatch")
}

func TestUserAuthServer(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	state, err := generateState()
	a.Nil(err)
	a.NotNil(state)

	code := "1234"

	userResponse := make(chan UserAuthorizationQueryResponse)

	go func() {
		res, err := userAuthServer()
		a.Nil(err)
		userResponse <- res
	}()

	_, err = loginRequest(http.MethodGet, fmt.Sprintf("http://localhost:3000?code=%s&state=%s", code, state), nil)
	a.Nil(err)

	ur := <-userResponse
	a.Equal(state, ur.State, "State mismatch")
	a.Equal(code, ur.Code, "Code mismatch")
}
