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

type UserInfo struct {
	ID   string
	Type string
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

	log.Print("Finished generation.")
	return nil
}

func generateUsers(ctx context.Context, count int) error {
	db := ctx.Value("db").(database.CLIDatabase)
	var users []UserInfo
	var categoryIds []string
	var streams []Stream
	var tagIds []string

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

	// create users
	log.Printf("Creating users...")
	for i := 0; i < count; i++ {
		id := util.RandomUserID()

		un := generateUsername()

		bt := ""
		// status check
		t := util.RandomInt(3)

		if i == 0 {
			t = 2
		}

		if t == 1 {
			bt = "affiliate"
		} else if t == 2 {
			bt = "partner"
		}

		users = append(users, UserInfo{ID: id, Type: bt})

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
			CategoryID:      sql.NullString{String: categoryIds[util.RandomInt(int64(len(categoryIds)-1))], Valid: true},
			Title:           "Sample stream!",
			Language:        "en",
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

	// drops entitlements static data
	benefitID := util.RandomGUID()
	dropsGameID := categoryIds[len(categoryIds)-1]

	// t is true for use with a *bool value used below
	t := true
	// create fake follows, blocks, mods, and team membership
	log.Printf("Creating channel points rewards and redemptions, follows, blocks, mods, bans, editors, and team members...")
	for i, broadcaster := range users {
		copoReward := database.ChannelPointsReward{}
		prediction := database.Prediction{}
		poll := database.Poll{}
		if broadcaster.Type != "" {
			copoReward = database.ChannelPointsReward{
				ID:                  util.RandomGUID(),
				BroadcasterID:       broadcaster.ID,
				RewardImage:         sql.NullString{},
				BackgroundColor:     "#9146FF",
				IsEnabled:           &t,
				Cost:                1,
				Title:               fmt.Sprintf("Fake reward for %v", broadcaster.ID),
				RewardPrompt:        "",
				IsUserInputRequired: false,
				MaxPerStream: database.MaxPerStream{
					StreamMaxEnabled: false,
					StreamMaxCount:   0,
				},
				MaxPerUserPerStream: database.MaxPerUserPerStream{
					StreamUserMaxEnabled: false,
					StreamMUserMaxCount:  0,
				},
				GlobalCooldown: database.GlobalCooldown{
					GlobalCooldownEnabled: false,
					GlobalCooldownSeconds: 0,
				},
				IsPaused:                         false,
				IsInStock:                        true,
				ShouldRedemptionsSkipQueue:       false,
				RedemptionsRedeemedCurrentStream: nil,
				CooldownExpiresAt:                sql.NullString{},
			}

			err := db.NewQuery(nil, 100).InsertChannelPointsReward(copoReward)
			if err != nil {
				log.Print(err.Error())
			}
		}

		entitlement := database.DropsEntitlement{
			ID:        util.RandomGUID(),
			BenefitID: benefitID,
			GameID:    dropsGameID,
			UserID:    broadcaster.ID,
			Timestamp: util.GetTimestamp().Format(time.RFC3339Nano),
		}
		err = db.NewQuery(nil, 1000).InsertDropsEntitlement(entitlement)
		if err != nil {
			log.Print(err.Error())
		}

		// create fake polls
		poll = database.Poll{
			ID:                         util.RandomGUID(),
			Title:                      "Test title",
			BroadcasterID:              broadcaster.ID,
			BitsVotingEnabled:          false,
			ChannelPointsVotingEnabled: false,
			Status:                     "ACTIVE",
			Duration:                   300,
			StartedAt:                  util.GetTimestamp().Format(time.RFC3339),
		}

		poll.Choices = []database.PollsChoice{
			{Title: "Choice 1", Votes: 0, ChannelPointsVotes: 0, BitsVotes: 0, PollID: poll.ID, ID: util.RandomGUID()},
			{Title: "Choice 2", Votes: 0, ChannelPointsVotes: 0, BitsVotes: 0, PollID: poll.ID, ID: util.RandomGUID()},
		}

		err = db.NewQuery(nil, 100).InsertPoll(poll)
		if err != nil {
			log.Print(err.Error())
		}

		// create fake predictions
		prediction = database.Prediction{
			ID:               util.RandomGUID(),
			BroadcasterID:    broadcaster.ID,
			Title:            "Test Prediction",
			PredictionWindow: 1000,
			Status:           "ACTIVE",
			StartedAt:        util.GetTimestamp().Format(time.RFC3339),
		}

		prediction.Outcomes = []database.PredictionOutcome{
			{
				ID:            util.RandomGUID(),
				Title:         "Choice1",
				Color:         "BLUE",
				Users:         0,
				ChannelPoints: 0,
				PredictionID:  prediction.ID,
			},
			{
				ID:            util.RandomGUID(),
				Title:         "Choice1",
				Color:         "PINK",
				Users:         0,
				ChannelPoints: 0,
				PredictionID:  prediction.ID,
			},
		}

		err = db.NewQuery(nil, 100).InsertPrediction(prediction)
		if err != nil {
			log.Print(err.Error())
		}

		for j, user := range users {
			// create a seed used for the below determination on if a user should follow one another- this simply simulates a social mesh
			userSeed := util.RandomInt(100 * 100)
			if copoReward.ID != "" {
				copoRedemption := database.ChannelPointsRedemption{
					ID:            util.RandomGUID(),
					BroadcasterID: broadcaster.ID,
					UserID:        user.ID,
					RewardID:      copoReward.ID,
					ChannelPointsRedemptionRewardInfo: database.ChannelPointsRedemptionRewardInfo{
						ID:           copoReward.ID,
						Title:        copoReward.Title,
						RewardPrompt: copoReward.RewardPrompt,
						Cost:         copoReward.Cost,
					},
					UserInput:        sql.NullString{},
					RedemptionStatus: "UNFULFILLED",
					RedeemedAt:       util.GetTimestamp().Format(time.RFC3339),
				}
				err := db.NewQuery(nil, 100).InsertChannelPointsRedemption(copoRedemption)
				if err != nil {
					log.Print(err.Error())
				}
			}

			if poll.ID != "" {
				choice := util.RandomInt(2)
				err = db.NewQuery(nil, 100).UpdatePollChoice(database.PollsChoice{ID: poll.Choices[choice].ID})
				if err != nil {
					log.Print(err.Error())
				}
			}

			if prediction.ID != "" {
				option := util.RandomInt(2)
				err = db.NewQuery(nil, 100).InsertPredictionPrediction(database.PredictionPrediction{PredictionID: prediction.ID, UserID: user.ID, Amount: int(util.RandomInt(10 * 1000)), OutcomeID: prediction.Outcomes[option].ID})
				if err != nil {
					log.Print(err.Error())
				}
			}
			// can't follow/block yourself :)
			if i == j {
				continue
			}

			// 1 in 25 chance roughly to block one another
			shouldBlock := userSeed%25 == 0
			if shouldBlock {
				err := db.NewQuery(nil, 100).AddBlock(database.UserRequestParams{UserID: user.ID, BroadcasterID: broadcaster.ID})
				if err != nil {
					log.Print(err.Error())
				}
				// since you're blocked, can't do any of the other things, so continue
				continue
			}
			shouldBan := userSeed%20 == 0
			if shouldBan {
				err := db.NewQuery(nil, 100).InsertBan(database.UserRequestParams{UserID: user.ID, BroadcasterID: broadcaster.ID})
				if err != nil {
					log.Print(err.Error())
				}
				// if banned, you wouldn't be able to follow or anything else- so continuing
				continue
			}

			shouldFollow := userSeed%5 == 0 || i == 0
			if shouldFollow {
				err := db.NewQuery(nil, 100).AddFollow(database.UserRequestParams{UserID: user.ID, BroadcasterID: broadcaster.ID})
				if err != nil {
					log.Print(err.Error())
				}
			}

			shouldSub := userSeed%10 == 0 || i == 0
			if shouldSub && broadcaster.Type != "" {
				err := db.NewQuery(nil, 100).InsertSubscription(database.SubscriptionInsert{
					UserID:        user.ID,
					BroadcasterID: broadcaster.ID,
					Tier:          fmt.Sprint((util.RandomInt(3) + 1) * 1000),
					CreatedAt:     util.GetTimestamp().Format(time.RFC3339),
					IsGift:        false,
				})
				if err != nil {
					log.Print(err.Error())
				}
			}

			shouldMod := userSeed%10 == 0 || i == 0
			if shouldMod {
				err := db.NewQuery(nil, 100).AddModerator(database.UserRequestParams{UserID: user.ID, BroadcasterID: broadcaster.ID})
				if err != nil {
					log.Print(err.Error())
				}
			}

			shouldAddEditor := userSeed%20 == 0 || i == 0
			if shouldAddEditor {
				err := db.NewQuery(nil, 100).AddEditor(database.UserRequestParams{BroadcasterID: broadcaster.ID, UserID: user.ID})
				if err != nil {
					log.Print(err.Error())
				}
			}

		}

		shouldBeTeamMember := util.RandomInt(100*100)%20 == 0

		if i == 0 {
			shouldBeTeamMember = true
		}
		if shouldBeTeamMember {
			err := db.NewQuery(nil, 100).InsertTeamMember(database.TeamMember{
				TeamID: team.ID,
				UserID: broadcaster.ID,
			})
			if err != nil {
				log.Print(err.Error())
			}
		}

	}

	// create fake streams
	log.Printf("Creating streams...")
	for i, u := range users {
		if util.RandomInt(100)%10 != 0 && i != 0 {
			continue
		}
		s := database.Stream{
			ID:          util.RandomGUID(),
			UserID:      u.ID,
			StreamType:  "live",
			ViewerCount: int(util.RandomViewerCount()),
			StartedAt:   util.GetTimestamp().Format(time.RFC3339),
			IsMature:    false,
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
	log.Printf("Creating stream tags, videos, clips, and stream markers...")
	for _, s := range streams {
		var prevTag string
		for i := 0; i < int(util.RandomInt(5)); i++ {
			st := database.StreamTag{
				UserID: s.Broacaster,
				TagID:  tagIds[util.RandomInt(int64(len(tagIds)-1))],
			}
			if prevTag == st.TagID {
				continue
			}

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
			StreamID:         &s.ID,
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

		// clips
		c := database.Clip{
			ID:            util.RandomGUID(), // does not follow the same slug format, but for this it'll do
			BroadcasterID: s.Broacaster,
			GameID:        categoryIds[util.RandomInt(int64(len(categoryIds)-1))],
			CreatorID:     s.Broacaster,
			VideoID:       v.ID,
			Title:         "Generated clip!",
			Language:      "en",
			ViewCount:     0,
			CreatedAt:     util.GetTimestamp().Format(time.RFC3339),
			Duration:      30.1,
		}
		err = db.NewQuery(nil, 100).InsertClip(c)
		if err != nil {
			log.Print(err.Error())
		}
	}

	if len(users) > 0 {
		// log out that the user X has all units for easier getting started
		log.Printf("User ID %v has all applicable units (streams, subs, and the like) and is a partner for use with the API", users[0].ID)
	}

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
		UserID:    userID,
		ExpiresAt: util.GetTimestamp().Add(24 * time.Hour).Format(time.RFC3339),
		Scopes:    "",
	}

	auth, err := db.NewQuery(nil, 100).CreateAuthorization(a)
	if err != nil {
		return err
	}
	log.Printf("Created authorization for user %v with token %v", userID, auth.Token)
	return nil
}
