// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package login

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/twitchdev/twitch-cli/internal/util"
)

type LoginParameters struct {
	ClientID     string
	ClientSecret string
	Scopes       string
	Token        string
	URL          string
	RedirectURL  string
	AuthorizeURL string
}

type RefreshParameters struct {
	ClientID     string
	ClientSecret string
	RefreshToken string
	URL          string
}

type AuthorizationResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int64    `json:"expires_in"`
	Scope        []string `json:"scope"`
	TokenType    string   `json:"token_type"`
}

type UserAuthorizationQueryResponse struct {
	Code  string
	State string
	Error error
}

type LoginResponse struct {
	Response  AuthorizationResponse
	ExpiresAt time.Time
}

type ValidateResponse struct {
	ClientID  string   `json:"client_id"`
	UserLogin string   `json:"login"`
	UserID    string   `json:"user_id"`
	Scopes    []string `json:"scopes"`
	ExpiresIn int64    `json:"expires_in"`
}

const ClientCredentialsURL = "https://id.twitch.tv/oauth2/token?grant_type=client_credentials"
const UserCredentialsURL = "https://id.twitch.tv/oauth2/token?grant_type=authorization_code"

const UserAuthorizeURL = "https://id.twitch.tv/oauth2/authorize?response_type=code"

const RefreshTokenURL = "https://id.twitch.tv/oauth2/token?grant_type=refresh_token"
const RevokeTokenURL = "https://id.twitch.tv/oauth2/revoke"
const ValidateTokenURL = "https://id.twitch.tv/oauth2/validate"

// Sends `https://id.twitch.tv/oauth2/token?grant_type=client_credentials`.
// Generates a new App Access Token. Stores new token information in the CLI's config.
func ClientCredentialsLogin(p LoginParameters) (LoginResponse, error) {
	u, err := url.Parse(p.URL)
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	q.Set("client_id", p.ClientID)
	q.Set("client_secret", p.ClientSecret)
	u.RawQuery = q.Encode()

	resp, err := loginRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return LoginResponse{}, errors.New("API responded with an error while revoking token: " + string(resp.Body))
	}

	r, err := handleLoginResponse(resp.Body, true)
	if err != nil {
		return LoginResponse{}, err
	}

	return r, nil
}

// Sends `https://id.twitch.tv/oauth2/token?grant_type=authorization_code`.
// Generates a new App Access Token, requiring the use of a web browser. Stores new token information in the CLI's config.
func UserCredentialsLogin(p LoginParameters, webserverIP string, webserverPort string) (LoginResponse, error) {
	u, err := url.Parse(p.AuthorizeURL)
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	q.Set("client_id", p.ClientID)
	q.Set("redirect_uri", p.RedirectURL)
	if p.Scopes != "" {
		q.Set("scope", p.Scopes)
	}

	state, err := generateState()
	if err != nil {
		log.Fatal(err.Error())
	}

	q.Set("state", state)
	u.RawQuery = q.Encode()

	execOpenBrowser := func() {
		fmt.Println("Opening browser. Press Ctrl+C to cancel...")
		err = openBrowser(u.String())
		if err != nil {
			fmt.Printf("Unable to open default browser. You can manually navigate to this URL to complete the login: %s\n", u.String())
		}
	}

	urp, err := userAuthServer(webserverIP, webserverPort, execOpenBrowser)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("Error processing request: %v", err.Error())
	}
	ur := *urp

	if ur.State != state {
		return LoginResponse{}, fmt.Errorf("Error processing request: state mismatch")
	}

	u2, err := url.Parse(p.URL)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("Internal error: %v", err.Error())
	}

	q = u2.Query()
	q.Set("client_id", p.ClientID)
	q.Set("client_secret", p.ClientSecret)
	q.Set("redirect_uri", p.RedirectURL)
	q.Set("code", ur.Code)
	u2.RawQuery = q.Encode()

	resp, err := loginRequest(http.MethodPost, u2.String(), nil)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("Error reading body: %v", err.Error())
	}

	r, err := handleLoginResponse(resp.Body, true)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("Error handling login: %v", err.Error())
	}

	return r, nil
}

// Sends `https://id.twitch.tv/oauth2/revoke`.
// Revokes the provided token. Does not change the CLI's config at all.
func CredentialsLogout(p LoginParameters) (LoginResponse, error) {
	u, err := url.Parse(p.URL)
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	q.Set("client_id", p.ClientID)
	q.Set("token", p.Token)
	u.RawQuery = q.Encode()

	resp, err := loginRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		log.Print(err.Error())
		return LoginResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("API responded with an error while revoking token: [%v] %v", resp.StatusCode, string(resp.Body))
		return LoginResponse{}, errors.New("API responded with an error while revoking token")
	}

	log.Printf("Token %s has been successfully revoked.", p.Token)
	return LoginResponse{}, nil
}

// Sends `POST https://id.twitch.tv/oauth2/token`.
// Refreshes the provided token and optionally stores the result in the CLI's config.
func RefreshUserToken(p RefreshParameters, shouldStoreInConfig bool) (LoginResponse, error) {
	u, err := url.Parse(p.URL)
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	q.Set("client_id", p.ClientID)
	q.Set("client_secret", p.ClientSecret)
	q.Set("refresh_token", p.RefreshToken)
	u.RawQuery = q.Encode()

	resp, err := loginRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return LoginResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return LoginResponse{}, errors.New(fmt.Sprintf("Error with client while refreshing: [%v - `%v`]", resp.StatusCode, strings.TrimSpace(string(resp.Body))))
	}

	r, err := handleLoginResponse(resp.Body, shouldStoreInConfig)
	if err != nil {
		log.Printf("Error handling login: %v", err)
		return LoginResponse{}, err
	}

	return r, nil
}

// Sends `GET https://id.twitch.tv/oauth2/validate`.
// Only validates. Does not store this information in the CLI's config.
func ValidateCredentials(p LoginParameters) (ValidateResponse, error) {
	u, err := url.Parse(p.URL)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := loginRequestWithHeaders(http.MethodGet, u.String(), nil, []loginHeader{
		{
			Key:   "Authorization",
			Value: "OAuth " + p.Token,
		},
	})
	if err != nil {
		return ValidateResponse{}, err
	}

	// Handle validate response body
	var r ValidateResponse
	if err = json.Unmarshal(resp.Body, &r); err != nil {
		return ValidateResponse{}, err
	}

	return r, nil
}

func handleLoginResponse(body []byte, shouldStoreInConfig bool) (LoginResponse, error) {
	var r AuthorizationResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return LoginResponse{}, err
	}
	expiresAt := util.GetTimestamp().Add(time.Duration(int64(time.Second) * int64(r.ExpiresIn)))

	if shouldStoreInConfig {
		storeInConfig(r.AccessToken, r.RefreshToken, r.Scope, expiresAt)
	}

	return LoginResponse{
		Response:  r,
		ExpiresAt: expiresAt,
	}, nil
}

func generateState() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func openBrowser(url string) error {
	const rundllParameters = "url.dll,FileProtocolHandler"
	var err error
	switch runtime.GOOS {
	case "linux":
		if util.IsWsl() {
			err = exec.Command("rundll32.exe", rundllParameters, url).Start()
		} else {
			err = exec.Command("xdg-open", url).Start()
		}
	case "windows":
		err = exec.Command("rundll32", rundllParameters, url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}

func userAuthServer(ip string, port string, onSuccessfulListenCallback func()) (*UserAuthorizationQueryResponse, error) {
	m := http.NewServeMux()
	s := http.Server{Addr: fmt.Sprintf("%v:%v", ip, port), Handler: m}
	userAuth := make(chan UserAuthorizationQueryResponse)
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/favicon.ico" {
			return
		}
		authError := r.URL.Query().Get("error")

		if authError != "" {
			w.Write([]byte(fmt.Sprintf("Error! %v\nError Details: %v", authError, r.URL.Query().Get("error_description"))))
			var u = UserAuthorizationQueryResponse{
				Error: fmt.Errorf("%v", r.URL.Query().Get("error_description")),
			}
			userAuth <- u
		} else {
			w.Write([]byte("Feel free to close this browser window."))

			var u = UserAuthorizationQueryResponse{
				Code:  r.URL.Query().Get("code"),
				State: r.URL.Query().Get("state"),
			}
			userAuth <- u
		}
	})

	ln, err := net.Listen("tcp", s.Addr)
	defer s.Shutdown(context.Background())
	if err != nil {
		return nil, err
	}

	if onSuccessfulListenCallback != nil {
		onSuccessfulListenCallback()
	}

	go func() {
		if err := s.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
			return
		}
	}()

	log.Printf("Waiting for authorization response ...")
	userAuthResponse := <-userAuth
	log.Printf("Closing local server ...")
	return &userAuthResponse, userAuthResponse.Error
}

func storeInConfig(token string, refresh string, scopes []string, expiresAt time.Time) {
	viper.Set("accessToken", token)
	viper.Set("refreshToken", refresh)
	viper.Set("tokenScopes", scopes)
	viper.Set("tokenExpiration", expiresAt.Format(time.RFC3339Nano))

	err := viper.WriteConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		err = viper.SafeWriteConfig()
	}

	if err != nil {
		log.Fatalf("Error writing configuration: %s", err)
	}
}
