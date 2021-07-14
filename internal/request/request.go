// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package request

import (
	"io"
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/util"
)

func NewRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	version := util.GetVersion()

	req.Header.Set("User-Agent", "twitch-cli/" + version)

	return req, nil
}
