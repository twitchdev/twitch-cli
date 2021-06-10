// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package mock_auth

import (
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
)

type AuthEndpoint interface {
	Path() string
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

var db database.CLIDatabase

const APP_ACCES_TOKEN = "app_access"
const USER_ACCESS_TOKEN = "user_access"

var validScopesByTokenType = map[string]map[string]bool{
	APP_ACCES_TOKEN: {
		"analytics:read:extensions": true,
		"analytics:read:games":      true,
	},
	USER_ACCESS_TOKEN: {
		"analytics:read:extensions":  true,
		"analytics:read:games":       true,
		"bits:read":                  true,
		"channel:edit:commercial":    true,
		"channel:manage:broadcast":   true,
		"channel:manage:extensions":  true,
		"channel:manage:polls":       true,
		"channel:manage:predictions": true,
		"channel:manage:redemptions": true,
		"channel:manage:videos":      true,
		"channel:read:editors":       true,
		"channel:read:hype_train":    true,
		"channel:read:polls":         true,
		"channel:read:predictions":   true,
		"channel:read:redemptions":   true,
		"channel:read:stream_key":    true,
		"channel:read:subscriptions": true,
		"clips:edit":                 true,
		"moderation:read":            true,
		"moderator:manage:automod":   true,
		"user:edit":                  true,
		"user:edit:follows":          true,
		"user:read:email":            true,
		"user:manage:blocked_users":  true,
		"user:read:blocked_users":    true,
		"user:read:broadcast":        true,
		"user:read:follows":          true,
		"user:read:subscriptions":    true,
	},
}

func All() []AuthEndpoint {
	return []AuthEndpoint{
		UserTokenEndpoint{},
		AppAccessTokenEndpoint{},
	}
}

func areValidScopes(scopes []string, tokenType string) bool {
	if tokenType != APP_ACCES_TOKEN && tokenType != USER_ACCESS_TOKEN {
		return false
	}
	if len(scopes) == 0 {
		return true
	}

	for _, s := range scopes {
		if s == "" {
			continue
		}
		if validScopesByTokenType[tokenType][s] != true {
			return false
		}
	}
	return true
}
