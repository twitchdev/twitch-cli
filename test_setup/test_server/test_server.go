// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package test_server

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
)

func SetupTestServer(next mock_api.MockEndpoint) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		// just stub it all
		db, err := database.NewConnection(true)
		if err != nil {
			log.Fatalf("Error connecting to database: %v", err.Error())
			return
		}

		defer db.DB.Close()

		ctx = context.WithValue(ctx, "db", db)
		ctx = context.WithValue(ctx, "auth", authentication.UserAuthentication{Scopes: []string{
			"analytics:read:extensions",
			"analytics:read:games",
			"bits:read",
			"channel:edit:commercial",
			"channel:manage:broadcast",
			"channel:manage:extensions",
			"channel:manage:polls",
			"channel:manage:predictions",
			"channel:manage:redemptions",
			"channel:manage:videos",
			"channel:read:editors",
			"channel:read:hype_train",
			"channel:read:polls",
			"channel:read:predictions",
			"channel:read:redemptions",
			"channel:read:stream_key",
			"channel:read:subscriptions",
			"clips:edit",
			"moderation:read",
			"moderator:manage:automod",
			"user:edit",
			"user:edit:follows",
			"user:read:email",
			"user:manage:blocked_users",
			"user:read:blocked_users",
			"user:read:broadcast",
			"user:read:follows",
			"user:read:subscriptions",
		}, UserID: "1", ClientID: "1"})
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}))
	return ts
}
