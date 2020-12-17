// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package cmd

import (
	"github.com/twitchdev/twitch-cli/internal/login"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var isUserToken bool
var userScopes string

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "token",
	Short: "Logs into Twitch and returns an access token according to your client id/secret in the configuration.",
	Run:   loginCmdRun,
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().BoolVarP(&isUserToken, "user-token", "u", false, "Whether to login as a user or getting an app access token.")
	loginCmd.Flags().StringVarP(&userScopes, "scopes", "s", "", "Space seperated list of scopes to request with your user token.")
}

func loginCmdRun(cmd *cobra.Command, args []string) {
	clientID = viper.GetString("clientId")
	clientSecret = viper.GetString("clientSecret")
	if clientID == "" || clientSecret == "" {
		println("No Client ID or Secret found in configuration. Triggering configuration now.")
		configureCmd.Run(cmd, args)
		clientID = viper.GetString("clientId")
		clientSecret = viper.GetString("clientSecret")
	}

	var p = login.LoginParameters{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       userScopes,
	}

	if isUserToken == true {
		login.UserCredentialsLogin(p)
	} else {
		login.ClientCredentialsLogin(p)
	}
}
