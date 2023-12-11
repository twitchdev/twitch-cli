// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package chat

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
)

var colorMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    true,
}

var colorScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {"user:manage:chat_color"},
}

type GetColorRequestBody struct {
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	UserLogin string `json:"user_login"`
	Color     string `json:"color"`
}

type Color struct{}

func (e Color) Path() string { return "/chat/color" }

func (e Color) GetRequiredScopes(method string) []string {
	return colorScopesByMethod[method]
}

func (e Color) ValidMethod(method string) bool {
	return colorMethodsSupported[method]
}

func (e Color) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getColor(w, r)
		break
	case http.MethodPut:
		putColor(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

var validHexColorRegexp *regexp.Regexp = regexp.MustCompile("^#[a-fA-F0-9]{6}$")
var validNamedColorsLower = map[string]string{
	"blue":         "#0000FF",
	"blue_violet":  "#8A2BE2",
	"cadet_blue":   "#5F9EA0",
	"chocolate":    "#D2691E",
	"coral":        "#FF7F50",
	"dodger_blue":  "#1E90FF",
	"firebrick":    "#B22222",
	"golden_rod":   "#DAA520",
	"green":        "#008000",
	"hot_pink":     "#FF69B4",
	"orange_red":   "#FF4500",
	"red":          "#FF0000",
	"sea_green":    "#2E8B57",
	"spring_green": "#00FF7F",
	"yellow_green": "#9ACD32",
}

func getColor(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	userIDs := q["user_id"]
	results := []GetColorRequestBody{}

	if len(userIDs) == 0 {
		mock_errors.WriteBadRequest(w, "Missing required parameter user_id")
		return
	}

	if len(userIDs) > 100 {
		mock_errors.WriteBadRequest(w, "You may only specify up to 100 user_id query parameters")
		return
	}

	for _, i := range userIDs {
		user := database.User{
			ID: i,
		}
		u, err := db.NewQuery(r, 100).GetUser(user)
		if err != nil {
			log.Print(err.Error())
			w.WriteHeader(500)
			return
		}
		if u.ID == "" {
			continue
		}

		duplicate := false
		for _, d := range results {
			if d.UserID == u.ID {
				duplicate = true
			}
		}
		if duplicate {
			continue
		}

		results = append(results, GetColorRequestBody{
			UserID:    u.ID,
			UserName:  u.DisplayName,
			UserLogin: u.UserLogin,
			Color:     u.ChatColor,
		})
	}

	bytes, _ := json.Marshal(models.APIResponse{
		Data: results,
	})
	w.Write(bytes)
}

func putColor(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	if !userCtx.MatchesUserIDParam(r) {
		mock_errors.WriteUnauthorized(w, "User ID does not match token.")
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter user_id")
		return
	}

	color := r.URL.Query().Get("color")
	if color == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter color")
		return
	}

	// Users need to input %23 instead of # into their query. We store as # internally, so this needs to be converted
	color = strings.ReplaceAll(color, "%23", "#")

	// Check if named color is valid. If so, change the above color variable to the hex so it can pass the regex below.
	// This allows us to store color directly into the database without an additional variable.
	if namedColorHex, ok := validNamedColorsLower[strings.ToLower(color)]; ok {
		color = namedColorHex
	}

	validHex := validHexColorRegexp.MatchString(color)

	if !validHex {
		mock_errors.WriteBadRequest(w, "The color specified in the color query paramter is not valid")
		return
	}

	// Store in database, and return no body, just HTTP 204
	u, err := db.NewQuery(r, 100).GetUser(database.User{ID: userID})
	if err != nil {
		mock_errors.WriteServerError(w, "Error fetching user: "+err.Error())
		return
	}

	u.ChatColor = color
	log.Printf("%v", u)

	err = db.NewQuery(r, 100).InsertUser(u, true)
	if err != nil {
		mock_errors.WriteServerError(w, "Error writing to database: "+err.Error())
	}

	w.WriteHeader(http.StatusNoContent)
}
