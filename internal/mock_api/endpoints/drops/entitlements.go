// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package drops

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var dropsEntitlementsMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  true,
	http.MethodPut:    false,
}

var dropsEntitlementsScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type DropsEntitlements struct{}

type PatchEntitlementsBody struct {
	FulfillmentStatus string   `json:"fulfillment_status"`
	EntitlementIDs    []string `json:"entitlement_ids"`
}

type PatchEntitlementsResponse struct {
	Status string   `json:"status"`
	IDs    []string `json:"ids"`
}

func (e DropsEntitlements) Path() string { return "/entitlements/drops" }

func (e DropsEntitlements) GetRequiredScopes(method string) []string {
	return dropsEntitlementsScopesByMethod[method]
}

func (e DropsEntitlements) ValidMethod(method string) bool {
	return dropsEntitlementsMethodsSupported[method]
}

func (e DropsEntitlements) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getEntitlements(w, r)
	case http.MethodPatch:
		patchEntitlements(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getEntitlements(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	id := r.URL.Query().Get("id")
	userID := r.URL.Query().Get("user_id")
	gameID := r.URL.Query().Get("game_id")
	status := r.URL.Query().Get("fulfillment_status")

	if userCtx.UserID != "" && userID != "" {
		mock_errors.WriteBadRequest(w, "user_id is invalid when using user access token")
		return
	}
	if userCtx.UserID != "" {
		userID = userCtx.UserID
	}
	e := database.DropsEntitlement{UserID: userID, GameID: gameID, ID: id, Status: status}
	dbr, err := db.NewQuery(r, 1000).GetDropsEntitlements(e)
	if err != nil {
		mock_errors.WriteServerError(w, "error fetching entitlements")
		return
	}
	entitlements := dbr.Data.([]database.DropsEntitlement)
	if len(entitlements) == 0 {
		entitlements = []database.DropsEntitlement{}
	}
	apiResponse := models.APIResponse{
		Data: entitlements,
	}
	if len(entitlements) == dbr.Limit {
		apiResponse.Pagination = &models.APIPagination{
			Cursor: dbr.Cursor,
		}
	}

	bytes, err := json.Marshal(apiResponse)
	w.Write(bytes)
}

func patchEntitlements(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)

	var body PatchEntitlementsBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		mock_errors.WriteBadRequest(w, "Invalid body")
		return
	}

	if body.FulfillmentStatus != "CLAIMED" && body.FulfillmentStatus != "FULFILLED" {
		mock_errors.WriteBadRequest(w, "fulfillment_status must be one of CLAIMED or FULFILLED")
		return
	}

	if len(body.EntitlementIDs) == 0 || len(body.EntitlementIDs) > 100 {
		mock_errors.WriteBadRequest(w, "entitlement_ids must be at least 1 and at most 100")
		return
	}
	s := PatchEntitlementsResponse{
		Status: "SUCCESS",
	}
	ua := PatchEntitlementsResponse{Status: "UNAUTHORIZED"}
	fail := PatchEntitlementsResponse{Status: "UPDATE_FAILED"}
	notFound := PatchEntitlementsResponse{Status: "NOT_FOUND"}
	for _, e := range body.EntitlementIDs {
		dbr, err := db.NewQuery(nil, 100).GetDropsEntitlements(database.DropsEntitlement{ID: e})
		if err != nil {
			fail.IDs = append(fail.IDs, e)
			continue
		}
		entitlement := dbr.Data.([]database.DropsEntitlement)
		if len(entitlement) == 0 {
			notFound.IDs = append(notFound.IDs, e)
			continue
		}

		if userCtx.UserID != "" && userCtx.UserID != entitlement[0].UserID {
			ua.IDs = append(ua.IDs, e)
			continue
		}

		err = db.NewQuery(nil, 100).UpdateDropsEntitlement(
			database.DropsEntitlement{
				ID:          e,
				UserID:      entitlement[0].UserID,
				Status:      body.FulfillmentStatus,
				LastUpdated: util.GetTimestamp().Format(time.RFC3339Nano),
			},
		)
		if err != nil {
			fail.IDs = append(fail.IDs, e)
			continue
		}
		s.IDs = append(s.IDs, e)
	}
	all := []PatchEntitlementsResponse{
		s,
		ua,
		fail,
		notFound,
	}
	resp := []PatchEntitlementsResponse{}
	for _, r := range all {
		if len(r.IDs) != 0 {
			resp = append(resp, r)
		}
	}

	apiResponse := models.APIResponse{
		Data: resp,
	}

	bytes, _ := json.Marshal(apiResponse)
	w.Write(bytes)
}
