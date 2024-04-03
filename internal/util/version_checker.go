// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/viper"

	vrs "github.com/hashicorp/go-version"
)

type updateCheckerReleasesResponse struct {
	// Only thing we really care about here is the tag name
	TagName string `json:"tag_name"`
}

// Check Github's Releases API to see if we're running the latest version.
// If there's any errors, quietly allow it to fail.
func CheckForUpdatesAndPrintNotice() {
	// Don't bother running this if current application is built from source
	if strings.EqualFold(GetVersion(), "source") {
		return
	}

	// Don't run if program is running in CI/CD (CI env variable will be set), or if TWITCH_DISABLE_UPDATE_CHECKS is set to true.
	if viper.GetBool("disable_update_checks") || strings.EqualFold(os.Getenv("CI"), "true") {
		return
	}

	// Don't run if this already ran successfully today
	today, _ := time.Parse(time.DateOnly, time.Now().Format(time.DateOnly))
	lastRunDate, _ := time.Parse(time.DateOnly, viper.GetString("last_update_check"))
	if !today.After(lastRunDate) {
		return
	}

	runningLatestVersion, latestVersionTag, err := areWeRunningLatestVersion()
	if err != nil {
		return // Drop errors without notifying
	}

	if !runningLatestVersion {
		// Messages to be displayed
		messages := []string{
			" A new Twitch CLI release is available: " + latestVersionTag + " ",
			"",
			" See upgrade instructions at: ",
			" https://github.com/twitchdev/twitch-cli/blob/main/README.md#update ",
			"",
			" Check out the release notes at: ",
			" https://github.com/twitchdev/twitch-cli/releases/latest ",
		}

		// Find longest message
		longestMessageLength := 0
		for _, str := range messages {
			if len(str) > longestMessageLength {
				longestMessageLength = len(str)
			}
		}

		// Print messages
		shaded := color.New(color.BgWhite, color.FgBlack).SprintfFunc()

		fmt.Println()
		for _, str := range messages {
			fmt.Printf(" %v\n", shaded(str+strings.Repeat(" ", longestMessageLength-len(str))))
		}
		fmt.Println()

		// Update config so this isn't repeated until tomorrow
		viper.Set("last_update_check", GetTimestamp().Format(time.DateOnly))
		configPath, err := GetConfigPath()
		if err != nil {
			return
		}
		_ = viper.WriteConfigAs(configPath)
	}
}

// Makes the call to Github's Releases API to check for the latest release version, and compares it to the current version.
func areWeRunningLatestVersion() (bool, string, error) {
	REGEX_TAG_NAME_VERSION := regexp.MustCompile("^(?:v)(.+)") // Removes the v from the start of the tag, if it exists

	client := &http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest("GET", "https://api.github.com/repos/twitchdev/twitch-cli/releases/latest", nil)
	if err != nil {
		return false, "", err
	}
	req.Header.Set("User-Agent", "twitch-cli/"+GetVersion())

	response, err := client.Do(req)
	if err != nil {
		return false, "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return false, "", err
	}

	var obj updateCheckerReleasesResponse
	err = json.Unmarshal(body, &obj)
	if err != nil {
		return false, "", err
	}

	latestReleaseVersion, err := vrs.NewVersion(REGEX_TAG_NAME_VERSION.FindAllStringSubmatch(obj.TagName, -1)[0][1])
	if err != nil {
		return false, "", err
	}

	currentVersion, err := vrs.NewVersion(GetVersion())
	if err != nil {
		return false, "", err
	}

	return currentVersion.GreaterThanOrEqual(latestReleaseVersion), obj.TagName, nil
}
