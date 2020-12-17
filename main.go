// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package main

import (
	"github.com/twitchdev/twitch-cli/cmd"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var buildVersion string

func main() {
	if len(buildVersion) > 0 {
		util.SetVersion(buildVersion)
	}
	cmd.Execute()
}
