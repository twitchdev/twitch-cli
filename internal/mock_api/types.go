// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package mock_api

import "net/http"

// MockEndpoint is an implementation of an endpoint in the API; this enables the quick building of new endpoints with minimal additional logic
type MockEndpoint interface {
	Path() string
	GetRequiredScopes(method string) []string
	ValidMethod(string) bool
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
