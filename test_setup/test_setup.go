// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package test_setup

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/twitchdev/twitch-cli/internal/util"
)

func SetupTestEnv(t *testing.T) *assert.Assertions {
	assert := assert.New(t)

	home, err := util.GetApplicationDir()
	assert.Nil(err)

	viper.AddConfigPath(home)
	viper.SetConfigName(".twitch-cli-test")
	viper.SetConfigType("env")

	viper.Set("DB_FILENAME", "test-eventCache.db")
	t.Setenv("GOLANG_TESTING", "true")
	return assert
}
