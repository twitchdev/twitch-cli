// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package util

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mitchellh/go-homedir"
)

var legacySubFolder = ".twitch-cli"
var subFolder = "twitch-cli"

// GetApplicationDir returns a string representation of the home path for use with configuration/data storage needs
func GetApplicationDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	// check if the home/.twitch-cli folder exists; if so, use that as the path
	if _, err := os.Stat(filepath.Join(home, ".twitch-cli")); !os.IsNotExist(err) {
		return filepath.Join(home, ".twitch-cli"), nil
	}

	path := ""

	xdg, exists := os.LookupEnv("XDG_CONFIG_HOME") // per comment in PR #71- using this env var if present
	if !exists || xdg == "" {
		// if not present, set sane defaults- APPDATA\twitch-cli for Windows, .config/twitch-cli for OSX/Linux
		if runtime.GOOS == "WINDOWS" {
			path = filepath.Join("$APPDATA", subFolder)
		} else {
			path = filepath.Join(home, ".config", subFolder)
		}
	} else {
		// if it does exist, then just use it and combine with the subfolder; example is: $HOME/.config/twitch-cli
		path = filepath.Join(xdg, subFolder)
	}

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

	return configPath, nil
}
