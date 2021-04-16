// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package login

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"syscall"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/twitchdev/twitch-cli/internal/util"
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
	a := util.SetupTestEnv(t)

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
	a := util.SetupTestEnv(t)

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
	a := util.SetupTestEnv(t)

	state, err := generateState()

	a.Nil(err)
	a.NotNil(state)
}

func TestStoreInConfig(t *testing.T) {
	a := util.SetupTestEnv(t)

	r := response.Response
	storeInConfig(r.AccessToken, r.RefreshToken, r.Scope, response.ExpiresAt)

	a.Equal(r.AccessToken, viper.Get("accesstoken"), "Invalid token in config.")
	a.Equal(r.RefreshToken, viper.Get("refreshtoken"), "Invalid refresh token in config.")
	a.Equal(r.Scope, viper.Get("tokenscopes"), "Invalid scopes in config.")
	a.Equal(response.ExpiresAt.Format(time.RFC3339Nano), viper.GetString("tokenexpiration"), "Invalid expiration in config.")
}

func TestRefreshUserToken(t *testing.T) {
	a := util.SetupTestEnv(t)

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
	a := util.SetupTestEnv(t)

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

func TestIsWsl(t *testing.T) {
	a := assert.New(t)

	var (
		// syscall.Utsname.Release value on various systems

		// Ubuntu 20.04 on WSL2 on Windows 10 x64 20H2
		ubuntu20Wsl2 = [65]int8{52, 46, 49, 57, 46, 49, 50, 56, 45, 109, 105, 99, 114, 111, 115, 111, 102, 116, 45, 115, 116, 97, 110, 100, 97, 114, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

		// Arch Linux on baremetal on 2021-04-02
		archReal = [65]int8{53, 46, 49, 49, 46, 49, 49, 45, 97, 114, 99, 104, 49, 45, 49, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	)

	result := isWsl(util.Syscall{
		Uname: func(buf *syscall.Utsname) (err error) {
			buf.Release = ubuntu20Wsl2
			return nil
		},
	})
	a.True(result)

	result = isWsl(util.Syscall{
		Uname: func(buf *syscall.Utsname) (err error) {
			buf.Release = archReal
			return nil
		},
	})
	a.False(result)

	result = isWsl(util.Syscall{
		Uname: func(buf *syscall.Utsname) (err error) {
			return errors.New("mocked error")
		},
	})
	a.False(result)
}
