// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package login

import (
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/twitchdev/twitch-cli/internal/request"
)

type loginRequestResponse struct {
	StatusCode int
	Body       []byte
}

func loginRequest(method string, url string, payload io.Reader) (loginRequestResponse, error) {
	req, err := request.NewRequest(method, url, payload)

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return loginRequestResponse{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		return loginRequestResponse{}, err
	}

	return loginRequestResponse{
		StatusCode: resp.StatusCode,
		Body:       body,
	}, nil
}
