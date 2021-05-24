// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package generate

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/util"
)

func Generate(userCount int) error {
	db, err := database.NewConnection()
	if err != nil {
		return err
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "db", db)

	// generate a client and fake secret
	c, err := generateClient(ctx)
	if err != nil {
		return err
	}
	generateAuthorization(ctx, c, "")

	// generate users and random related info (follows, bans, etc)
	generateUsers(ctx, userCount)

	return nil
}

func generateUsers(ctx context.Context, count int) error {
	db := ctx.Value("db").(database.CLIDatabase)
	var userIds []string

	// create users
	for i := 0; i < count; i++ {
		id := util.RandomUserID()
		userIds = append(userIds, id)

		un := generateUsername()

		bt := ""
		// status check
		t := util.RandomInt(3)

		if t == 1 {
			bt = "affiliate"
		} else if t == 2 {
			bt = "partner"
		}

		u := database.User{
			ID:              id,
			UserLogin:       strings.ToLower(un),
			DisplayName:     un,
			Email:           fmt.Sprintf("%v@testing.mocks", un),
			BroadcasterType: bt,
			UserType:        "",
			UserDescription: "",
			CreatedAt:       util.GetTimestamp().Format(time.RFC3339),
			ModifiedAt:      util.GetTimestamp().Format(time.RFC3339),
		}

		err := db.InsertUser(u, false)
		if err != nil {
			log.Print(err.Error())
		}
	}

	// create fake follows, blocks, mods
	for i, broadcaster := range userIds {
		for j, user := range userIds {
			// can't follow/block yourself :)
			if i == j {
				continue
			}

			userSeed := util.RandomInt(100 * 100)

			// 1 in 5 to follow
			shouldFollow := userSeed%5 == 0
			if shouldFollow {
				err := db.AddFollow(broadcaster, user)
				if err != nil {
					log.Print(err.Error())
				}
			}
			// 1 in 10 chance roughly to block one another
			shouldBlock := userSeed%10 == 0
			if shouldBlock {
				err := db.AddBlock(broadcaster, user)
				if err != nil {
					log.Print(err.Error())
				}
			}

			// 1 in 50 chance to mod one another, plus adds to the moderator events
			shouldMod := userSeed%50 == 0
			if shouldMod {
				err := db.AddModerator(broadcaster, user)
				if err != nil {
					log.Print(err.Error())
				}
			}
		}
	}

	// seed categories

	// create fake streams

	// create fake stream_markers

	// create fake videos

	// create fake polls

	// create fake predictions

	return nil
}

func generateClient(ctx context.Context) (database.AuthenticationClient, error) {
	db := ctx.Value("db").(database.CLIDatabase)

	client := database.AuthenticationClient{
		ID:          util.RandomClientID(),
		Name:        "Mock API Client",
		IsExtension: false,
	}

	client, err := db.InsertOrUpdateAuthenticationClient(client, false)
	log.Printf("Created Client. Details:\nClient-ID: %v\nSecret: %v\nName: %v", client.ID, client.Secret, client.Name)
	return client, err
}

func generateAuthorization(ctx context.Context, c database.AuthenticationClient, userID string) error {
	db := ctx.Value("db").(database.CLIDatabase)

	a := database.Authorization{
		ClientID:  c.ID,
		UserID:    sql.NullString{String: userID},
		ExpiresAt: util.GetTimestamp().Add(24 * time.Hour).Format(time.RFC3339),
	}

	auth, err := db.CreateAuthorization(a)
	if err != nil {
		return err
	}
	log.Printf("Created authorization for user %v with token %v", userID, auth.Token)
	return nil
}
