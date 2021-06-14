// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package authentication

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/util"
	"github.com/twitchdev/twitch-cli/test_setup"
)

var a *assert.Assertions
var ac = database.AuthenticationClient{ID: "1234", Secret: "1234", Name: "test_client", IsExtension: false}
var token = "potato"
var firstRun = true

func TestHasScope(t *testing.T) {
	a = test_setup.SetupTestEnv(t)

	u := UserAuthentication{UserID: "1", Scopes: []string{"user:read:email"}}

	a.Equal(true, u.HasScope("user:read:email"))
	a.Equal(false, u.HasScope("user:read"))

	a.Equal(true, u.HasOneOfRequiredScope([]string{}))
	a.Equal(true, u.HasOneOfRequiredScope([]string{"user:read:email", "user:read"}))
	a.Equal(false, u.HasOneOfRequiredScope([]string{"user:read"}))

}

func TestNatchesBroadcasterIDParam(t *testing.T) {
	a = test_setup.SetupTestEnv(t)

	req, _ := http.NewRequest(http.MethodGet, "http://google.com", nil)
	u := UserAuthentication{UserID: "1", Scopes: []string{"user:read:email"}}

	q := req.URL.Query()
	q.Set("broadcaster_id", "2")
	req.URL.RawQuery = q.Encode()

	a.Equal(false, u.MatchesBroadcasterIDParam(req))

	q.Set("broadcaster_id", "1")
	req.URL.RawQuery = q.Encode()

	a.Equal(true, u.MatchesBroadcasterIDParam(req))
}

func TestAuthenticationMiddleware(t *testing.T) {
	a = test_setup.SetupTestEnv(t)
	ts := httptest.NewServer(baseMiddleware(AuthenticationMiddleware(testEndpoint{})))

	req, _ := http.NewRequest(http.MethodGet, ts.URL+testEndpoint{}.Path(), nil)

	resp, err := http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(401, resp.StatusCode)

	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(401, resp.StatusCode)

	req.Header.Set("Client-ID", ac.ID)
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(401, resp.StatusCode)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer%v", token))
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(401, resp.StatusCode)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer"))
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(401, resp.StatusCode)
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
			a.Nil(err)
			auth, err := db.NewQuery(r, 100).CreateAuthorization(database.Authorization{ClientID: ac.ID, UserID: "1", Scopes: "user:read:email bits:read", Token: token, ExpiresAt: util.GetTimestamp().Add(7 * 24 * time.Hour).Format(time.RFC3339)})
			token = auth.Token
			a.Nil(err)

			firstRun = false
		}

		defer db.DB.Close()

		ctx = context.WithValue(ctx, "db", db)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

type testEndpoint struct{}

func (e testEndpoint) Path() string { return "/endpoint" }

func (e testEndpoint) GetRequiredScopes(method string) []string {
	return []string{}
}

func (e testEndpoint) ValidMethod(method string) bool {
	return true
}

func (e testEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(UserAuthentication)

	a.NotNil(userCtx)
	w.WriteHeader(200)
}
