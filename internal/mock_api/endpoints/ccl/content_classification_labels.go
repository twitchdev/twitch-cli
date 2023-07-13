// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package ccl

import (
	"encoding/json"
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/models"
)

var cclMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var cclScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type ContentClassificationLabels struct{}

func (e ContentClassificationLabels) Path() string { return "/content_classification_labels" }

func (e ContentClassificationLabels) GetRequiredScopes(method string) []string {
	return cclScopesByMethod[method]
}

func (e ContentClassificationLabels) ValidMethod(method string) bool {
	return cclMethodsSupported[method]
}

func (e ContentClassificationLabels) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getContentClassificationLabels(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
func getContentClassificationLabels(w http.ResponseWriter, r *http.Request) {
	// TODO: locale param

	allCCLs := []models.ContentClassificationLabel{}
	for _, ccl := range models.CCL_MAP {
		allCCLs = append(allCCLs, ccl)
	}

	bytes, _ := json.Marshal(
		models.APIResponse{
			Data: allCCLs,
		},
	)
	w.Write(bytes)
}
