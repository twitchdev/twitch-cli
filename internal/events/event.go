// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package events

// MockEventParameters are used to craft the event; most of this data is prepopulated by lower services, such as the from/to users to avoid
// replicating logic across files
type MockEventParameters struct {
	ID           string
	Transport    string
	Trigger      string
	FromUserID   string
	FromUserName string
	ToUserID     string
	ToUserName   string
	IsAnonymous  bool
	IsGift       bool
	Status       string
	ItemID       string
	ItemName     string
	Cost         int64
	IsPermanent  bool
	Description  string
}

type MockEventResponse struct {
	ID        string
	JSON      []byte
	FromUser  string
	ToUser    string
	Timestamp string
}

// MockEvent represents an event to be triggered using the `twitch event trigger <event>` command. Given that
// both WebSub and EventSub need to be supported, it's required to have logic for both currently.
type MockEvent interface {

	// Returns the Mock Response for the given transport
	GenerateEvent(p MockEventParameters) (MockEventResponse, error)

	// Returns the trigger for the event (e.g. cheer for cheer events, or add-reward for channel points add rewards)
	ValidTrigger(trigger string) bool

	// Returns whether a given event supports a supplied transport
	ValidTransport(transport string) bool

	// Returns the string of the topic
	GetTopic(transport string, trigger string) string
}
