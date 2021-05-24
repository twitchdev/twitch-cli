// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package endpoints

import (
	"github.com/twitchdev/twitch-cli/internal/mock_api"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/follows"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/users"
)

func All() []mock_api.MockEndpoint {
	return []mock_api.MockEndpoint{
		users.Endpoint{},
		follows.Endpoint{},
	}
}
