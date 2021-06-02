// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package users

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/models"
)

var methodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    true,
}

var scopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

var db database.CLIDatabase

type Endpoint struct{}

func (e Endpoint) Path() string { return "/users" }

func (e Endpoint) GetRequiredScopes(method string) []string {
	return scopesByMethod[method]
}

func (e Endpoint) ValidMethod(method string) bool {
	return methodsSupported[method]
}

func (e Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getUsers(w, r)
	case http.MethodPut:
		putUsers(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func convertUsers(users []database.User, includeEmail bool) []models.UserAPIData {
	r := []models.UserAPIData{}

	for _, u := range users {
		if u.ID == "" {
			continue
		}

		if includeEmail == false {
			u.Email = ""
		}

		r = append(r, models.UserAPIData{
			ID:              u.ID,
			Login:           u.UserLogin,
			DisplayName:     u.DisplayName,
			Type:            u.UserType,
			BroadcasterType: u.BroadcasterType,
			Description:     u.UserDescription,
			ViewCount:       0,
			Email:           u.Email,
			CreatedAt:       u.CreatedAt,
			OfflineImageURL: "https://static-cdn.jtvnw.net/jtv_user_pictures/3f13ab61-ec78-4fe6-8481-8682cb3b0ac2-channel_offline_image-1920x1080.png",
			ProfileImageURL: "https://static-cdn.jtvnw.net/jtv_user_pictures/8a6381c7-d0c0-4576-b179-38bd5ce1d6af-profile_image-300x300.png",
		})
	}
	return r
}

func getUsers(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()
	users := []database.User{}
	userIDs := q["id"]
	logins := q["login"]
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)

	shouldIncludeEmail := userCtx.HasScope("user:read:email")

	// TODO: update to handle in case the token contains the user info vs. simple params
	if len(logins) == 0 && len(userIDs) == 0 && len(userCtx.UserID) == 0 {
		w.WriteHeader(400)
		return
	}

	for _, i := range userIDs {
		user := database.User{
			ID: i,
		}
		println(i)
		u, err := db.NewQuery(r, 100).GetUser(user)
		log.Printf("%#v", u)
		if err != nil {
			log.Print(err.Error())
			w.WriteHeader(500)
			return
		}
		users = append(users, u)
	}

	for _, l := range logins {
		user := database.User{
			UserLogin: l,
		}
		u, err := db.NewQuery(r, 100).GetUser(user)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		users = append(users, u)
	}

	data := models.APIResponse{
		Data: convertUsers(users, shouldIncludeEmail),
	}

	body, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(500)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(body)
}

func putUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	users := []database.User{}
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	description := q["description"]
	shouldIncludeEmail := userCtx.HasScope("user:read:email")

	if userCtx.UserID == "" || len(description) == 0 {
		w.WriteHeader(400)
		return
	}

	u, err := db.NewQuery(r, 100).GetUser(database.User{ID: userCtx.UserID})
	if err != nil {
		log.Printf(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	u.UserDescription = description[len(description)-1]

	err = db.NewQuery(r, 100).InsertUser(u, true)
	if err != nil {
		log.Print(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := models.APIResponse{
		Data: convertUsers(append(users, u), shouldIncludeEmail),
	}

	body, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(500)
	}
	w.Write(body)
}
