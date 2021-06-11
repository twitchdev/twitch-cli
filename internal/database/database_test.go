// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"github.com/twitchdev/twitch-cli/internal/util"
	"github.com/twitchdev/twitch-cli/test_setup"
)

const TEST_USER_ID = "1"
const TEST_USER_LOGIN = "testing_user1"

func TestGetDatabase(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	p, _ := util.GetApplicationDir()

	dbFileName = viper.GetString("DB_FILENAME")

	// delete the existing temp db if it exists
	path := filepath.Join(p, dbFileName)
	err := os.Remove(path)

	// if the error is not that the file doesn't exist, fail the test
	if !os.IsNotExist(err) {
		a.Nil(err)
	}

	// since this creates a new db, will check those codepaths
	db, err := getDatabase()
	a.Nil(err)
	a.NotNil(db)

	// get again, making sure that this works
	db, err = getDatabase()
	a.Nil(err)
	a.NotNil(db)
}

func TestRetriveFromDB(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	db, err := NewConnection()
	a.Nil(err)

	ecParams := *&EventCacheParameters{
		ID:        util.RandomGUID(),
		Event:     "foo",
		JSON:      "bar",
		FromUser:  "1234",
		ToUser:    "5678",
		Transport: "test",
		Timestamp: util.GetTimestamp().Format(time.RFC3339Nano),
	}

	q := Query{DB: db.DB}

	err = q.InsertIntoDB(ecParams)
	a.Nil(err)

	dbResponse, err := q.GetEventByID(ecParams.ID)
	a.Nil(err)

	a.NotNil(dbResponse)
	a.Equal("test", dbResponse.Transport)
}

func TestGenerateString(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	s := generateString(10)

	a.Len(s, 10)
}

func TestAuthentication(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	client := AuthenticationClient{ID: "1234", Secret: "4567", Name: "for_testing", IsExtension: false}
	db, _ := NewConnection()
	q := db.NewQuery(nil, 100)

	// test true insert
	ac, err := q.InsertOrUpdateAuthenticationClient(client, false)
	a.Nil(err)
	a.NotNil(ac)

	// if duped, should give a fresh client ID
	ac, err = q.InsertOrUpdateAuthenticationClient(client, false)
	a.Nil(err)
	a.NotNil(ac)
	a.NotEqual(ac.ID, client.ID, fmt.Sprintf("%v %v", ac.ID, client.ID))

	// test upsert
	ac, err = q.InsertOrUpdateAuthenticationClient(client, true)
	a.Nil(err)
	a.NotNil(ac)

	// create a fake auth
	auth, err := q.CreateAuthorization(Authorization{ClientID: ac.ID})
	a.Nil(err)
	a.NotNil(auth)

	// test fetching client
	dbr, err := q.GetAuthenticationClient(AuthenticationClient{ID: auth.ClientID})
	a.Nil(err)
	c := dbr.Data.([]AuthenticationClient)
	a.NotNil(c)
	a.Len(c, 1)
	a.Equal(c[0].ID, client.ID)

	authorization, err := q.GetAuthorizationByToken(auth.Token)
	a.Nil(err)
	a.Equal(client.ID, authorization.ClientID)
}

func TestAPI(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	db, _ := NewConnection()
	b := db.IsFirstRun()

	a.Equal(b, true)
}

func TestCategories(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	db, err := NewConnection()
	a.Nil(err)

	q := db.NewQuery(nil, 100)

	c := Category{Name: "test", ID: "1"}
	err = q.InsertCategory(c, false)
	a.Nil(err)

	// get categories
	dbr, err := q.GetCategories(Category{ID: c.ID})
	a.Nil(err)
	categories := dbr.Data.([]Category)
	a.Len(categories, 1)
	a.Equal(c.ID, categories[0].ID)

	// search
	dbr, err = q.SearchCategories("es")
	a.Nil(err)
	categories = dbr.Data.([]Category)
	a.Len(categories, 1)
	a.Equal(c.ID, categories[0].ID)

	// top
	dbr, err = q.GetTopGames()
	a.Nil(err)
	categories = dbr.Data.([]Category)
	a.Len(categories, 1)
	a.Equal(c.ID, categories[0].ID)
}

func TestUsers(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	db, err := NewConnection()
	a.Nil(err)

	q := db.NewQuery(nil, 100)

	err = q.InsertUser(User{
		ID:              TEST_USER_ID,
		UserLogin:       TEST_USER_LOGIN,
		DisplayName:     TEST_USER_LOGIN,
		Email:           "",
		BroadcasterType: "partner",
		UserType:        "testing",
		UserDescription: "hi mom",
		CreatedAt:       util.GetTimestamp().Format(time.RFC3339),
		ModifiedAt:      util.GetTimestamp().Format(time.RFC3339),
		CategoryID:      sql.NullString{String: "1", Valid: true},
		Title:           "hello",
		Language:        "en",
		Delay:           0,
	}, false)
	a.Nil(err)

	err = q.InsertUser(User{
		ID:              "2",
		UserLogin:       "second_user",
		DisplayName:     "second_user",
		Email:           "",
		BroadcasterType: "partner",
		UserType:        "testing",
		UserDescription: "hi mom",
		CreatedAt:       util.GetTimestamp().Format(time.RFC3339),
		ModifiedAt:      util.GetTimestamp().Format(time.RFC3339),
		CategoryID:      sql.NullString{String: "", Valid: false},
		Title:           "hello",
		Language:        "en",
		Delay:           0,
	}, false)
	a.Nil(err)

	u, err := q.GetUser(User{ID: TEST_USER_ID})
	a.Nil(err)
	a.Equal(TEST_USER_ID, u.ID)
	a.Equal(TEST_USER_LOGIN, u.UserLogin)

	dbr, err := q.GetUsers(User{ID: TEST_USER_ID})
	a.Nil(err)
	users := dbr.Data.([]User)
	a.Len(users, 1)

	dbr, err = q.GetChannels(User{ID: TEST_USER_ID})
	a.Nil(err)
	channels := dbr.Data.([]User)
	a.Len(channels, 1)
	a.Equal(channels[0].CategoryID.String, "1")

	// urp
	urp := UserRequestParams{BroadcasterID: TEST_USER_ID, UserID: "2"}
	err = q.AddFollow(urp)
	a.Nil(err)

	dbr, err = q.GetFollows(urp)
	a.Nil(err)
	follows := dbr.Data.([]Follow)
	a.Len(follows, 1)

	err = q.DeleteFollow(urp.UserID, urp.BroadcasterID)
	a.Nil(err)

	err = q.AddBlock(urp)
	a.Nil(err)

	dbr, err = q.GetBlocks(urp)
	a.Nil(err)
	blocks := dbr.Data.([]Block)
	a.Len(blocks, 1)

	err = q.DeleteBlock(urp.UserID, urp.BroadcasterID)
	a.Nil(err)

	err = q.AddEditor(urp)
	a.Nil(err)

	dbr, err = q.GetEditors(User{ID: urp.BroadcasterID})
	a.Nil(err)
	editors := dbr.Data.([]Editor)
	a.Len(editors, 1)

	err = q.UpdateChannel(urp.BroadcasterID, User{ID: urp.BroadcasterID, UserDescription: "hi mom2"})
	a.Nil(err)

	dbr, err = q.GetUsers(User{ID: TEST_USER_ID})
	a.Nil(err)
	users = dbr.Data.([]User)
	a.Len(users, 1)
	a.Equal("hi mom2", users[0].UserDescription)

	dbr, err = q.SearchChannels("testing_", false)
	a.Nil(err)
	search := dbr.Data.([]SearchChannel)
	a.Len(search, 1)
	a.Equal(TEST_USER_ID, search[0].ID)

	dbr, err = q.SearchChannels("testing_", true)
	a.Nil(err)
	search = dbr.Data.([]SearchChannel)
	a.Len(search, 0)
}

func TestChannelPoints(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	db, err := NewConnection()
	a.Nil(err)

	q := db.NewQuery(nil, 100)

	bTrue := true

	reward := ChannelPointsReward{
		ID:                         "1",
		BroadcasterID:              TEST_USER_ID,
		BackgroundColor:            "#fff",
		IsEnabled:                  &bTrue,
		Cost:                       100,
		Title:                      "1234",
		RewardPrompt:               "",
		IsUserInputRequired:        false,
		IsPaused:                   false,
		IsInStock:                  false,
		ShouldRedemptionsSkipQueue: false,
	}

	err = q.InsertChannelPointsReward(reward)
	a.Nil(err)

	reward.Cost = 101
	err = q.UpdateChannelPointsReward(reward)
	a.Nil(err)

	dbr, err := q.GetChannelPointsReward(reward)
	a.Nil(err)
	rewards := dbr.Data.([]ChannelPointsReward)
	a.Len(rewards, 1)
	a.Equal(101, rewards[0].Cost)

	redemption := ChannelPointsRedemption{
		ID:               "1",
		BroadcasterID:    TEST_USER_ID,
		UserID:           "2",
		RedemptionStatus: "CANCELED",
		RewardID:         reward.ID,
		RedeemedAt:       util.GetTimestamp().Format(time.RFC3339),
	}

	err = q.InsertChannelPointsRedemption(redemption)
	a.Nil(err)

	redemption.RedemptionStatus = "TEST"
	err = q.UpdateChannelPointsRedemption(redemption)
	a.Nil(err)

	dbr, err = q.GetChannelPointsRedemption(redemption, "")
	a.Nil(err)
	redemptions := dbr.Data.([]ChannelPointsRedemption)
	a.Len(redemptions, 1)

	err = q.DeleteChannelPointsReward(redemption.RewardID)
	a.Nil(err)
}

func TestDrops(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	db, err := NewConnection()
	a.Nil(err)

	q := db.NewQuery(nil, 100)

	e := DropsEntitlement{
		ID:        "1",
		UserID:    TEST_USER_ID,
		BenefitID: util.RandomGUID(),
		GameID:    "1",
		Timestamp: util.GetTimestamp().Format(time.RFC3339),
	}

	err = q.InsertDropsEntitlement(e)
	a.Nil(err)

	dbr, err := q.GetDropsEntitlements(e)
	a.Nil(err)
	entitlements := dbr.Data.([]DropsEntitlement)
	a.Len(entitlements, 1)
	a.Equal(e.BenefitID, entitlements[0].BenefitID)
}

func TestErrors(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	a.False(DatabaseErrorIs(errors.New(""), sqlite3.ErrReadonlyRollback))
}

func TestModeration(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	db, err := NewConnection()
	a.Nil(err)
	q := db.NewQuery(nil, 100)

	urp := UserRequestParams{BroadcasterID: TEST_USER_ID, UserID: "2"}
	err = q.AddModerator(urp)
	a.Nil(err)

	dbr, err := q.GetModerationActionsByBroadcaster(TEST_USER_ID)
	a.Nil(err)
	moderatorActions := dbr.Data.([]ModeratorAction)
	a.Len(moderatorActions, 1)

	dbr, err = q.GetModerators(urp)
	a.Nil(err)
	moderators := dbr.Data.([]Moderator)
	a.Len(moderators, 1)

	dbr, err = q.GetModeratorsForBroadcaster(TEST_USER_ID, "2")
	a.Nil(err)
	moderators = dbr.Data.([]Moderator)
	a.Len(moderators, 1)

	dbr, err = q.GetModeratorEvents(urp)
	a.Nil(err)
	moderatorActions = dbr.Data.([]ModeratorAction)
	a.Len(moderatorActions, 1)

	err = q.RemoveModerator(urp.BroadcasterID, urp.UserID)
	a.Nil(err)

	err = q.InsertBan(urp)
	a.Nil(err)

	dbr, err = q.GetBans(urp)
	a.Nil(err)
	bans := dbr.Data.([]Ban)
	a.Len(bans, 1)

	dbr, err = q.GetBanEvents(urp)
	a.Nil(err)
	banEvents := dbr.Data.([]BanEvent)
	a.Len(banEvents, 1)

}

func TestPolls(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	db, err := NewConnection()
	a.Nil(err)
	q := db.NewQuery(nil, 100)

	poll := Poll{
		ID:                         "1",
		BroadcasterID:              TEST_USER_ID,
		Title:                      "test",
		BitsVotingEnabled:          false,
		ChannelPointsVotingEnabled: false,
		Status:                     "ACTIVE",
		Duration:                   150,
		StartedAt:                  util.GetTimestamp().Format(time.RFC3339),
		Choices: []PollsChoice{
			{
				ID:                 "1",
				Title:              "1234",
				Votes:              0,
				ChannelPointsVotes: 0,
				BitsVotes:          0,
				PollID:             "1",
			},
			{
				ID:                 "2",
				Title:              "234",
				Votes:              0,
				ChannelPointsVotes: 0,
				BitsVotes:          0,
				PollID:             "1",
			},
		},
	}

	err = q.InsertPoll(poll)
	a.Nil(err)

	err = q.UpdatePoll(Poll{ID: "1", BroadcasterID: TEST_USER_ID, Title: "test2"})
	a.Nil(err)

	dbr, err := q.GetPolls(Poll{ID: "1"})
	a.Nil(err)
	polls := dbr.Data.([]Poll)
	a.Len(polls, 1)
	a.Equal("test2", polls[0].Title)
}

func TestPredictions(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	db, err := NewConnection()
	a.Nil(err)
	q := db.NewQuery(nil, 100)

	prediction := Prediction{
		ID:               "1",
		BroadcasterID:    TEST_USER_ID,
		Title:            "1234",
		WinningOutcomeID: nil,
		PredictionWindow: 60,
		Status:           "ACTIVE",
		StartedAt:        util.GetTimestamp().Format(time.RFC3339),
		Outcomes: []PredictionOutcome{
			{
				ID:            "1",
				Title:         "111",
				Users:         0,
				ChannelPoints: 0,
				Color:         "BLUE",
				PredictionID:  "1",
			},
			{
				ID:            "2",
				Title:         "222",
				Users:         0,
				ChannelPoints: 0,
				Color:         "PINK",
				PredictionID:  "1",
			},
		},
	}

	err = q.InsertPrediction(prediction)
	a.Nil(err)

	predictionPredition := PredictionPrediction{
		PredictionID: "1",
		UserID:       TEST_USER_ID,
		Amount:       1000,
		OutcomeID:    prediction.Outcomes[0].ID,
	}

	err = q.InsertPredictionPrediction(predictionPredition)
	a.Nil(err)

	prediction.WinningOutcomeID = &prediction.Outcomes[0].ID
	err = q.UpdatePrediction(prediction)
	a.Nil(err)

	dbr, err := q.GetPredictions(Prediction{ID: "1"})
	a.Nil(err)
	predictions := dbr.Data.([]Prediction)
	a.Len(predictions, 1)
	prediction = predictions[0]
	a.NotNil(prediction.WinningOutcomeID)
	a.NotNil(prediction.Outcomes[0].TopPredictors)
}

func TestQuery(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	db, err := NewConnection()
	a.Nil(err)
	request, err := http.NewRequest(http.MethodGet, "https://google.com", nil)
	a.Nil(err)

	q := request.URL.Query()
	q.Set("first", "50")
	request.URL.RawQuery = q.Encode()

	query := db.NewQuery(request, 100)
	a.Equal(50, query.Limit)

	q.Set("after", query.PaginationCursor)
	request.URL.RawQuery = q.Encode()
	query = db.NewQuery(request, 100)

	q.Set("before", query.PaginationCursor)
	request.URL.RawQuery = q.Encode()
	query = db.NewQuery(request, 100)

	q.Set("after", "notbase64")
	request.URL.RawQuery = q.Encode()
	query = db.NewQuery(request, 100)
}

func TestStreams(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	db, err := NewConnection()
	a.Nil(err)
	q := db.NewQuery(nil, 100)
	s := Stream{
		ID:          "1",
		UserID:      TEST_USER_ID,
		StreamType:  "live",
		ViewerCount: 100,
		StartedAt:   util.GetTimestamp().Format(time.RFC3339),
		IsMature:    false,
	}
	err = q.InsertStream(s, false)
	a.Nil(err)

	tag := Tag{
		ID:     "1",
		Name:   "test",
		IsAuto: false,
	}

	err = q.InsertTag(tag)
	a.Nil(err)

	dbr, err := q.GetTags(tag)
	a.Nil(err)
	tags := dbr.Data.([]Tag)
	a.Len(tags, 1)

	err = q.InsertStreamTag(StreamTag{TagID: "1", UserID: "1"})
	a.Nil(err)

	dbr, err = q.GetStreamTags(TEST_USER_ID)
	a.Nil(err)
	tags = dbr.Data.([]Tag)
	a.Len(tags, 1)

	dbr, err = q.GetFollowedStreams(s.UserID)
	a.Nil(err)
	streams := dbr.Data.([]Stream)
	a.Len(streams, 0)

	err = q.AddFollow(UserRequestParams{BroadcasterID: s.UserID, UserID: "2"})
	a.Nil(err)

	dbr, err = q.GetFollowedStreams("2")
	a.Nil(err)
	streams = dbr.Data.([]Stream)
	a.Len(streams, 1)

	dbr, err = q.GetStream(s)
	a.Nil(err)
	streams = dbr.Data.([]Stream)
	a.Len(streams, 1)
	stream := streams[0]
	a.Len(stream.TagIDs, 1)

	err = q.DeleteAllStreamTags(s.UserID)
	a.Nil(err)

	v := Video{
		ID:               "1",
		StreamID:         &s.ID,
		BroadcasterID:    s.UserID,
		Title:            "1234",
		VideoDescription: "1234",
		CreatedAt:        util.GetTimestamp().Format(time.RFC3339),
		PublishedAt:      util.GetTimestamp().Format(time.RFC3339),
		Viewable:         "public",
		ViewCount:        100,
		Duration:         "1h0m0s",
		VideoLanguage:    "en",
		CategoryID:       nil,
		Type:             "archive",
	}

	err = q.InsertVideo(v)
	a.Nil(err)

	sm := StreamMarker{
		VideoID:         v.ID,
		CreatedAt:       util.GetTimestamp().Format(time.RFC3339),
		PositionSeconds: 10,
		Description:     "1234",
		BroadcasterID:   TEST_USER_ID,
		ID:              "1",
	}

	err = q.InsertStreamMarker(sm)
	a.Nil(err)

	dbr, err = q.GetStreamMarkers(StreamMarker{BroadcasterID: s.UserID})
	a.Nil(err)
	streamTags := dbr.Data.([]StreamMarkerUser)
	a.Len(streamTags, 1)
}

func TestSubscriptions(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	db, err := NewConnection()
	a.Nil(err)
	q := db.NewQuery(nil, 100)

	sub := Subscription{
		BroadcasterID: TEST_USER_ID,
		UserID:        "2",
		IsGift:        true,
		GifterID:      &sql.NullString{String: TEST_USER_ID, Valid: true},
		Tier:          "3000",
		CreatedAt:     util.GetTimestamp().Format(time.RFC3339),
	}

	err = q.InsertSubscription(SubscriptionInsert{BroadcasterID: sub.BroadcasterID, UserID: sub.UserID, IsGift: sub.IsGift, GifterID: sub.GifterID, Tier: sub.Tier, CreatedAt: sub.CreatedAt})
	a.Nil(err)

	dbr, err := q.GetSubscriptions(Subscription{BroadcasterID: sub.BroadcasterID, UserID: sub.UserID})
	a.Nil(err)
	subs := dbr.Data.([]Subscription)
	a.Len(subs, 1)
	a.Equal(subs[0].IsGift, true)
}

func TestTeams(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	db, err := NewConnection()
	a.Nil(err)
	q := db.NewQuery(nil, 100)

	team := Team{
		ID:              "1",
		CreatedAt:       util.GetTimestamp().Format(time.RFC3339),
		UpdatedAt:       util.GetTimestamp().Format(time.RFC3339),
		Info:            "",
		ThumbnailURL:    "",
		TeamName:        "test",
		TeamDisplayName: "Test",
	}

	err = q.InsertTeam(team)
	a.Nil(err)

	err = q.InsertTeamMember(TeamMember{TeamID: team.ID, UserID: TEST_USER_ID})
	a.Nil(err)

	dbr, err := q.GetTeam(Team{ID: team.ID})
	a.Nil(err)
	teams := dbr.Data.([]Team)
	a.Len(teams, 1)

	dbr, err = q.GetTeamByBroadcaster(TEST_USER_ID)
	a.Nil(err)
	teams = dbr.Data.([]Team)
	a.Len(teams, 1)

	dbr, err = q.GetTeamByBroadcaster("2")
	a.Nil(err)
	teams = dbr.Data.([]Team)
	a.Len(teams, 0)
}

func TestVideos(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	db, err := NewConnection()
	a.Nil(err)
	q := db.NewQuery(nil, 100)

	v := Video{
		ID:               util.RandomGUID(),
		StreamID:         nil,
		BroadcasterID:    TEST_USER_ID,
		Title:            "1234",
		VideoDescription: "1234",
		CreatedAt:        util.GetTimestamp().Format(time.RFC3339),
		PublishedAt:      util.GetTimestamp().Format(time.RFC3339),
		Viewable:         "public",
		ViewCount:        100,
		Duration:         "1h0m0s",
		VideoLanguage:    "en",
		CategoryID:       nil,
		Type:             "archive",
	}

	err = q.InsertVideo(v)
	a.Nil(err)

	vms := VideoMutedSegment{
		VideoID:     v.ID,
		VideoOffset: 20,
		Duration:    30,
	}

	err = q.InsertMutedSegmentsForVideo(vms)
	a.Nil(err)

	dbr, err := q.GetVideos(Video{ID: v.ID}, "", "time")
	a.Nil(err)
	videos := dbr.Data.([]Video)
	a.Len(videos, 1)
	a.Len(videos[0].MutedSegments, 1)

	c := Clip{
		ID:            "1",
		BroadcasterID: TEST_USER_ID,
		CreatorID:     TEST_USER_ID,
		VideoID:       vms.VideoID,
		GameID:        "1",
		Language:      "en",
		Title:         "?",
		ViewCount:     100,
		Duration:      1234.5,
		CreatedAt:     util.GetTimestamp().Format(time.RFC3339),
	}

	err = q.InsertClip(c)
	a.Nil(err)

	dbr, err = q.GetClips(Clip{ID: c.ID}, "", "")
	a.Nil(err)
	clips := dbr.Data.([]Clip)
	a.Len(clips, 1)

	err = q.DeleteVideo(vms.VideoID)
	a.Nil(err)
}
