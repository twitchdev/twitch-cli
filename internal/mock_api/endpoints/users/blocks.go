// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package users

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
)

var blocksMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: true,
	http.MethodPatch:  false,
	http.MethodPut:    true,
}

var blocksScopesByMethod = map[string][]string{
	http.MethodGet:    {"user:read:blocked_users", "user:manage:blocked_users"},
	http.MethodPost:   {},
	http.MethodDelete: {"user:manage:blocked_users"},
	http.MethodPatch:  {},
	http.MethodPut:    {"user:manage:blocked_users"},
}

type Blocks struct{}

func (e Blocks) Path() string { return "/users/blocks" }

func (e Blocks) GetRequiredScopes(method string) []string {
	return blocksScopesByMethod[method]
}

func (e Blocks) ValidMethod(method string) bool {
	return blocksMethodsSupported[method]
}

func (e Blocks) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getBlocks(w, r)
	case http.MethodPut:
		putBlocks(w, r)
	case http.MethodDelete:
		deleteBlocks(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func getBlocks(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)

	if !userCtx.MatchesBroadcasterIDParam(r) {
		mock_errors.WriteUnauthorized(w, "broadcaster_id must match the token")
		return
	}

	dbr, err := db.NewQuery(r, 100).GetBlocks(database.UserRequestParams{BroadcasterID: userCtx.UserID})
	if err != nil {
		log.Print(err)
		mock_errors.WriteServerError(w, "error fetching blocks")
		return
	}

	apiResponse := models.APIResponse{
		Data:       dbr.Data,
		Pagination: &models.APIPagination{},
	}

	if len(dbr.Data.([]database.Block)) == 0 {
		apiResponse.Data = []database.Block{}
	}

	if len(dbr.Data.([]database.Block)) == dbr.Limit {
		apiResponse.Pagination.Cursor = dbr.Cursor
	}

	bytes, _ := json.Marshal(apiResponse)
	w.Write(bytes)
}

func putBlocks(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)

	target := r.URL.Query().Get("target_user_id")
	if target == "" {
		mock_errors.WriteBadRequest(w, "target_user_id is required")
		return
	}

	if target == userCtx.UserID {
		mock_errors.WriteBadRequest(w, "you can't block yourself")
		return
	}

	// other parameters are not surfaced anywhere, so they won'tb be inserted, but will validate them to match the API
	validContext := []string{"chat", "whisper"}
	sourceContext := r.URL.Query().Get("source_context")

	validReason := []string{"spam", "harassment", "other"}
	reason := r.URL.Query().Get("reason")

	if sourceContext != "" && !isOneOf(validContext, sourceContext) {
		mock_errors.WriteBadRequest(w, fmt.Sprintf("source_context must be one of %v", strings.Join(validContext, " or ")))
		return
	}

	if reason != "" && !isOneOf(validReason, reason) {
		mock_errors.WriteBadRequest(w, fmt.Sprintf("reason must be one of %v", strings.Join(validReason, " or ")))
		return
	}

	err := db.NewQuery(r, 100).AddBlock(database.UserRequestParams{BroadcasterID: userCtx.UserID, UserID: target})
	if err != nil {
		mock_errors.WriteServerError(w, "error adding block")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteBlocks(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)

	target := r.URL.Query().Get("target_user_id")
	if target == "" {
		mock_errors.WriteBadRequest(w, "target_user_id is required")
		return
	}

	err := db.NewQuery(r, 100).DeleteBlock(target, userCtx.UserID)
	if err != nil {
		mock_errors.WriteServerError(w, "error deleting block")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func isOneOf(valid []string, value string) bool {
	for _, v := range valid {
		if value == v {
			return true
		}
	}
	return false
}
