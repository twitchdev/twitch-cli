// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package bits

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var leaderboardMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var leaderboardScopesByMethod = map[string][]string{
	http.MethodGet:    {"bits:read"},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type BitsLeaderboard struct{}

type BitsLeaderboardResponse struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
	Rank      int    `json:"rank"`
	Score     int    `json:"score"`
}

var validPeriod = map[string]time.Duration{
	"all":   1,
	"day":   24 * time.Hour,
	"week":  7 * 24 * time.Hour,
	"month": 7 * 24 * time.Hour,
	"year":  365 * 24 * time.Hour,
}

func (e BitsLeaderboard) Path() string { return "/bits/leaderboard" }

func (e BitsLeaderboard) GetRequiredScopes(method string) []string {
	return leaderboardScopesByMethod[method]
}

func (e BitsLeaderboard) ValidMethod(method string) bool {
	return leaderboardMethodsSupported[method]
}

func (e BitsLeaderboard) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getBitsLeaderboard(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getBitsLeaderboard(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	startedAt := r.URL.Query().Get("started_at")
	userID := r.URL.Query().Get("user_id")
	count := r.URL.Query().Get("count")
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)

	bl := []BitsLeaderboardResponse{}
	dateRange := models.BitsLeaderboardDateRange{}

	// check for empty auth context (aka app access tokens)
	if userCtx.UserID == "" {
		w.Write(mock_errors.GetErrorBytes(http.StatusUnauthorized, errors.New("Unauthorized"), "App access tokens are unable to use /bits/leaderboard"))
		return
	}

	// default the period value
	if period == "" {
		period = "all"
	}

	// check if provided period is not valid (not one of the valid values)
	if isValidPeriod(period) == false {
		w.Write(mock_errors.GetErrorBytes(http.StatusBadRequest, errors.New("Invalid Request"), "Period is invalid"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// default the count value if not provided
	if count == "" {
		count = "10"
	}

	// validate count param
	c, err := strconv.Atoi(count)
	if err != nil {
		w.Write(mock_errors.GetErrorBytes(http.StatusBadRequest, errors.New("Invalid Request"), "count is not a valid integer"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if c < 0 || c > 100 {
		w.Write(mock_errors.GetErrorBytes(http.StatusBadRequest, errors.New("Invalid Request"), "count must be between 1 and 100"))
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	// check if the started_at date is valid and then add it to the start/end range
	if period != "all" {
		if startedAt == "" {
			startedAt = time.Now().Format(time.RFC3339)
		}

		sa, err := time.Parse(time.RFC3339, startedAt)
		if err != nil {
			w.Write(mock_errors.GetErrorBytes(http.StatusBadRequest, errors.New("Invalid Request"), "started_at is not in RFC3339 format"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sa = time.Date(sa.Year(), sa.Month(), sa.Day(), 0, 0, 0, 0, sa.Location())
		switch period {
		case "day":
			dateRange = models.BitsLeaderboardDateRange{
				StartedAt: sa.Format(time.RFC3339),
				EndedAt:   time.Date(sa.Year(), sa.Month(), sa.Day()+1, 0, 0, 0, 0, sa.Location()).Format(time.RFC3339),
			}
			break
		case "week":
			weekday := time.Duration(sa.Weekday())
			sa = time.Date(sa.Year(), sa.Month(), sa.Day(), 0, 0, 0, 0, sa.Location()).Add(-1 * (weekday - 1) * 24 * time.Hour)
			dateRange = models.BitsLeaderboardDateRange{
				StartedAt: sa.Format(time.RFC3339),
				EndedAt:   time.Date(sa.Year(), sa.Month(), sa.Day()+7, 0, 0, 0, 0, sa.Location()).Format(time.RFC3339),
			}
			break
		case "month":
			sa = time.Date(sa.Year(), sa.Month(), 1, 0, 0, 0, 0, sa.Location())
			dateRange = models.BitsLeaderboardDateRange{
				StartedAt: sa.Format(time.RFC3339),
				EndedAt:   time.Date(sa.Year(), sa.Month()+1, 1, 0, 0, 0, 0, sa.Location()).Format(time.RFC3339),
			}
		case "year":
			sa = time.Date(sa.Year(), 1, 1, 0, 0, 0, 0, sa.Location())
			dateRange = models.BitsLeaderboardDateRange{
				StartedAt: sa.Format(time.RFC3339),
				EndedAt:   time.Date(sa.Year()+1, 1, 1, 0, 0, 0, 0, sa.Location()).Format(time.RFC3339),
			}
		case "all":
		default:
			break
		}
	}

	p := database.User{}
	if userID != "" {
		p.ID = userID
	}

	u, err := db.NewQuery(r, 100).GetUsers(p)
	if err != nil {
		w.Write(mock_errors.GetErrorBytes(http.StatusInternalServerError, err, "Error getting users for mock."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	users := u.Data.([]database.User)

	total := int(util.RandomInt(100 * 1000))

	for i := 0; i <= c-1; i++ {
		if len(users) == i {
			break
		}
		u := users[i]
		localTotal := math.Round(float64(total / (i + 1)))
		bl = append(bl, BitsLeaderboardResponse{
			UserID:    u.ID,
			UserLogin: u.UserLogin,
			UserName:  u.DisplayName,
			Rank:      i + 1,
			Score:     int(localTotal),
		})
	}

	length := len(bl)
	apiR := models.APIResponse{
		Data:  bl,
		Total: &length,
	}

	if dateRange.StartedAt != "" {
		apiR.DateRange = &dateRange
	}

	body, _ := json.Marshal(apiR)

	w.Write(body)
	return
}

func isValidPeriod(period string) bool {
	return validPeriod[period] > 0
}
