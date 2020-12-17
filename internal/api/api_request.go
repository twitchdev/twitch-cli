// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/twitchdev/twitch-cli/internal/request"
)

type apiRequestParameters struct {
	Token    string
	ClientID string
}
type apiRequestResponse struct {
	StatusCode int
	Body       []byte
}

func apiRequest(method string, url string, payload []byte, p apiRequestParameters) (apiRequestResponse, error) {
	req, err := request.NewRequest(method, url, bytes.NewBuffer(payload))

	req.Header.Set("Client-ID", p.ClientID)
	req.Header.Set("Content-Type", "application/json")

	if p.Token != "" {
		req.Header.Set("Authorization", "Bearer "+p.Token)
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error reading body: %v", err)
		return apiRequestResponse{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		fmt.Printf("Error reading body: %v", err)
		return apiRequestResponse{}, err
	}

	return apiRequestResponse{
		StatusCode: resp.StatusCode,
		Body:       body,
	}, nil
}
