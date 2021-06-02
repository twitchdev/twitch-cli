// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package mock_units

import (
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/mock_units/categories"
	"github.com/twitchdev/twitch-cli/internal/mock_units/streams"
	"github.com/twitchdev/twitch-cli/internal/mock_units/subscriptions"
	"github.com/twitchdev/twitch-cli/internal/mock_units/tags"
	"github.com/twitchdev/twitch-cli/internal/mock_units/teams"
	"github.com/twitchdev/twitch-cli/internal/mock_units/users"
	"github.com/twitchdev/twitch-cli/internal/mock_units/videos"
)

// MockEndpoint is an implementation of an endpoint in the API; this enables the quick building of new endpoints with minimal additional logic
type UnitEndpoint interface {
	Path() string
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func All() []UnitEndpoint {
	return []UnitEndpoint{
		categories.Endpoint{},
		users.Endpoint{},
		teams.Endpoint{},
		videos.Endpoint{},
		streams.Endpoint{},
		tags.Endpoint{},
		subscriptions.Endpoint{},
	}
}
