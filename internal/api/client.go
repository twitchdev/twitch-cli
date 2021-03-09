// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package api

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type RLClient struct {
	client      *http.Client
	RateLimiter *rate.Limiter
}

func (c *RLClient) Do(req *http.Request) (*http.Response, error) {
	err := c.RateLimiter.Wait(context.Background())
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewClient(l *rate.Limiter) *RLClient {
	client := http.Client{
		Timeout: time.Second * 10,
	}

	return &RLClient{
		client:      &client,
		RateLimiter: l,
	}
}
