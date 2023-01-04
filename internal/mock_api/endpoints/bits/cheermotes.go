// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package bits

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var cheermotesMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var cheermotesScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type Cheermotes struct{}

type CheermotesResponse struct {
	Data []CheermotesResponseData `json:"data"`
}

type CheermotesResponseData struct {
	Prefix       string      `json:"prefix"`
	Order        int         `json:"order"`
	LastUpdated  string      `json:"last_updated"`
	IsCharitable bool        `json:"is_charitable"`
	Tiers        []Cheermote `json:"tiers"`
	Type         string      `json:"type"`
}

type Cheermote struct {
	MinBits        int            `json:"min_bits"`
	ID             string         `json:"id"`
	Color          string         `json:"color"`
	Images         CheermoteImage `json:"images"`
	CanCheer       bool           `json:"can_cheer"`
	ShowInBitsCard bool           `json:"show_in_bits_card"`
}

type CheermoteImage struct {
	Dark  CheermoteImageType `json:"dark"`
	Light CheermoteImageType `json:"light"`
}
type CheermoteImageType struct {
	Animated CheermoteImageSizes `json:"animated"`
	Static   CheermoteImageSizes `json:"static"`
}

type CheermoteImageSizes struct {
	One         string `json:"1"`
	OneAndAHalf string `json:"1.5"`
	Two         string `json:"2"`
	Three       string `json:"3"`
	Four        string `json:"4"`
}

var defaultMinCheers = []int{1, 10, 100, 1000, 5000, 10000}
var defaultTypes = []string{"global_first_party", "global_third_party", "channel_custom", "display_only", "sponsored"}
var defaultPrefixes = []string{"Cheer", "Charity", "Concrete", "Goal", "CLI"}

func (e Cheermotes) Path() string { return "/bits/cheermotes" }

func (e Cheermotes) GetRequiredScopes(method string) []string {
	return cheermotesScopesByMethod[method]
}

func (e Cheermotes) ValidMethod(method string) bool {
	return cheermotesMethodsSupported[method]
}

func (e Cheermotes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getCheermotes(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getCheermotes(w http.ResponseWriter, r *http.Request) {
	cheermoteBody := CheermotesResponse{}
	for i := 0; i < len(defaultTypes); i++ {
		cheermote := CheermotesResponseData{
			Prefix:       defaultPrefixes[i],
			Order:        i + 1,
			LastUpdated:  util.GetTimestamp().Format(time.RFC3339),
			IsCharitable: defaultPrefixes[i] == "Charity",
			Tiers:        []Cheermote{},
			Type:         defaultTypes[i],
		}

		// generate the mock object
		for _, bits := range defaultMinCheers {
			cheermote.Tiers = append(cheermote.Tiers, Cheermote{
				MinBits:        bits,
				ID:             fmt.Sprint(bits),
				Color:          "#9146FF",
				CanCheer:       true,
				ShowInBitsCard: true,
				Images: CheermoteImage{
					Dark: CheermoteImageType{
						Animated: generateCheermoteImageSizes(strings.ToLower(defaultPrefixes[i]), "dark", "animated", bits),
						Static:   generateCheermoteImageSizes(strings.ToLower(defaultPrefixes[i]), "dark", "static", bits),
					},
					Light: CheermoteImageType{
						Animated: generateCheermoteImageSizes(strings.ToLower(defaultPrefixes[i]), "light", "animated", bits),
						Static:   generateCheermoteImageSizes(strings.ToLower(defaultPrefixes[i]), "light", "static", bits),
					},
				},
			})
		}

		cheermoteBody.Data = append(cheermoteBody.Data, cheermote)
	}

	response, _ := json.Marshal(cheermoteBody)
	w.Write(response)
}

func generateCheermoteImageSizes(prefix string, theme string, imageType string, bits int) CheermoteImageSizes {
	fileType := "png"
	if imageType == "animated" {
		fileType = "gif"
	}
	return CheermoteImageSizes{
		One:         fmt.Sprintf("https://d3aqoihi2n8ty8.cloudfront.net/actions/%v/%v/%v/%v/1.%v", prefix, theme, imageType, fmt.Sprint(bits), fileType),
		OneAndAHalf: fmt.Sprintf("https://d3aqoihi2n8ty8.cloudfront.net/actions/%v/%v/%v/%v/1.5.%v", prefix, theme, imageType, fmt.Sprint(bits), fileType),
		Two:         fmt.Sprintf("https://d3aqoihi2n8ty8.cloudfront.net/actions/%v/%v/%v/%v/2.%v", prefix, theme, imageType, fmt.Sprint(bits), fileType),
		Three:       fmt.Sprintf("https://d3aqoihi2n8ty8.cloudfront.net/actions/%v/%v/%v/%v/3.%v", prefix, theme, imageType, fmt.Sprint(bits), fileType),
		Four:        fmt.Sprintf("https://d3aqoihi2n8ty8.cloudfront.net/actions/%v/%v/%v/%v/4.%v", prefix, theme, imageType, fmt.Sprint(bits), fileType),
	}
}
