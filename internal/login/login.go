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

type DeviceCodeFlowInitResponse struct {
	DeviceCode      string `json:"device_code"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
	UserCode        string `json:"user_code"`
	VerificationUri string `json:"verification_uri"`
}

const ClientCredentialsURL = "https://id.twitch.tv/oauth2/token?grant_type=client_credentials"
const UserCredentialsURL = "https://id.twitch.tv/oauth2/token?grant_type=authorization_code"

const UserAuthorizeURL = "https://id.twitch.tv/oauth2/authorize?response_type=code"

const RefreshTokenURL = "https://id.twitch.tv/oauth2/token?grant_type=refresh_token"
const RevokeTokenURL = "https://id.twitch.tv/oauth2/revoke"
const ValidateTokenURL = "https://id.twitch.tv/oauth2/validate"

const DeviceCodeFlowUrl = "https://id.twitch.tv/oauth2/device"
const DeviceCodeFlowTokenURL = "https://id.twitch.tv/oauth2/token"
const DeviceCodeFlowGrantType = "urn:ietf:params:oauth:grant-type:device_code"

// Sends `https://id.twitch.tv/oauth2/token?grant_type=client_credentials`.
// Generates a new App Access Token. Stores new token information in the CLI's config.
func ClientCredentialsLogin(p LoginParameters) (LoginResponse, error) {
	u, err := url.Parse(p.URL)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("Internal error: %v", err.Error())
	}
	q := u.Query()
	q.Set("client_id", p.ClientID)
	q.Set("client_secret", p.ClientSecret)
	u.RawQuery = q.Encode()

	resp, err := loginRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("Error processing request: %v", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return LoginResponse{}, errors.New("API responded with an error while revoking token: " + string(resp.Body))
	}

	r, err := handleLoginResponse(resp.Body, true)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("Error processing login response: %v", err.Error())
	}

	return r, nil
}

// Uses Authorization Code Flow: https://dev.twitch.tv/docs/authentication/getting-tokens-oauth/#authorization-code-grant-flow
// Sends `https://id.twitch.tv/oauth2/token?grant_type=authorization_code`.
// Generates a new User Access Token, requiring the use of a web browser. Stores new token information in the CLI's config.
func UserCredentialsLogin_AuthorizationCodeFlow(p LoginParameters, webserverIP string, webserverPort string) (LoginResponse, error) {
	u, err := url.Parse(p.AuthorizeURL)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("Internal error (parsing AuthorizeURL): %v", err.Error())
	}
	q := u.Query()
	q.Set("client_id", p.ClientID)
	q.Set("redirect_uri", p.RedirectURL)
	if p.Scopes != "" {
		q.Set("scope", p.Scopes)
	}

	state, err := generateState()
	if err != nil {
		return LoginResponse{}, fmt.Errorf("Internal error (generating state): %v", err.Error())
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
		return LoginResponse{}, fmt.Errorf("Internal error (parsing URL): %v", err.Error())
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

	if resp.StatusCode == 400 {
		// If 400 is returned, the applications' Client Type was set up as "Public", and you can only use Implicit Auth or Device Code Flow to get a User Access Token
		return LoginResponse{}, fmt.Errorf(
			"This Client Type of this Client ID is set to \"Public\", which doesn't allow the use of Authorization Code Grant Flow.\n" +
				"Please call the token command with the --dcf flag to use Device Code Flow. For example: twitch token -u --dcf",
		)
	}

	r, err := handleLoginResponse(resp.Body, true)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("Error handling login: %v", err.Error())
	}

	return r, nil
}

// Uses Device Code Flow: https://dev.twitch.tv/docs/authentication/getting-tokens-oauth/#device-code-grant-flow
// Generates a new User Access Token, requiring the use of a web browser from any device. Stores new token information in the CLI's config.
func UserCredentialsLogin_DeviceCodeFlow(p LoginParameters) (LoginResponse, error) {
	// Initiate DCF flow
	deviceResp, err := dcfInitiateRequest(DeviceCodeFlowUrl, p.ClientID, p.Scopes)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("Error initiating Device Code Flow: %v", err.Error())
	}

	var deviceObj DeviceCodeFlowInitResponse
	if err := json.Unmarshal(deviceResp.Body, &deviceObj); err != nil {
		return LoginResponse{}, fmt.Errorf("Error reading body: %v", err.Error())
	}
	expirationTime := time.Now().Add(time.Second * time.Duration(deviceObj.ExpiresIn))

	fmt.Printf("Started Device Code Flow login.\n")
	fmt.Printf("Use this URL to log in: %v\n", deviceObj.VerificationUri)
	fmt.Printf("Use this code when prompted at the above URL: %v\n\n", deviceObj.UserCode)
	fmt.Printf("This system will check every %v seconds, and will expire after %v minutes.\n", deviceObj.Interval, (deviceObj.ExpiresIn / 60))

	// Loop and check for user login. Respects given interval, and times out after expiration
	tokenResp := loginRequestResponse{StatusCode: 999}
	for tokenResp.StatusCode != 0 {
		// Check for expiration
		if time.Now().After(expirationTime) {
			return LoginResponse{}, fmt.Errorf("The Device Code used for getting access token has expired. Run token command again to generate a new user.")
		}

		// Wait interval
		time.Sleep(time.Second * time.Duration(deviceObj.Interval))

		// Check for token
		tokenResp, err = dcfTokenRequest(DeviceCodeFlowTokenURL, p.ClientID, p.Scopes, deviceObj.DeviceCode, DeviceCodeFlowGrantType)
		if err != nil {
			return LoginResponse{}, fmt.Errorf("Error getting token via Device Code Flow: %v", err)
		}

		if tokenResp.StatusCode == 200 {
			r, err := handleLoginResponse(tokenResp.Body, true)
			if err != nil {
				return LoginResponse{}, fmt.Errorf("Error handling login: %v", err.Error())
			}
			return r, nil
		}
	}

	return LoginResponse{}, nil
}

// Sends `https://id.twitch.tv/oauth2/revoke`.
// Revokes the provided token. Does not change the CLI's config at all.
func CredentialsLogout(p LoginParameters) (LoginResponse, error) {
	u, err := url.Parse(p.URL)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("Internal error (parsing URL): %v", err.Error())
	}
	q := u.Query()
	q.Set("client_id", p.ClientID)
	q.Set("token", p.Token)
	u.RawQuery = q.Encode()

	resp, err := loginRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("Error reading body: %v", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return LoginResponse{}, fmt.Errorf("API responded with an error while revoking token: [%v] %v", resp.StatusCode, string(resp.Body))
	}

	return LoginResponse{}, nil
}

// Sends `POST https://id.twitch.tv/oauth2/token`.
// Refreshes the provided token and optionally stores the result in the CLI's config.
func RefreshUserToken(p RefreshParameters, shouldStoreInConfig bool) (LoginResponse, error) {
	u, err := url.Parse(p.URL)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("Internal error (parsing URL): %v", err)
	}
	q := u.Query()
	q.Set("client_id", p.ClientID)
	q.Set("client_secret", p.ClientSecret)
	q.Set("refresh_token", p.RefreshToken)
	u.RawQuery = q.Encode()

	resp, err := loginRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("Error processing request: %v", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return LoginResponse{}, fmt.Errorf("Error with client while refreshing: [%v - `%v`]", resp.StatusCode, strings.TrimSpace(string(resp.Body)))
	}

	r, err := handleLoginResponse(resp.Body, shouldStoreInConfig)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("Error handling login: %v", err.Error())
	}

	return r, nil
}

// Sends `GET https://id.twitch.tv/oauth2/validate`.
// Only validates. Does not store this information in the CLI's config.
func ValidateCredentials(p LoginParameters) (ValidateResponse, error) {
	u, err := url.Parse(p.URL)
	if err != nil {
		return ValidateResponse{}, fmt.Errorf("Internal error (parsing URL): %v", err)
	}

	resp, err := loginRequestWithHeaders(http.MethodGet, u.String(), nil, []loginHeader{
		{
			Key:   "Authorization",
			Value: "OAuth " + p.Token,
		},
	})
	if err != nil {
		return ValidateResponse{}, fmt.Errorf("Error processing request: %v", err)
	}

	// Handle validate response body
	var r ValidateResponse
	if err = json.Unmarshal(resp.Body, &r); err != nil {
		return ValidateResponse{}, fmt.Errorf("Error handling response: %v", err)
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
