// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package cmd

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/twitchdev/twitch-cli/internal/login"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var isUserToken bool
var userScopes string
var revokeToken string
var validateToken string
var refreshToken string
var overrideClientId string
var overrideClientSecret string
var tokenServerPort int
var tokenServerIP string
var redirectHost string

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
	loginCmd.Flags().StringVarP(&refreshToken, "refresh", "R", "", "Instead of generating a new token, refresh the token associated with the Refresh Token passed to this parameter.")
	loginCmd.Flags().StringVar(&overrideClientId, "client-id", "", "Override/manually set Client ID for token actions. By default Client ID from CLI config will be used.")
	loginCmd.Flags().StringVar(&overrideClientSecret, "secret", "", "Override/manually set Client Secret for token actions. By default Client Secret from CLI config will be used.")
	loginCmd.Flags().StringVar(&tokenServerIP, "ip", "", "Manually set the IP address to be bound to for the User Token web server.")
	loginCmd.Flags().IntVarP(&tokenServerPort, "port", "p", 3000, "Manually set the port to be used for the User Token web server.")
	loginCmd.Flags().StringVar(&redirectHost, "redirect-host", "localhost", "Manually set the host to be used for the redirect URL")
}

func loginCmdRun(cmd *cobra.Command, args []string) error {
	clientID = viper.GetString("clientId")
	clientSecret = viper.GetString("clientSecret")

	webserverPort := strconv.Itoa(tokenServerPort)
	redirectURL := fmt.Sprintf("http://%v:%v", redirectHost, webserverPort)

	if clientID == "" || clientSecret == "" {
		println("No Client ID or Secret found in configuration. Triggering configuration now.")
		err := configureCmd.RunE(cmd, args)
		if err != nil {
			return err
		}

		clientID = viper.GetString("clientId")
		clientSecret = viper.GetString("clientSecret")
	}

	if overrideClientId != "" {
		clientID = overrideClientId
	}

	if overrideClientSecret != "" {
		clientSecret = overrideClientSecret
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
		_, err := login.CredentialsLogout(p)

		if err != nil {
			return err
		}

		log.Printf("Token %s has been successfully revoked", p.Token)

	} else if validateToken != "" {
		p.Token = validateToken
		p.URL = login.ValidateTokenURL
		r, err := login.ValidateCredentials(p)
		if err != nil {
			return err
		}

		tokenType := "App Access Token"
		if r.UserID != "" {
			tokenType = "User Access Token"
		}

		expiresInTimestamp := time.Now().Add(time.Duration(r.ExpiresIn) * time.Second).UTC().Format(time.RFC1123)

		lightYellow := color.New(color.FgHiYellow).PrintfFunc()
		white := color.New(color.FgWhite).SprintfFunc()

		lightYellow("Client ID: %v\n", white(r.ClientID))
		lightYellow("Token Type: %v\n", white(tokenType))
		if r.UserID != "" {
			lightYellow("User ID: %v\n", white(r.UserID))
			lightYellow("User Login: %v\n", white(r.UserLogin))
		}
		lightYellow("Expires In: %v\n", white("%v (%v)", strconv.FormatInt(r.ExpiresIn, 10), expiresInTimestamp))

		if len(r.Scopes) == 0 {
			lightYellow("User ID: %v\n", white("None"))
		} else {
			lightYellow("Scopes:\n")
			for _, s := range r.Scopes {
				fmt.Println(white("- %v\n", s))
			}
		}

	} else if refreshToken != "" {
		p.URL = login.RefreshTokenURL

		// If we are overriding the Client ID then we shouldn't store this in the config.
		shouldStoreInConfig := (overrideClientId == "")

		resp, err := login.RefreshUserToken(login.RefreshParameters{
			RefreshToken: refreshToken,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			URL:          login.RefreshTokenURL,
		}, shouldStoreInConfig)

		if err != nil {
			errDescription := ""
			if overrideClientId == "" {
				errDescription = "Check `--refresh` flag to ensure the provided Refresh Token is valid for the Client ID set with `twitch config`."
			} else {
				errDescription = "Check `--refresh` and `--client-id` flags to ensure the provided Refresh Token is valid for the provided Client ID."
			}

			return fmt.Errorf("%v\n%v", err.Error(), errDescription)
		}

		lightYellow := color.New(color.FgHiYellow).SprintfFunc()

		log.Println("Successfully refreshed Access Token.")
		log.Println(lightYellow("Access Token: ") + resp.Response.AccessToken)
		log.Println(lightYellow("Refresh Token: ") + resp.Response.RefreshToken)
		log.Println(lightYellow("Expires At: ") + resp.ExpiresAt.String())

	} else if isUserToken {
		p.URL = login.UserCredentialsURL
		resp, err := login.UserCredentialsLogin(p, tokenServerIP, webserverPort)

		if err != nil {
			return err
		}

		lightYellow := color.New(color.FgHiYellow).SprintfFunc()

		log.Println("Successfully generated User Access Token.")
		log.Println(lightYellow("User Access Token: ") + resp.Response.AccessToken)
		log.Println(lightYellow("Refresh Token: ") + resp.Response.RefreshToken)
		log.Println(lightYellow("Expires At: ") + resp.ExpiresAt.String())
		log.Println(lightYellow("Scopes: ") + fmt.Sprintf("%v", resp.Response.Scope))

	} else {
		p.URL = login.ClientCredentialsURL
		resp, err := login.ClientCredentialsLogin(p)

		if err != nil {
			return err
		}

		lightYellow := color.New(color.FgHiYellow).SprintfFunc()

		log.Println("Successfully generated App Access Token.")
		log.Println(lightYellow("App Access Token: ") + resp.Response.AccessToken)
		log.Println(lightYellow("Expires At: ") + resp.ExpiresAt.String())
		log.Println(lightYellow("Scopes: ") + fmt.Sprintf("%v", resp.Response.Scope))
	}

	return nil
}
