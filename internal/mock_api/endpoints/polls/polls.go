// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package polls

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

var pollsMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   true,
	http.MethodDelete: false,
	http.MethodPatch:  true,
	http.MethodPut:    false,
}

var pollsScopesByMethod = map[string][]string{
	http.MethodGet:    {"channel:read:polls", "channel:manage:polls"},
	http.MethodPost:   {"channel:manage:polls"},
	http.MethodDelete: {},
	http.MethodPatch:  {"channel:manage:polls"},
	http.MethodPut:    {},
}

type Polls struct{}

type PostPollsBody struct {
	BroadcasterID              string                `json:"broadcaster_id"`
	Title                      string                `json:"title"`
	Choices                    []PostPollsBodyChoice `json:"choices"`
	Duration                   int                   `json:"duration"`
	BitsVotingEnabled          bool                  `json:"bits_voting_enabled"`
	BitsPerVote                int                   `json:"bits_per_vote"`
	ChannelPointsVotingEnabled bool                  `json:"channel_points_voting_enabled"`
	ChannelPointsPerVote       int                   `json:"channel_points_per_vote"`
}

type PostPollsBodyChoice struct {
	Title string `json:"title"`
}

type PatchPollsBody struct {
	BroadcasterID string `json:"broadcaster_id"`
	ID            string `json:"id"`
	Status        string `json:"status"`
}

func (e Polls) Path() string { return "/polls" }

func (e Polls) GetRequiredScopes(method string) []string {
	return pollsScopesByMethod[method]
}

func (e Polls) ValidMethod(method string) bool {
	return pollsMethodsSupported[method]
}

func (e Polls) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getPolls(w, r)
		break
	case http.MethodPost:
		postPolls(w, r)
		break
	case http.MethodPatch:
		patchPolls(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
func getPolls(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	var dbr *database.DBResponse
	var err error
	polls := []database.Poll{}

	if !userCtx.MatchesBroadcasterIDParam(r) {
		mock_errors.WriteBadRequest(w, "broadcaster_id does not match token")
		return
	}

	ids := r.URL.Query()["id"]

	for _, id := range ids {
		dbr, err := db.NewQuery(r, 100).GetPolls(database.Poll{ID: id})
		if err != nil {
			mock_errors.WriteServerError(w, "error fetching polls")
			return
		}
		polls = append(polls, dbr.Data.([]database.Poll)...)
	}
	if len(ids) == 0 {
		dbr, err = db.NewQuery(r, 100).GetPolls(database.Poll{BroadcasterID: userCtx.UserID})
		if err != nil {
			mock_errors.WriteServerError(w, "error fetching polls")
			return
		}
		polls = append(polls, dbr.Data.([]database.Poll)...)
	}

	apiResposne := models.APIResponse{
		Data: polls,
	}

	if dbr.Cursor != "" {
		apiResposne.Pagination = &models.APIPagination{
			Cursor: dbr.Cursor,
		}
	}

	bytes, _ := json.Marshal(apiResposne)
	w.Write(bytes)
}

func postPolls(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)

	u, err := db.NewQuery(r, 100).GetUser(database.User{ID: userCtx.UserID})
	if err != nil {
		mock_errors.WriteBadRequest(w, "error getting user")
		return
	}

	var body PostPollsBody

	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		mock_errors.WriteBadRequest(w, "error reading body")
		return
	}

	if body.BroadcasterID != userCtx.UserID {
		mock_errors.WriteBadRequest(w, "broadcaster_id does not match token")
		return
	}

	if body.Title == "" {
		mock_errors.WriteBadRequest(w, "title is required")
		return
	}

	if len(body.Choices) < 2 || len(body.Choices) > 5 {
		mock_errors.WriteBadRequest(w, "you may only have between 2 and 5 choices")
		return
	}

	if body.Duration < 15 || body.Duration > 1800 {
		mock_errors.WriteBadRequest(w, "duation must be at least15 and at most 1800")
		return
	}
	poll := database.Poll{
		ID:                         util.RandomGUID(),
		BroadcasterID:              userCtx.UserID,
		BroadcasterLogin:           u.UserLogin,
		BroadcasterName:            u.DisplayName,
		Title:                      body.Title,
		BitsVotingEnabled:          body.BitsVotingEnabled,
		BitsPerVote:                body.BitsPerVote,
		ChannelPointsVotingEnabled: body.ChannelPointsVotingEnabled,
		ChannelPointsPerVote:       body.ChannelPointsPerVote,
		Status:                     "ACTIVE",
		Duration:                   body.Duration,
		StartedAt:                  util.GetTimestamp().Format(time.RFC3339),
	}

	for _, c := range body.Choices {
		if c.Title == "" {
			mock_errors.WriteBadRequest(w, "each choice must have a title")
			return
		}

		poll.Choices = append(poll.Choices, database.PollsChoice{
			ID:                 util.RandomGUID(),
			Title:              c.Title,
			Votes:              0,
			ChannelPointsVotes: 0,
			BitsVotes:          0,
			PollID:             poll.ID,
		})
	}

	err = db.NewQuery(r, 100).InsertPoll(poll)
	if err != nil {
		mock_errors.WriteServerError(w, "error inserting poll")
		return
	}

	json.NewEncoder(w).Encode(models.APIResponse{Data: []database.Poll{poll}})
}

func patchPolls(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)

	var body PatchPollsBody

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		mock_errors.WriteBadRequest(w, "error reading body")
		return
	}
	if body.BroadcasterID != userCtx.UserID {
		mock_errors.WriteBadRequest(w, "broadcaster_id does not match token")
		return
	}
	if body.ID == "" {
		mock_errors.WriteBadRequest(w, "id is required")
		return
	}
	if body.Status != "TERMINATED" && body.Status != "ARCHIVED" {
		mock_errors.WriteBadRequest(w, "status must be one of TERMINATED or ARCHIVED")
		return
	}

	err = db.NewQuery(r, 100).UpdatePoll(database.Poll{ID: body.ID, Status: body.Status, EndedAt: util.GetTimestamp().Format(time.RFC3339)})

	dbr, err := db.NewQuery(r, 100).GetPolls(database.Poll{BroadcasterID: userCtx.UserID, ID: body.ID})
	if err != nil {
		println(err.Error())
		mock_errors.WriteServerError(w, "error fetching polls")
		return
	}

	apiResposne := models.APIResponse{
		Data: dbr.Data,
	}

	if dbr.Cursor != "" {
		apiResposne.Pagination = &models.APIPagination{
			Cursor: dbr.Cursor,
		}
	}

	bytes, _ := json.Marshal(apiResposne)
	w.Write(bytes)
}
