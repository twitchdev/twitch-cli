// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package trigger

import (
	"encoding/json"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/models"
)

func TestWebusbTransaction(t *testing.T) {
	params := *&TransactionParams{
		FromUser:  fromUser,
		ToUser:    toUser,
		Transport: "websub",
	}

	r, err := GenerateTransactionBody(params)
	if err != nil {
		t.Error(err)
	}

	var body models.TransactionWebSubResponse
	if err = json.Unmarshal(r.JSON, &body); err != nil {
		t.Error("Error unmarshalling JSON")
	}

	if body.Data[0].BroadcasterID != toUser {
		t.Errorf("Expected to user %v, got %v", toUser, body.Data[0].BroadcasterID)
	}

	if body.Data[0].UserID != fromUser {
		t.Errorf("Expected from user %v, got %v", r.ToUser, body.Data[0].UserID)
	}
}
