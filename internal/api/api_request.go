// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/twitchdev/twitch-cli/internal/request"
	"golang.org/x/time/rate"
)

type apiRequestParameters struct {
	Token    string
	ClientID string
}
type apiRequestResponse struct {
	StatusCode int
	Body       []byte

	HttpMethod      string
	RequestPath     string
	RequestHeaders  http.Header
	ResponseHeaders http.Header
	HttpVersion     string
}

func apiRequest(method string, url string, payload []byte, p apiRequestParameters) (apiRequestResponse, error) {
	req, err := request.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return apiRequestResponse{}, err
	}
	req.Header.Set("Client-ID", p.ClientID)
	req.Header.Set("Content-Type", "application/json")
	rl := rate.NewLimiter(rate.Every(time.Minute), 800)

	client := NewClient(rl)

	if p.Token != "" {
		req.Header.Set("Authorization", "Bearer "+p.Token)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error reading body: %v", err)
		return apiRequestResponse{}, err
	}

	body, err := io.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		fmt.Printf("Error reading body: %v", err)
		return apiRequestResponse{}, err
	}

	return apiRequestResponse{
		StatusCode: resp.StatusCode,
		Body:       body,

		HttpMethod:      req.Method,
		RequestPath:     req.URL.RequestURI(),
		RequestHeaders:  req.Header,
		ResponseHeaders: resp.Header,
		HttpVersion:     req.Proto,
	}, nil
}
