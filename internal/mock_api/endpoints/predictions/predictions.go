// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package predictions

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var predictionsMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   true,
	http.MethodDelete: false,
	http.MethodPatch:  true,
	http.MethodPut:    false,
}

var predictionsScopesByMethod = map[string][]string{
	http.MethodGet:    {"channel:read:predictions", "channel:manage:predictions"},
	http.MethodPost:   {"channel:manage:predictions"},
	http.MethodDelete: {},
	http.MethodPatch:  {"channel:manage:predictions"},
	http.MethodPut:    {},
}

type Predictions struct{}

type PostPredictionsBody struct {
	BroadcasterID    string                        `json:"broadcaster_id"`
	Title            string                        `json:"title"`
	Outcomes         []PostPredictionsBodyOutcomes `json:"outcomes"`
	PredictionWindow int                           `json:"prediction_window"`
}

type PostPredictionsBodyOutcomes struct {
	Title string `json:"title"`
}

type PatchPredictionsBody struct {
	BroadcasterID    string `json:"broadcaster_id"`
	ID               string `json:"id"`
	Status           string `json:"status"`
	WinningOutcomeID string `json:"winning_outcome_id"`
}

func (e Predictions) Path() string { return "/predictions" }

func (e Predictions) GetRequiredScopes(method string) []string {
	return predictionsScopesByMethod[method]
}

func (e Predictions) ValidMethod(method string) bool {
	return predictionsMethodsSupported[method]
}

func (e Predictions) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getPredictions(w, r)
		break
	case http.MethodPost:
		postPredictions(w, r)
		break
	case http.MethodPatch:
		patchPredictions(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getPredictions(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	predictions := []database.Prediction{}
	var dbr *database.DBResponse
	var err error

	if !userCtx.MatchesBroadcasterIDParam(r) {
		mock_errors.WriteUnauthorized(w, "broadcaster_id does not match token")
		return
	}

	ids := r.URL.Query()["id"]

	for _, id := range ids {
		dbr, err := db.NewQuery(r, 100).GetPredictions(database.Prediction{ID: id, BroadcasterID: userCtx.UserID})
		if err != nil {
			mock_errors.WriteServerError(w, "error fetching predictions")
			return
		}
		predictions = append(predictions, dbr.Data.([]database.Prediction)...)
	}
	if len(ids) == 0 {
		dbr, err = db.NewQuery(r, 100).GetPredictions(database.Prediction{BroadcasterID: userCtx.UserID})
		if err != nil {
			log.Print(err)
			mock_errors.WriteServerError(w, "error fetching predictions")
			return
		}
		predictions = append(predictions, dbr.Data.([]database.Prediction)...)
	}

	apiResponse := models.APIResponse{
		Data: predictions,
	}

	if dbr != nil && dbr.Cursor != "" {
		apiResponse.Pagination = &models.APIPagination{
			Cursor: dbr.Cursor,
		}
	}

	bytes, _ := json.Marshal(apiResponse)
	w.Write(bytes)
}

func postPredictions(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	var body PostPredictionsBody

	u, err := db.NewQuery(r, 100).GetUser(database.User{ID: userCtx.UserID})
	if err != nil {
		mock_errors.WriteBadRequest(w, "error getting user")
		return
	}

	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		mock_errors.WriteBadRequest(w, "error reading body")
		return
	}

	if userCtx.UserID != body.BroadcasterID {
		mock_errors.WriteUnauthorized(w, "broadcaster_id does not match token")
		return
	}

	if body.Title == "" || len(body.Title) > 45 {
		mock_errors.WriteBadRequest(w, "title is required and must be less than 45 characters")
		return
	}

	if body.PredictionWindow < 1 || body.PredictionWindow > 1800 {
		mock_errors.WriteBadRequest(w, "prediction_window is required and must between 1 and 1800")
		return
	}

	if len(body.Outcomes) < 2 || len(body.Outcomes) > 10 {
		mock_errors.WriteBadRequest(w, "Number of outcomes in the prediction must be equal to or above 2, and equal to or below 10")
		return
	}

	prediction := database.Prediction{
		ID:               util.RandomGUID(),
		BroadcasterID:    userCtx.UserID,
		BroadcasterLogin: u.UserLogin,
		BroadcasterName:  u.DisplayName,
		Title:            body.Title,
		WinningOutcomeID: nil,
		Status:           "ACTIVE",
		StartedAt:        util.GetTimestamp().Format(time.RFC3339),
		PredictionWindow: body.PredictionWindow,
	}

	for i, o := range body.Outcomes {
		color := "BLUE"
		if o.Title == "" {
			mock_errors.WriteBadRequest(w, "title is required for each outcome")
			return
		}
		if i == 1 {
			color = "PINK"
		}
		prediction.Outcomes = append(prediction.Outcomes, database.PredictionOutcome{
			ID:            util.RandomGUID(),
			Title:         o.Title,
			Users:         0,
			ChannelPoints: 0,
			Color:         color,
			PredictionID:  prediction.ID,
		})
	}
	err = db.NewQuery(r, 100).InsertPrediction(prediction)
	if err != nil {
		mock_errors.WriteBadRequest(w, "error inserting prediction")
		return
	}

	json.NewEncoder(w).Encode(models.APIResponse{Data: []database.Prediction{prediction}})
}

func patchPredictions(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	var body PatchPredictionsBody

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		mock_errors.WriteBadRequest(w, "error reading body")
		return
	}

	if userCtx.UserID != body.BroadcasterID {
		mock_errors.WriteUnauthorized(w, "broadcaster_id does not match token")
		return
	}

	if body.ID == "" {
		mock_errors.WriteBadRequest(w, "id is required")
		return
	}

	if body.Status != "RESOLVED" && body.Status != "CANCELED" && body.Status != "LOCKED" {
		mock_errors.WriteBadRequest(w, "status must be one of RESOLVED or CANCELED or LOCKED")
		return
	}

	if body.Status == "RESOLVED" && body.WinningOutcomeID == "" {
		mock_errors.WriteBadRequest(w, "winning_outcome_id is required if status is RESOLVED")
		return
	}

	err = db.NewQuery(r, 100).UpdatePrediction(database.Prediction{ID: body.ID, Status: body.Status, WinningOutcomeID: &body.WinningOutcomeID, BroadcasterID: body.BroadcasterID})
	if err != nil {
		mock_errors.WriteBadRequest(w, "error updating prediction")
		return
	}

	dbr, err := db.NewQuery(r, 100).GetPredictions(database.Prediction{ID: body.ID, BroadcasterID: body.BroadcasterID})
	if err != nil {
		mock_errors.WriteBadRequest(w, "error fetching prediction")
		return
	}

	prediction := dbr.Data.([]database.Prediction)

	json.NewEncoder(w).Encode(models.APIResponse{Data: prediction})
}
