// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package generate

import (
	"fmt"

	"github.com/twitchdev/twitch-cli/internal/util"
)

var usernamePossibilities = []string{
	"KomodoHype",
	"Dave",
	"Concrete",
	"Entree",
	"Steve",
	"Chief",
	"Drake",
	"Lara",
	"Shepard",
	"Kid",
	"Jack",
	"Gordon",
	"Fisher",
	"Banjo",
	"Marcus",
	"Duke",
	"Marston",
	"Developer",
	"Skateboard",
	"Egg",
	"Lion",
}

func generateUsername() string {
	return fmt.Sprintf("%v%v%v",
		usernamePossibilities[util.RandomInt(int64(len(usernamePossibilities)))],
		usernamePossibilities[util.RandomInt(int64(len(usernamePossibilities)))],
		util.RandomInt(1000),
	)
}
