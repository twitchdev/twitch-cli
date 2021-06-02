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

type Stream struct {
	ID         string
	Broacaster string
}

func Generate(userCount int) error {
	db, err := database.NewConnection()
	if err != nil {
		return err
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "db", db)

	// generate users and random related info (follows, bans, etc)
	generateUsers(ctx, userCount)

	// generate a client and fake secret
	c, err := generateClient(ctx)
	if err != nil {
		return err
	}
	generateAuthorization(ctx, c, "")

	return nil
}

func generateUsers(ctx context.Context, count int) error {
	db := ctx.Value("db").(database.CLIDatabase)
	var userIds []string
	var categoryIds []string
	var streams []Stream
	var tagIds []string

	// create users
	log.Printf("Creating users...")
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

		err := db.NewQuery(nil, 100).InsertUser(u, false)
		if err != nil {
			log.Print(err.Error())
		}
	}
	// fake team
	log.Printf("Creating team...")
	team := database.Team{
		ID:              fmt.Sprint(util.RandomInt(10 * 1000)),
		TeamName:        "clidev",
		TeamDisplayName: "CLI Developers",
		CreatedAt:       util.GetTimestamp().Format(time.RFC3339),
		UpdatedAt:       util.GetTimestamp().Format(time.RFC3339),
	}

	err := db.NewQuery(nil, 100).InsertTeam(team)
	if err != nil {
		log.Print(err.Error())
	}

	// create fake follows, blocks, mods, and team membership
	log.Printf("Creating follows, blocks, mods, bans, and team members...")
	for i, broadcaster := range userIds {
		for j, user := range userIds {
			// can't follow/block yourself :)
			if i == j {
				continue
			}

			userSeed := util.RandomInt(100 * 100)
			// 1 in 25 chance roughly to block one another
			shouldBlock := userSeed%25 == 0
			if shouldBlock {
				err := db.NewQuery(nil, 100).AddBlock(database.UserRequestParams{UserID: user, BroadcasterID: broadcaster})
				if err != nil {
					log.Print(err.Error())
				}
				// since you're blocked, can't do any of the other things, so continue
				continue
			}

			// 1 in 5 to follow
			shouldFollow := userSeed%5 == 0
			if shouldFollow {
				err := db.NewQuery(nil, 100).AddFollow(database.UserRequestParams{UserID: user, BroadcasterID: broadcaster})
				if err != nil {
					log.Print(err.Error())
				}
			}

			// 1 in 50 chance to mod one another, plus adds to the moderator events
			shouldMod := userSeed%50 == 0
			if shouldMod {
				err := db.NewQuery(nil, 100).AddModerator(database.UserRequestParams{UserID: user, BroadcasterID: broadcaster})
				if err != nil {
					log.Print(err.Error())
				}
			}

			// 1 in 100 chance to ban one another, plus adds to banned events
			shouldBan := userSeed%100 == 0
			if shouldBan {
				err := db.NewQuery(nil, 100).InsertBan(database.UserRequestParams{UserID: user, BroadcasterID: broadcaster})
				if err != nil {
					log.Print(err.Error())
				}
			}

			shouldSub := userSeed%10 == 0
			if shouldSub {
				err := db.NewQuery(nil, 100).InsertSubscription(database.SubscriptionInsert{
					UserID:        user,
					BroadcasterID: broadcaster,
					Tier:          fmt.Sprint((util.RandomInt(3) + 1) * 1000),
					CreatedAt:     util.GetTimestamp().Format(time.RFC3339),
					IsGift:        false,
				})
				if err != nil {
					log.Print(err.Error())
				}
			}
		}

		shouldBeTeamMember := util.RandomInt(100*100)%20 == 0

		if shouldBeTeamMember {
			err := db.NewQuery(nil, 100).InsertTeamMember(database.TeamMember{
				TeamID: team.ID,
				UserID: broadcaster,
			})
			if err != nil {
				log.Print(err.Error())
			}
		}
	}

	// seed categories
	log.Printf("Creating categories...")
	for _, c := range categories {
		category := database.Category{
			ID:   fmt.Sprintf("%v", util.RandomInt(10*100*100)),
			Name: c,
		}

		err := db.NewQuery(nil, 100).InsertCategory(category, false)
		if err != nil {
			log.Print(err.Error())
		}
		category = database.Category{
			ID: category.ID,
		}
		categoryIds = append(categoryIds, category.ID)
	}

	// create fake streams
	log.Printf("Creating streams...")
	for _, u := range userIds {
		if util.RandomInt(100)%25 != 0 {
			continue
		}
		s := database.Stream{
			ID:             util.RandomGUID(),
			UserID:         u,
			CategoryID:     categoryIds[util.RandomInt(int64(len(categoryIds)-1))],
			StreamType:     "live",
			Title:          "Sample stream!",
			ViewerCount:    int(util.RandomViewerCount()),
			StartedAt:      util.GetTimestamp().Format(time.RFC3339),
			StreamLanguage: "en",
			IsMature:       false,
		}
		err := db.NewQuery(nil, 100).InsertStream(s, false)
		if err != nil {
			log.Print(err.Error())
		}
		streams = append(streams, Stream{ID: s.ID, Broacaster: s.UserID})
	}

	log.Printf("Creating tags...")
	for _, t := range tags {
		tag := database.Tag{
			ID:   util.RandomGUID(),
			Name: t,
		}
		err := db.NewQuery(nil, 100).InsertTag(tag)
		if err != nil {
			log.Print(err.Error())
		}
		tagIds = append(tagIds, tag.ID)
	}

	// creates fake stream tags, videos, and markers
	log.Printf("Creating stream tags, videos, and stream markers...")
	for _, s := range streams {
		var prevTag string
		for i := 0; i < int(util.RandomInt(5)); i++ {
			st := database.StreamTag{
				StreamID: s.ID,
				TagID:    tagIds[util.RandomInt(int64(len(tagIds)-1))],
			}
			if prevTag == st.TagID {
				continue
			}
			log.Printf("%#v", st)

			err := db.NewQuery(nil, 100).InsertStreamTag(st)
			if err != nil {
				log.Print(err.Error())
			}
			prevTag = st.TagID
		}
		// markers

		// videos
		v := database.Video{
			ID:               fmt.Sprint(util.RandomInt(10 * 1000 * 1000)),
			StreamID:         s.ID,
			BroadcasterID:    s.Broacaster,
			Title:            "Sample stream!",
			VideoDescription: "",
			CreatedAt:        util.GetTimestamp().Format(time.RFC3339),
			PublishedAt:      util.GetTimestamp().Format(time.RFC3339),
			Viewable:         "public",
			ViewCount:        0,
			Duration:         "1h0m0s",
			VideoLanguage:    "en",
		}
		err := db.NewQuery(nil, 100).InsertVideo(v)
		if err != nil {
			log.Print(err.Error())
		}

		// muted sections
		for i := 0; i < int(util.RandomInt(2)); i++ {
			vms := database.VideoMutedSegment{
				VideoID:     v.ID,
				VideoOffset: int(util.RandomInt(45) * 60),
				Duration:    int(util.RandomInt(30)),
			}
			err := db.NewQuery(nil, 100).InsertMutedSegmentsForVideo(vms)
			if err != nil {
				log.Print(err.Error())
			}
		}
	}

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

	client, err := db.NewQuery(nil, 100).InsertOrUpdateAuthenticationClient(client, false)
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

	auth, err := db.NewQuery(nil, 100).CreateAuthorization(a)
	if err != nil {
		return err
	}
	log.Printf("Created authorization for user %v with token %v", userID, auth.Token)
	return nil
}