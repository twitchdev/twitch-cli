// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package authentication

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
)

type UserAuthentication struct {
	Scopes   []string
	UserID   string
	ClientID string
}

func AuthenticationMiddleware(next mock_api.MockEndpoint) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := r.Context().Value("db").(database.CLIDatabase)

		// skip auth check for unsupported methods
		if next.ValidMethod(r.Method) == false {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if len(r.URL.Query()["skip_auth"]) > 0 && r.URL.Query()["skip_auth"][0] == "true" {
			fakeAuth := UserAuthentication{}
			r = r.WithContext(context.WithValue(r.Context(), "auth", fakeAuth))
			next.ServeHTTP(w, r)
			log.Printf("Skipping auth...")
			return
		}

		clientID := r.Header.Get("Client-ID")
		bearerToken := r.Header.Get("Authorization")
		unauthroizedError := mock_errors.GetErrorBytes(http.StatusUnauthorized, errors.New("Unauthorized"), "Missing Client ID or OAuth token")
		if clientID == "" || bearerToken == "" || len(bearerToken) < 7 {
			w.Write(unauthroizedError)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		prefix := strings.ToLower(bearerToken[:6])
		token := bearerToken[7:]

		// check if the client ID is invalid or missing the proper token prefix
		if len(clientID) < 30 || prefix != "bearer" {
			w.Write(unauthroizedError)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// get the authorization from the db
		auth, err := db.NewQuery(r, 100).GetAuthorizationByToken(token)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// check if invalid
		if auth.Token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(mock_errors.GetErrorBytes(http.StatusUnauthorized, errors.New("Unauthorized"), "Invalid OAuth token"))
			return
		}

		// check for mismatches
		if auth.ClientID != clientID {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(mock_errors.GetErrorBytes(http.StatusUnauthorized, errors.New("Unauthorized"), "Missing Client ID or OAuth token"))
			return
		}

		// check if expired
		expiration, err := time.Parse(time.RFC3339, auth.ExpiresAt)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if time.Now().After(expiration) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(mock_errors.GetErrorBytes(http.StatusUnauthorized, errors.New("Unauthorized"), "Token expired"))
			return
		}

		// pass as context
		authContext := UserAuthentication{
			Scopes:   strings.Split(auth.Scopes, " "),
			UserID:   auth.UserID,
			ClientID: auth.ClientID,
		}

		if authContext.HasOneOfRequiredScope(next.GetRequiredScopes(r.Method)) == false {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(mock_errors.GetErrorBytes(http.StatusUnauthorized, errors.New("Unauthorized"), fmt.Sprintf("Missing required scope %v", next.GetRequiredScopes(r.Method))))
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "trace-id", "1234"))
		r = r.WithContext(context.WithValue(r.Context(), "auth", authContext))

		next.ServeHTTP(w, r)
	})
}

func (u UserAuthentication) HasScope(scope string) bool {
	for _, s := range u.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}

func (u UserAuthentication) HasOneOfRequiredScope(scopes []string) bool {
	if len(scopes) == 0 {
		return true
	}

	for _, s := range scopes {
		for _, us := range u.Scopes {
			if s == us {
				return true
			}
		}
	}
	return false
}

func (u *UserAuthentication) MatchesBroadcasterIDParam(r *http.Request) bool {
	bid := r.URL.Query().Get("broadcaster_id")
	if bid == "" || bid != u.UserID {
		return false
	}
	return true
}
