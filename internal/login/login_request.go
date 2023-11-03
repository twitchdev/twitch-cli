// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package login

import (
	"io"
	"net/http"
	"time"

	"github.com/twitchdev/twitch-cli/internal/request"
)

type loginRequestResponse struct {
	StatusCode int
	Body       []byte
}

type loginHeader struct {
	Key   string
	Value string
}

func loginRequest(method string, url string, payload io.Reader) (loginRequestResponse, error) {
	return loginRequestWithHeaders(method, url, payload, []loginHeader{})
}

func loginRequestWithHeaders(method string, url string, payload io.Reader, headers []loginHeader) (loginRequestResponse, error) {
	req, err := request.NewRequest(method, url, payload)

	if err != nil {
		return loginRequestResponse{}, err
	}

	for _, header := range headers {
		req.Header.Add(header.Key, header.Value)
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return loginRequestResponse{}, err
	}

	body, err := io.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		return loginRequestResponse{}, err
	}

	return loginRequestResponse{
		StatusCode: resp.StatusCode,
		Body:       body,
	}, nil
}
