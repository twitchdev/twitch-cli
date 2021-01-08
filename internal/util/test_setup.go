// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package util

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func SetupTestEnv(t *testing.T) *assert.Assertions {
	assert := assert.New(t)

	home, err := GetApplicationDir()
	assert.Nil(err)

	viper.AddConfigPath(home)
	viper.SetConfigName(".twitch-cli-test")
	viper.SetConfigType("env")

	viper.Set("DB_FILENAME", "test-eventCache.db")
	return assert
}
