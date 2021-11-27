// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package user_update

import (
    "encoding/json"
    "time"

    "github.com/twitchdev/twitch-cli/internal/events"
    "github.com/twitchdev/twitch-cli/internal/models"
    "github.com/twitchdev/twitch-cli/internal/util"
)

var transportsSupported = map[string]bool{
    models.TransportEventSub: true,
}
var triggers = []string{"user.update"}

var triggerMapping = map[string]map[string]string{
    models.TransportEventSub: {
        "user.update": "user.update",
    },
}

type Event struct{}

func (e Event) GenerateEvent(p events.MockEventParameters) (events.MockEventResponse, error) {
    var event []byte
    var err error

    switch p.Transport {
    case models.TransportEventSub:
        body := models.EventsubResponse{
            Subscription: models.EventsubSubscription{
                ID:      p.ID,
                Status:  "enabled",
                Type:    "user.update",
                Version: "1",
                Condition: models.EventsubCondition{
                    UserID : p.ToUserID,
                },
                Transport: models.EventsubTransport{
                    Method:   "webhook",
                    Callback: "null",
                },
                Cost:      0,
                CreatedAt: util.GetTimestamp().Format(time.RFC3339Nano),
            },
            Event: models.UserUpdateEventSubEvent{
                UserID:         p.ToUserID,
                UserLogin:      p.ToUserName,
                UserName:       p.ToUserName,
                Description:    p.Description,
            },
        }

        event, err = json.Marshal(body)
        if err != nil {
            return events.MockEventResponse{}, err
        }
    default:
        return events.MockEventResponse{}, nil
    }

    return events.MockEventResponse{
        ID:             p.ID,
        JSON:           event,
        FromUser:       p.FromUserID,
        ToUser:         p.ToUserID,
    }, nil
}

func (e Event) ValidTransport(transport string) bool {
    return transportsSupported[transport]
}

func (e Event) ValidTrigger(trigger string) bool {
    for _, t := range triggers {
        if t == trigger {
            return true
        }
    }
    return false
}
func (e Event) GetTopic(transport string, trigger string) string {
    return triggerMapping[transport][trigger]
}
func (e Event) GetEventSubAlias(t string) string {
    // check for aliases
    for trigger, topic := range triggerMapping[models.TransportEventSub] {
        if topic == t {
            return trigger
        }
    }
    return ""
}
