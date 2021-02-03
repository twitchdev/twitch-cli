// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/twitchdev/twitch-cli/internal/login"
	"github.com/twitchdev/twitch-cli/internal/util"

	"github.com/TylerBrock/colorjson"
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

var baseURL = "https://api.twitch.tv/helix"

type clientInformation struct {
	ClientID string
	Token    string
}

// NewRequest is used to request data from the Twitch API using a HTTP GET request- this function is a wrapper for the apiRequest function that handles the network call
func NewRequest(method string, path string, queryParameters []string, body []byte, prettyPrint bool) {
	client, err := getClientInformation()

	if viper.GetString("BASE_URL") != "" {
		baseURL = viper.GetString("BASE_URL")
	}

	if err != nil {
		fmt.Println("Error fetching client information", err.Error())
	}

	Parameters := url.Values{}

	if queryParameters != nil {
		path += "?"
		for _, param := range queryParameters {
			value := strings.Split(param, "=")
			Parameters.Add(value[0], value[1])
		}
		path += Parameters.Encode()
	}
	resp, err := apiRequest(strings.ToUpper(method), baseURL+path, body, apiRequestParameters{
		ClientID: client.ClientID,
		Token:    client.Token,
	})
	if err != nil {
		fmt.Printf("Error reading body: %v", err)
		return
	}

	if prettyPrint == true {
		var obj map[string]interface{}
		if err := json.Unmarshal(resp.Body, &obj); err != nil {
			fmt.Printf("Error pretty-printing body: %v", err)
			return
		}
		f := colorjson.NewFormatter()
		f.Indent = 2
		f.KeyColor = color.New(color.FgBlue).Add(color.Bold)
		s, _ := f.Marshal(obj)

		fmt.Println(string(s))
		return
	}
	fmt.Println(string(resp.Body))
	return
}

// ValidOptions returns a list of supported endpoints given a specified method as noted in the map endpointMethodSupports, which is located in resources.go of this package.
func ValidOptions(method string) []string {
	names := []string{}

	for n, m := range endpointMethodSupports {
		if m[method] {
			names = append(names, n)
		}
	}

	sort.Strings(names)

	return names
}

func getClientInformation() (clientInformation, error) {
	clientID := viper.GetString("clientID")
	expiration := viper.GetString("tokenexpiration")
	token := viper.GetString("accessToken")

	// Handle legacy nonexpiring tokens
	if expiration == "0" {
		return clientInformation{
			Token:    token,
			ClientID: clientID,
		}, nil
	}

	ex, _ := time.Parse(time.RFC3339Nano, expiration)
	if ex.Before(util.GetTimestamp()) {
		refreshToken := viper.GetString("refreshToken")

		if refreshToken == "" {
			log.Fatal("Please run twitch token")
		}

		clientSecret := viper.GetString("clientSecret")

		var err error
		r, err := login.RefreshUserToken(login.RefreshParameters{
			RefreshToken: refreshToken,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			URL:          login.RefreshTokenURL,
		})
		if err != nil {
			return clientInformation{}, err
		}
		fmt.Printf("%v", r)
		token = r.Response.AccessToken
	}

	return clientInformation{Token: token, ClientID: clientID}, nil
}
