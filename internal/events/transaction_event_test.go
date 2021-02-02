// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package events

import (
	"encoding/json"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

func TestWebusbTransaction(t *testing.T) {
	a := util.SetupTestEnv(t)

	params := *&TransactionParams{
		FromUser:  fromUser,
		ToUser:    toUser,
		Transport: TransportWebSub,
	}

	r, err := GenerateTransactionBody(params)
	a.Nil(err)

	var body models.TransactionWebSubResponse
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal(toUser, body.Data[0].BroadcasterID, "Expected to user %v, got %v", toUser, body.Data[0].BroadcasterID)
	a.Equal(fromUser, body.Data[0].UserID, "Expected from user %v, got %v", r.ToUser, body.Data[0].UserID)

	params = *&TransactionParams{
		Transport: TransportWebSub,
	}

	r, err = GenerateTransactionBody(params)
	a.Nil(err)

	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.NotNil(body.Data[0].BroadcasterID)
	a.NotNil(body.Data[0].UserID)

}

func TestEventsubTransaction(t *testing.T) {
	a := util.SetupTestEnv(t)

	params := *&TransactionParams{
		FromUser:  fromUser,
		ToUser:    toUser,
		Transport: TransportEventSub,
	}

	_, err := GenerateTransactionBody(params)
	a.NotNil(err)
}
