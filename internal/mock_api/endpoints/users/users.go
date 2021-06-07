// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package users

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/models"
)

var userMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    true,
}

var userScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {"user:edit"},
}

type User struct {
	ID              string `db:"id" json:"id"`
	UserLogin       string `db:"user_login" json:"login"`
	DisplayName     string `db:"display_name" json:"display_name"`
	Email           string `db:"email" json:"email,omitempty"`
	UserType        string `db:"user_type" json:"type"`
	BroadcasterType string `db:"broadcaster_type" json:"broadcaster_type"`
	UserDescription string `db:"user_description" json:"description"`
	CreatedAt       string `db:"created_at" json:"created_at"`
	ModifiedAt      string `db:"modified_at" json:"-"`
	ProfileImageURL string `dbi:"false" json:"profile_image_url" `
	OfflineImageURL string `dbi:"false" json:"offline_image_url" `
	ViewCount       int    `dbi:"false" json:"view_count"`
}

type UsersEndpoint struct{}

func (e UsersEndpoint) Path() string { return "/users" }

func (e UsersEndpoint) GetRequiredScopes(method string) []string {
	return userScopesByMethod[method]
}

func (e UsersEndpoint) ValidMethod(method string) bool {
	return userMethodsSupported[method]
}

func (e UsersEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func getUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	users := []User{}
	userIDs := q["id"]
	logins := q["login"]
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)

	shouldIncludeEmail := userCtx.HasScope("user:read:email")

	if len(logins) == 0 && len(userIDs) == 0 && len(userCtx.UserID) == 0 {
		log.Print(userCtx)
		w.WriteHeader(400)
		return
	}

	// add to list to get users if no logins or ids specified
	if len(userCtx.UserID) > 0 && len(logins) == 0 && len(userIDs) == 0 {
		userIDs = append(userIDs, userCtx.UserID)
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
		users = append(users, convertUsers([]database.User{u})...)
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
		if u.ID == "" {
			continue
		}
		users = append(users, convertUsers([]database.User{u})...)
	}

	// filter out emails for everyone but the authorized user
	for i, _ := range users {
		if users[i].ID != userCtx.UserID {
			users[i].Email = ""
		} else if !shouldIncludeEmail {
			users[i].Email = ""
		}
	}

	data := models.APIResponse{
		Data: users,
	}

	body, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(500)
	}
	w.Write(body)
}

func putUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
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
	users := convertUsers([]database.User{u})

	for i, _ := range users {
		if users[i].ID != userCtx.UserID {
			users[i].Email = ""
		} else if !shouldIncludeEmail {
			users[i].Email = ""
		}
	}

	if !shouldIncludeEmail {
		u.Email = ""
	}

	data := models.APIResponse{
		Data: users,
	}

	body, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(500)
	}
	w.Write(body)
}

func convertUsers(users []database.User) []User {
	response := []User{}
	if len(users) == 0 {
		return response
	}
	for _, u := range users {
		response = append(response, User{
			ID:              u.ID,
			UserLogin:       u.UserLogin,
			DisplayName:     u.DisplayName,
			Email:           u.Email,
			UserType:        u.UserType,
			BroadcasterType: u.BroadcasterType,
			CreatedAt:       u.CreatedAt,
			ProfileImageURL: u.ProfileImageURL,
			OfflineImageURL: u.OfflineImageURL,
			ViewCount:       u.ViewCount,
			UserDescription: u.UserDescription,
		})
	}
	return response
}
