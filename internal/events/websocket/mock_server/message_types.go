package mock_server

import "github.com/twitchdev/twitch-cli/internal/models"

// Generic response message Metadata; Always the same

type MessageMetadata struct {
	MessageID        string `json:"message_id"`
	MessageType      string `json:"message_type"`
	MessageTimestamp string `json:"message_timestamp"`

	SubscriptionType    string `json:"subscription_type,omitempty"`
	SubscriptionVersion string `json:"subscription_version,omitempty"`
}

/* ** Welcome message **
{ // <1>
	"metadata": { // <MessageMetadata>
		"message_id": "befa7b53-d79d-478f-86b9-120f112b044e",
		"message_type": "session_welcome",
		"message_timestamp": "2019-11-16T10:11:12.123Z"
	},
	"payload": { // <2>
		"session": { // <3>
			"id": "AQoQexAWVYKSTIu4ec_2VAxyuhAB",
			"status": "connected",
			"keepalive_timeout_seconds": 10,
			"reconnect_url": null,
			"connected_at": "2019-11-16T10:11:12.123Z"
		}
	}
}
*/

type WelcomeMessage struct { // <1>
	Metadata MessageMetadata       `json:"metadata"`
	Payload  WelcomeMessagePayload `json:"payload"`
}

type WelcomeMessagePayload struct { // <2>
	Session WelcomeMessagePayloadSession `json:"session"`
}

type WelcomeMessagePayloadSession struct { // <3>
	ID                      string  `json:"id"`
	Status                  string  `json:"status"`
	KeepaliveTimeoutSeconds int     `json:"keepalive_timeout_seconds"`
	ReconnectUrl            *string `json:"reconnect_url"`
	ConnectedAt             string  `json:"connected_at"`
}

/* ** Reconnect message **
{ // <1>
	"metadata": { // <MessageMetadata>
		"message_id": "84c1e79a-2a4b-4c13-ba0b-4312293e9308",
		"message_type": "session_reconnect",
		"message_timestamp": "2019-11-18T09:10:11.234Z"
	},
	"payload": { // <2>
		"session": { // <3>
			"id": "AQoQexAWVYKSTIu4ec_2VAxyuhAB",
			"status": "reconnecting",
			"keepalive_timeout_seconds": null,
			"reconnect_url": "wss://eventsub-experimental.wss.twitch.tv?...",
			"connected_at": "2019-11-16T10:11:12.123Z"
		}
	}
}
*/

type ReconnectMessage struct { // <1>
	Metadata MessageMetadata         `json:"metadata"`
	Payload  ReconnectMessagePayload `json:"payload"`
}

type ReconnectMessagePayload struct { // <2>
	Session ReconnectMessagePayloadSession `json:"session"`
}

type ReconnectMessagePayloadSession struct { // <3>
	ID                      string `json:"id"`
	Status                  string `json:"status"`
	KeepaliveTimeoutSeconds *int   `json:"keepalive_timeout_seconds"`
	ReconnectUrl            string `json:"reconnect_url"`
	ConnectedAt             string `json:"connected_at"`
}

/* ** Keepalive message **
{ // <1>
    "metadata": { // <MessageMetadata>
        "message_id": "84c1e79a-2a4b-4c13-ba0b-4312293e9308",
        "message_type": "session_keepalive",
        "message_timestamp": "2019-11-16T10:11:12.123Z"
    },
    "payload": {} // struct{}
}
*/

type KeepaliveMessage struct { // <1>
	Metadata MessageMetadata         `json:"metadata"`
	Payload  KeepaliveMessagePayload `json:"payload"`
}

type KeepaliveMessagePayload struct{}

/* ** Notification message **
{ // <1>
    "metadata": { // <MessageMetadata>
        "message_id": "befa7b53-d79d-478f-86b9-120f112b044e",
        "message_type": "notification",
        "message_timestamp": "2019-11-16T10:11:12.464757833Z",
        "subscription_type": "channel.follow",
        "subscription_version": "1"
    },
    "payload": { // <models.EventsubResponse>
        "subscription": {
           ...
        },
        "event": {
            ...
        }
    }
}
*/

type NotificationMessage struct { // <1>
	Metadata MessageMetadata         `json:"metadata"`
	Payload  models.EventsubResponse `json:"payload"`
}
