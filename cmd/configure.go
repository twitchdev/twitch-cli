// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/twitchdev/twitch-cli/internal/util"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var clientID string
var clientSecret string

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configures your Twitch CLI with your Client ID and Secret",
	Run:   configureCmdRun,
}

func init() {
	rootCmd.AddCommand(configureCmd)
}

func configureCmdRun(cmd *cobra.Command, args []string) {
	clientIDPrompt := promptui.Prompt{
		Label: "Client ID",
		Validate: func(s string) error {
			if len(s) == 30 || len(s) == 31 {
				return nil
			}
			return errors.New("Invalid length for Client ID")
		},
	}

	clientID, err := clientIDPrompt.Run()

	clientSecretPrompt := promptui.Prompt{
		Label: "Client Secret",
		Validate: func(s string) error {
			if len(s) == 30 || len(s) == 31 {
				return nil
			}
			return errors.New("Invalid length for Client Secret")
		},
	}

	clientSecret, err := clientSecretPrompt.Run()

	if clientID == "" && clientSecret == "" {
		fmt.Println("Must specify either the Client ID or Secret")
		return
	}

	viper.Set("clientId", clientID)
	viper.Set("clientSecret", clientSecret)

	home, _ := util.GetApplicationDir()
	configPath := filepath.Join(home, ".twitch-cli.env")

	_, err = os.Stat(configPath)

	if !os.IsExist(err) {
		if _, err := os.Create(configPath); err != nil {
			fmt.Println("Failed to create configuration.", err.Error())
		}
		if err := viper.WriteConfig(); err != nil {
			fmt.Println("Failed to write configuration.", err.Error())
		}
	}

	fmt.Println("Updated configuration.")
	return
}
