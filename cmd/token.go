// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/twitchdev/twitch-cli/internal/login"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var isUserToken bool
var userScopes string
var revokeToken string
var validateToken string
var overrideClientId string
var tokenServerPort int
var tokenServerIP string

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "token",
	Short: "Logs into Twitch and returns an access token according to your client id/secret in the configuration.",
	RunE:  loginCmdRun,
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().BoolVarP(&isUserToken, "user-token", "u", false, "Whether to login as a user or getting an app access token.")
	loginCmd.Flags().StringVarP(&userScopes, "scopes", "s", "", "Space separated list of scopes to request with your user token.")
	loginCmd.Flags().StringVarP(&revokeToken, "revoke", "r", "", "Instead of generating a new token, revoke the one passed to this parameter.")
	loginCmd.Flags().StringVarP(&validateToken, "validate", "v", "", "Instead of generating a new token, validate the one passed to this parameter.")
	loginCmd.Flags().StringVar(&overrideClientId, "client-id", "", "Override/manually set client ID for token actions. By default client ID from CLI config will be used.")
	loginCmd.Flags().StringVar(&tokenServerIP, "ip", "localhost", "Manually set the IP address to be binded to for the User Token web server.")
	loginCmd.Flags().IntVarP(&tokenServerPort, "port", "p", 3000, "Manually set the port to be used for the User Token web server.")
}

func loginCmdRun(cmd *cobra.Command, args []string) error {
	clientID = viper.GetString("clientId")
	clientSecret = viper.GetString("clientSecret")

	webserverPort := strconv.Itoa(tokenServerPort)
	redirectURL := fmt.Sprintf("http://%v:%v", tokenServerIP, webserverPort)

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
	} else if validateToken != "" {
		p.Token = validateToken
		p.URL = login.ValidateTokenURL
		r, err := login.ValidateCredentials(p)
		if err != nil {
			return fmt.Errorf("Failed to validate: %v", err.Error())
		}

		tokenType := "App Access Token"
		if r.UserID != "" {
			tokenType = "User Access Token"
		}

		expiresInTimestamp := time.Now().Add(time.Duration(r.ExpiresIn) * time.Second).UTC().Format(time.RFC1123)

		println("Client ID: " + r.ClientID)
		println("Token Type: " + tokenType)
		if r.UserID != "" {
			println("User ID: " + r.UserID)
			println("User Login: " + r.UserLogin)
		}
		println("Expires In: " + strconv.FormatInt(r.ExpiresIn, 10) + " (" + expiresInTimestamp + ")")

		if len(r.Scopes) == 0 {
			println("Scopes: None")
		} else {
			println("Scopes:")
			for _, s := range r.Scopes {
				println("- " + s)
			}
		}
	} else if isUserToken == true {
		p.URL = login.UserCredentialsURL
		login.UserCredentialsLogin(p, tokenServerIP, webserverPort)
	} else {
		p.URL = login.ClientCredentialsURL
		login.ClientCredentialsLogin(p)
	}

	return nil
}
