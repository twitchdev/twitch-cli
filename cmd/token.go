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
var revokeToken string
var overrideClientId string

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
	loginCmd.Flags().StringVarP(&revokeToken, "revoke", "r", "", "Instead of generating a new token, revoke the one passed to this parameter.")
	loginCmd.Flags().StringVar(&overrideClientId, "client-id", "", "Override/manually set client ID for token actions. By default client ID from CLI config will be used.")
}

func loginCmdRun(cmd *cobra.Command, args []string) {
	clientID = viper.GetString("clientId")
	clientSecret = viper.GetString("clientSecret")

	redirectURL := "http://localhost:3000"

	if clientID == "" || clientSecret == "" {
		println("No Client ID or Secret found in configuration. Triggering configuration now.")
		configureCmd.Run(cmd, args)
		clientID = viper.GetString("clientId")
		clientSecret = viper.GetString("clientSecret")
	}

	if overrideClientId != "" {
		clientID = overrideClientId
	}

	var p = login.LoginParameters{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       userScopes,
		RedirectURL:  redirectURL,
		AuthorizeURL: login.UserAuthorizeURL,
	}

	if revokeToken != "" {
		p.Token = revokeToken
		p.URL = login.RevokeTokenURL
		login.CredentialsLogout(p)
	} else if isUserToken == true {
		p.URL = login.UserCredentialsURL
		login.UserCredentialsLogin(p)
	} else {
		p.URL = login.ClientCredentialsURL
		login.ClientCredentialsLogin(p)
	}
}
