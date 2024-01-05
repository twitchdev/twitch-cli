// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package util

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var legacySubFolder = ".twitch-cli"
var subFolder = "twitch-cli"

// GetApplicationDir returns a string representation of the home path for use with configuration/data storage needs
func GetApplicationDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// check if the home/.twitch-cli folder exists; if so, use that as the path
	if _, err := os.Stat(filepath.Join(home, ".twitch-cli")); !os.IsNotExist(err) {
		return filepath.Join(home, ".twitch-cli"), nil
	}

	// handles the XDG_CONFIG_HOME var as well as using AppData
	configPath, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(configPath, subFolder)

	// if the full path doesn't exist, make all the folders to get there
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	// if the user ends up in this state, provide some basic diagnostic info
	if path == "" {
		triageMessage := fmt.Sprintf("Invalid path generated; Please file a bugreport here: https://github.com/twitchdev/twitch-cli/issues/new\nInclude this in the report:\n-----\nOS: %v\nArchitecture: %v\nVersion: %v\n-----", runtime.GOOS, runtime.GOARCH, GetVersion())
		return "", errors.New(triageMessage)
	}

	return path, nil
}

// GetConfigPath returns a string representation of the configuration's path
func GetConfigPath() (string, error) {
	home, err := GetApplicationDir()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(home, ".twitch-cli.env")

	// purely for testing purposes- this allows us to run tests without overwriting the user's config
	if os.Getenv("GOLANG_TESTING") == "true" {
		configPath = filepath.Join(home, ".twitch-cli-test.env")
	}

	return configPath, nil
}
