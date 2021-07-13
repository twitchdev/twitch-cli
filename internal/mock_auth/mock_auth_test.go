// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package mock_auth

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/util"
	"github.com/twitchdev/twitch-cli/test_setup"
)

var a *assert.Assertions
var firstRun = true
var ac = database.AuthenticationClient{ID: "222", Secret: "333", Name: "test_client", IsExtension: false}

func TestAreValidScopes(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	a.Equal(true, areValidScopes([]string{"user:read:email"}, USER_ACCESS_TOKEN))
	a.Equal(false, areValidScopes([]string{"user:read:email"}, APP_ACCES_TOKEN))
}

func TestUserToken(t *testing.T) {
	a = test_setup.SetupTestEnv(t)
	ts := httptest.NewServer(baseMiddleware(UserTokenEndpoint{}))

	req, _ := http.NewRequest(http.MethodPost, ts.URL+UserTokenEndpoint{}.Path(), nil)
	q := req.URL.Query()

	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(400, resp.StatusCode)

	// valid values
	q.Set("client_id", ac.ID)
	q.Set("client_secret", ac.Secret)
	q.Set("grant_type", "user_token")
	q.Set("user_id", "1")

	q.Set("scope", "potato")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	q.Set("scope", "")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	q.Set("client_id", "1234")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	q.Set("client_id", ac.ID)
	q.Set("user_id", util.RandomGUID())
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)
}

func TestAppAccessToken(t *testing.T) {
	a = test_setup.SetupTestEnv(t)
	ts := httptest.NewServer(baseMiddleware(AppAccessTokenEndpoint{}))

	req, _ := http.NewRequest(http.MethodPost, ts.URL+AppAccessTokenEndpoint{}.Path(), nil)
	q := req.URL.Query()

	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(400, resp.StatusCode)

	// valid values
	q.Set("client_id", ac.ID)
	q.Set("client_secret", ac.Secret)
	q.Set("grant_type", "client_credentials")

	q.Set("scope", "potato")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	q.Set("scope", "")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)
}
func baseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		// just stub it all
		db, err := database.NewConnection()
		if err != nil {
			log.Fatalf("Error connecting to database: %v", err.Error())
			return
		}
		if firstRun == true {
			ac, err = db.NewQuery(r, 100).InsertOrUpdateAuthenticationClient(ac, false)
			a.Nil(err, err)

			firstRun = false
		}

		defer db.DB.Close()

		ctx = context.WithValue(ctx, "db", db)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
