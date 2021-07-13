// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package util

import (
	"os"
	"path/filepath"

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

	legacyFolder := false

	// check if the home/.twitch-cli folder exists; if so, use that as the path
	if _, err := os.Stat(filepath.Join(home, ".twitch-cli")); !os.IsNotExist(err) {
		legacyFolder = true
	}

	path := filepath.Join(home, legacySubFolder)

	if !legacyFolder {
		path = filepath.Join(home, ".config", subFolder) // issue #33- putting into a subfolder to avoid clutter
		subpath := filepath.Join(home, ".config")
		if _, err := os.Stat(subpath); os.IsNotExist(err) {
			os.Mkdir(subpath, 0700)
		}
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0700)
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
