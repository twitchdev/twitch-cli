// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package login

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/spf13/viper"
)

type LoginParameters struct {
	ClientID     string
	ClientSecret string
	Scopes       string
	Token        string
}

type RefreshParameters struct {
	ClientID     string
	ClientSecret string
	RefreshToken string
}

type RefreshTokenResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    float64  `json:"expires_in"`
	Scope        []string `json:"scope"`
	TokenType    string   `json:"token_type"`
}

type ClientCredentialsResponse struct {
	AccessToken string  `json:"access_token"`
	ExpiresIn   float64 `json:"expires_in"`
	TokenType   string  `json:"token_type"`
}

type AuthorizationCodeResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    float64  `json:"expires_in"`
	Scope        []string `json:"scope"`
	TokenType    string   `json:"token_type"`
}

type UserAuthorizationQueryResponse struct {
	Code  string
	State string
}

var redirectURI = "http://localhost:3000"

func ClientCredentialsLogin(p LoginParameters) {
	twitchClientCredentialsURL := fmt.Sprintf(`https://id.twitch.tv/oauth2/token?grant_type=client_credentials&client_id=%s&client_secret=%s`, p.ClientID, p.ClientSecret)

	resp, err := loginRequest(http.MethodPost, twitchClientCredentialsURL, nil)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	var r ClientCredentialsResponse
	if err := json.Unmarshal(resp.Body, &r); err != nil {
		println(err.Error())
		return
	}

	expiresAt := time.Now().Add(time.Duration(int64(time.Second) * int64(r.ExpiresIn)))
	println("App Access Token: ", r.AccessToken)
	var scopes []string
	storeInConfig(r.AccessToken, "", scopes, expiresAt)
	return
}

func UserCredentialsLogin(p LoginParameters) {
	twitchAuthorizeURL := fmt.Sprintf(`https://id.twitch.tv/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s&force_verify=true`, p.ClientID, redirectURI)

	if p.Scopes != "" {
		twitchAuthorizeURL += "&scope=" + p.Scopes
	}
	state, err := generateState()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	twitchAuthorizeURL += "&state=" + state

	openBrowser(twitchAuthorizeURL)

	ur, err := userAuthServer()
	if err != nil {
		println(err.Error())
		return
	}

	if ur.State != state {
		println("state mismatch")
		return
	}

	twitchUserTokenURL := fmt.Sprintf(`https://id.twitch.tv/oauth2/token?grant_type=authorization_code&client_id=%s&client_secret=%s&redirect_uri=%s&code=%s`, p.ClientID, p.ClientSecret, redirectURI, ur.Code)
	resp, err := loginRequest(http.MethodPost, twitchUserTokenURL, nil)
	if err != nil {
		fmt.Printf("Error reading body: %v", err)
		return
	}

	var r AuthorizationCodeResponse
	if err := json.Unmarshal(resp.Body, &r); err != nil {
		println(err.Error())
		return
	}

	expiresAt := time.Now().Add(time.Duration(int64(time.Second) * int64(r.ExpiresIn)))
	println(fmt.Sprintf("User Access Token: %s\nRefresh Token: %s\nExpires At: %s\nScopes: %s", r.AccessToken, r.RefreshToken, expiresAt, r.Scope))
	storeInConfig(r.AccessToken, r.RefreshToken, r.Scope, expiresAt)
	return
}

func CredentialsLogout(p LoginParameters) {
	twitchClientCredentialsURL := fmt.Sprintf(`https://id.twitch.tv/oauth2/revoke?client_id=%s&token=%s`, p.ClientID, p.Token)

	resp, err := loginRequest(http.MethodPost, twitchClientCredentialsURL, nil)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	if resp.StatusCode != 200 {
		println("API responded with an error:")
		println(string(resp.Body))
	} else {
		println("Token '" + p.Token + "' has been successfully revoked.")
	}
}

func RefreshUserToken(p RefreshParameters) (string, error) {
	twitchRefreshTokenURL := fmt.Sprintf(`https://id.twitch.tv/oauth2/token?grant_type=refresh_token&client_id=%s&client_secret=%s&redirect_uri=&refresh_token=%s`, p.ClientID, p.ClientSecret, p.RefreshToken)
	resp, err := loginRequest(http.MethodPost, twitchRefreshTokenURL, nil)
	if err != nil {
		return "", err
	}
	var r RefreshTokenResponse

	if err := json.Unmarshal(resp.Body, &r); err != nil {
		return "", nil
	}
	expiresAt := time.Now().Add(time.Duration(int64(time.Second) * int64(r.ExpiresIn)))
	storeInConfig(r.AccessToken, r.RefreshToken, r.Scope, expiresAt)
	return r.AccessToken, nil
}

func generateState() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func userAuthServer() (UserAuthorizationQueryResponse, error) {
	m := http.NewServeMux()
	s := http.Server{Addr: ":3000", Handler: m}
	userAuth := make(chan UserAuthorizationQueryResponse)
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Feel free to close this browser window."))

		var u = UserAuthorizationQueryResponse{
			Code:  r.URL.Query().Get("code"),
			State: r.URL.Query().Get("state"),
		}
		userAuth <- u
	})
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
			return
		}
	}()

	userAuthResponse := <-userAuth

	s.Shutdown(context.Background())
	return userAuthResponse, nil
}

func storeInConfig(token string, refresh string, scopes []string, expiresAt time.Time) {
	viper.Set("accessToken", token)
	viper.Set("refreshToken", refresh)
	viper.Set("tokenScopes", scopes)
	viper.Set("tokenExpiration", expiresAt.Format(time.RFC3339))

	viper.WriteConfig()
}
