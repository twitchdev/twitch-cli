// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package drops

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/spf13/viper"
	"github.com/twitchdev/twitch-cli/internal/api"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/request"
	"github.com/twitchdev/twitch-cli/internal/util"
	"golang.org/x/time/rate"
)

type DropsEntitlementExportParameters struct {
	GameID string
	UserID string
	URL    string
	Cursor string
}

var (
	TOKEN        string
	CLIENT_ID    string
	ENTITLEMENTS []models.DropsEntitlementsData
)

var BASE_URL = "https://api.twitch.tv/helix/entitlements/drops"

func ExportEntitlements(filename string, gameID string, userID string) {
	c, err := api.GetClientInformation()
	if err != nil {
		return
	}

	CLIENT_ID = c.ClientID
	TOKEN = c.Token

	if viper.GetString("BASE_URL") != "" {
		BASE_URL = viper.GetString("BASE_URL")
	}

	p := DropsEntitlementExportParameters{
		URL:    BASE_URL,
		GameID: gameID,
		UserID: userID,
	}
	for {
		var e models.DropsEntitlementsResponse
		resp, err := makeAPIRequest(p)
		body, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			fmt.Printf("Error reading body: %v", err)
			return
		}
		err = json.Unmarshal(body, &e)
		if err != nil {
			fmt.Printf("Error reading body: %v", err)
			return
		}

		ENTITLEMENTS = append(ENTITLEMENTS, e.Data...)

		if resp.StatusCode == 500 {
			fmt.Println(fmt.Sprintf("[%v] Got 500 from endpoint; Make sure that your client is marked in the correct organization.", util.GetTimestamp().Format(time.RFC3339)))
			break
		}

		if len(e.Data) == 0 {
			fmt.Println("No results, stopping.")
			break
		}

		if e.Pagination.Cursor == "" {
			fmt.Println(fmt.Sprintf("[%v] End of records, found %v records.", util.GetTimestamp().Format(time.RFC3339), len(ENTITLEMENTS)))
			break
		}

		p.Cursor = e.Pagination.Cursor

		fmt.Println(fmt.Sprintf("[%v] Found %v records, hitting next page.", util.GetTimestamp().Format(time.RFC3339), len(e.Data)))

	}

	// don't make a csv if empty :)
	if len(ENTITLEMENTS) == 0 {
		return
	}

	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error writing file %v", err)
		return
	}

	w := csv.NewWriter(file)

	headers := []string{
		"id",
		"user_id",
		"benefit_id",
		"timestamp",
		"game_id",
	}

	w.Write(headers)

	for _, e := range ENTITLEMENTS {
		v := make([]string, 0)
		v = append(v,
			e.ID,
			e.UserID,
			e.BenefitID,
			e.Timestamp,
			e.GameID,
		)

		w.Write(v)
	}
	//finish writing
	w.Flush()
}

func makeAPIRequest(p DropsEntitlementExportParameters) (*http.Response, error) {
	u, err := url.Parse(p.URL)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Add("first", "100")

	if p.GameID != "" {
		q.Add("game_id", p.GameID)
	}
	if p.UserID != "" {
		q.Add("user_id", p.UserID)
	}
	if p.Cursor != "" {
		q.Add("after", p.Cursor)
	}

	u.RawQuery = q.Encode()

	req, _ := request.NewRequest(http.MethodGet, u.String(), nil)
	req.Header.Set("Client-ID", CLIENT_ID)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", TOKEN))

	rl := rate.NewLimiter(rate.Every(10*time.Second), 100)
	client := NewClient(rl)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
